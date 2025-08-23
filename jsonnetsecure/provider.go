package jsonnetsecure

import (
	"context"
	"os"
	"runtime"
	"testing"
)

type (
	VMProvider interface {
		// JsonnetVM creates a new secure process-isolated Jsonnet VM whose
		// execution is bound to the provided context, i.e.,
		// cancelling the context will terminate the VM process.
		JsonnetVM(context.Context) (VM, error)
	}

	// TestProvider provides a secure VM by running go build on github.
	// com/ory/x/jsonnetsecure/cmd.
	TestProvider struct {
		jsonnetBinary string
		pool          Pool
	}

	// DefaultProvider provides a secure VM by calling the currently
	// running the current binary with the provided subcommand.
	DefaultProvider struct {
		Subcommand string
		Pool       Pool
	}

	vmOptions struct {
		jsonnetBinaryPath string
		args              []string
		ctx               context.Context
		pool              *pool
	}

	Option func(o *vmOptions)
)

func newVMOptions() *vmOptions {
	jsonnetBinaryPath, _ := os.Executable()
	return &vmOptions{
		jsonnetBinaryPath: jsonnetBinaryPath,
		ctx:               context.Background(),
	}
}

func WithProcessPool(p Pool) Option {
	return func(o *vmOptions) {
		pool, _ := p.(*pool)
		o.pool = pool
	}
}

func WithJsonnetBinary(jsonnetBinaryPath string) Option {
	return func(o *vmOptions) {
		o.jsonnetBinaryPath = jsonnetBinaryPath
	}
}

func WithProcessArgs(args ...string) Option {
	return func(o *vmOptions) {
		o.args = args
	}
}

func NewTestProvider(t testing.TB) *TestProvider {
	pool := NewProcessPool(runtime.GOMAXPROCS(0))
	t.Cleanup(pool.Close)
	return &TestProvider{JsonnetTestBinary(t), pool}
}

func (p *TestProvider) JsonnetVM(ctx context.Context) (VM, error) {
	return MakeSecureVM(
		WithProcessPool(p.pool),
		WithJsonnetBinary(p.jsonnetBinary),
	), nil
}

func (p *DefaultProvider) JsonnetVM(ctx context.Context) (VM, error) {
	self, err := os.Executable()
	if err != nil {
		return nil, err
	}
	return MakeSecureVM(
		WithJsonnetBinary(self),
		WithProcessArgs(p.Subcommand),
		WithProcessPool(p.Pool),
	), nil
}
