package jsonnetsecure

import (
	"runtime"
	"testing"
)

type (
	VMProvider interface {
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
