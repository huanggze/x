package jsonnetsecure

import (
	"bufio"
	"context"
	"io"
	"math"
	"os/exec"
	"time"

	"github.com/jackc/puddle/v2"
	"github.com/pkg/errors"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/huanggze/x/otelx"
)

type (
	Pool interface {
		Close()
		private()
	}
	pool struct {
		puddle *puddle.Pool[worker]
	}
	worker struct {
		cmd    *exec.Cmd
		stdin  chan<- []byte
		stdout <-chan string
		stderr <-chan string
	}
	contextKeyType string
)

const (
	contextValuePath contextKeyType = "argc"
	contextValueArgs contextKeyType = "argv"
)

func NewProcessPool(size int) Pool {
	size = max(5, min(size, math.MaxInt32))
	pud, err := puddle.NewPool(&puddle.Config[worker]{
		MaxSize:     int32(size),
		Constructor: newWorker,
		Destructor:  worker.destroy,
	})
	if err != nil {
		panic(err) // this should never happen, see implementation of puddle.NewPool
	}
	for range size {
		// warm pool
		go pud.CreateResource(context.Background())
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			for _, proc := range pud.AcquireAllIdle() {
				if proc.Value().cmd.ProcessState != nil {
					proc.Destroy()
				} else {
					proc.Release()
				}
			}
		}
	}()
	return &pool{pud}
}

func (*pool) private() {}

func (p *pool) Close() {
	p.puddle.Close()
}
func newWorker(ctx context.Context) (_ worker, err error) {
	tracer := trace.SpanFromContext(ctx).TracerProvider().Tracer("")
	ctx, span := tracer.Start(ctx, "jsonnetsecure.newWorker")
	defer otelx.End(span, &err)

	path, _ := ctx.Value(contextValuePath).(string)
	if path == "" {
		return worker{}, errors.New("newWorker: missing binary path in context")
	}
	args, _ := ctx.Value(contextValueArgs).([]string)
	cmd := exec.Command(path, append(args, "-0")...)
	cmd.Env = []string{"GOMAXPROCS=1"}
	cmd.WaitDelay = 100 * time.Millisecond

	span.SetAttributes(semconv.ProcessCommand(cmd.Path), semconv.ProcessCommandArgs(cmd.Args...))

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return worker{}, errors.Wrap(err, "newWorker: failed to create stdin pipe")
	}

	in := make(chan []byte, 1)
	go func(c <-chan []byte) {
		for input := range c {
			if _, err := stdin.Write(append(input, 0)); err != nil {
				stdin.Close()
				return
			}
		}
	}(in)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return worker{}, errors.Wrap(err, "newWorker: failed to create stdout pipe")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return worker{}, errors.Wrap(err, "newWorker: failed to create stderr pipe")
	}

	if err := cmd.Start(); err != nil {
		return worker{}, errors.Wrap(err, "newWorker: failed to start process")
	}

	span.SetAttributes(semconv.ProcessPID(cmd.Process.Pid))

	scan := func(c chan<- string, r io.Reader) {
		defer close(c)
		// NOTE: `bufio.Scanner` has its own internal limit of 64 KiB.
		scanner := bufio.NewScanner(r)

		scanner.Split(splitNull)
		for scanner.Scan() {
			c <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			c <- "ERROR: scan: " + err.Error()
		}
	}
	out := make(chan string, 1)
	go scan(out, stdout)
	errs := make(chan string, 1)
	go scan(errs, stderr)

	w := worker{
		cmd:    cmd,
		stdin:  in,
		stdout: out,
		stderr: errs,
	}

	_, err = w.eval(ctx, []byte("{}")) // warm up
	if err != nil {
		w.destroy()
		return worker{}, errors.Wrap(err, "newWorker: warm up failed")
	}

	return w, nil
}

func (w worker) destroy() {
	close(w.stdin)
	w.cmd.Process.Kill()
	w.cmd.Wait()
}

func (w worker) eval(ctx context.Context, processParams []byte) (output string, err error) {
	tracer := trace.SpanFromContext(ctx).TracerProvider().Tracer("")
	ctx, span := tracer.Start(ctx, "jsonnetsecure.worker.eval", trace.WithAttributes(
		semconv.ProcessPID(w.cmd.Process.Pid)))
	defer otelx.End(span, &err)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case w.stdin <- processParams:
		break
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case output := <-w.stdout:
		return output, nil
	case err := <-w.stderr:
		return "", errors.New(err)
	}
}
