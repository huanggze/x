package jsonnetsecure

import (
	"context"
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
)

func NewTestProvider(t testing.TB) *TestProvider {
	pool := NewProcessPool(runtime.GOMAXPROCS(0))
	t.Cleanup(pool.Close)
	return &TestProvider{JsonnetTestBinary(t), pool}
}
