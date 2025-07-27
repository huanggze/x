package reqlog

import (
	"context"
	"sync"
	"time"
)

// WithEnableExternalLatencyMeasurement returns a context that measures external latencies.
func WithEnableExternalLatencyMeasurement(ctx context.Context) context.Context {
	container := contextContainer{
		latencies: make([]externalLatency, 0),
	}
	return context.WithValue(ctx, externalLatencyKey, &container)
}

// totalExternalLatency returns the total duration of all external calls.
func totalExternalLatency(ctx context.Context) (total time.Duration) {
	if _, ok := ctx.Value(disableExternalLatencyMeasurement).(bool); ok {
		return 0
	}
	container, ok := ctx.Value(externalLatencyKey).(*contextContainer)
	if !ok {
		return 0
	}

	container.Lock()
	defer container.Unlock()
	for _, l := range container.latencies {
		total += l.Took
	}
	return total
}

type (
	externalLatency = struct {
		Took          time.Duration
		Cause, Detail string
	}
	contextContainer = struct {
		latencies []externalLatency
		sync.Mutex
	}
	contextKey int
)

const (
	externalLatencyKey                contextKey = 1
	disableExternalLatencyMeasurement contextKey = 2
)
