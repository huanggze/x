package otelx

import (
	"github.com/huanggze/x/logrusx"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Tracer struct {
	tracer trace.Tracer
}

// Creates a new no-op tracer.
func NewNoop(_ *logrusx.Logger, c *Config) *Tracer {
	tp := noop.NewTracerProvider()
	t := &Tracer{tracer: tp.Tracer("")}
	return t
}
