package otelx

import (
	"github.com/huanggze/x/logrusx"
	"github.com/huanggze/x/stringsx"

	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

type Tracer struct {
	tracer trace.Tracer
}

// Creates a new tracer. If name is empty, a default tracer name is used
// instead. See: https://godocs.io/go.opentelemetry.io/otel/sdk/trace#TracerProvider.Tracer
func New(name string, l *logrusx.Logger, c *Config) (*Tracer, error) {
	t := &Tracer{}

	if err := t.setup(name, l, c); err != nil {
		return nil, err
	}

	return t, nil
}

// Creates a new no-op tracer.
func NewNoop(_ *logrusx.Logger, c *Config) *Tracer {
	tp := noop.NewTracerProvider()
	t := &Tracer{tracer: tp.Tracer("")}
	return t
}

// setup constructs the tracer based on the given configuration.
func (t *Tracer) setup(name string, l *logrusx.Logger, c *Config) error {
	switch f := stringsx.SwitchExact(c.Provider); {
	case f.AddCase("jaeger"):
		tracer, err := SetupJaeger(t, name, c)
		if err != nil {
			return err
		}

		t.tracer = tracer
		l.Infof("Jaeger tracer configured! Sending spans to %s", c.Providers.Jaeger.LocalAgentAddress)
	case f.AddCase("zipkin"):
		tracer, err := SetupZipkin(t, name, c)
		if err != nil {
			return err
		}

		t.tracer = tracer
		l.Infof("Zipkin tracer configured! Sending spans to %s", c.Providers.Zipkin.ServerURL)
	case f.AddCase("otel"):
		tracer, err := SetupOTLP(t, name, c)
		if err != nil {
			return err
		}

		t.tracer = tracer
		l.Infof("OTLP tracer configured! Sending spans to %s", c.Providers.OTLP.ServerURL)
	case f.AddCase(""):
		l.Infof("No tracer configured - skipping tracing setup")
		t.tracer = noop.NewTracerProvider().Tracer(name)
	default:
		return f.ToUnknownCaseErr()
	}

	return nil
}
