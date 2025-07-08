package httpx

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.opentelemetry.io/otel/trace"

	"github.com/huanggze/x/logrusx"
)

type resilientOptions struct {
	c                    *http.Client
	l                    interface{}
	retryWaitMin         time.Duration
	retryWaitMax         time.Duration
	retryMax             int
	noInternalIPs        bool
	internalIPExceptions []string
	ipV6                 bool
	tracer               trace.Tracer
}

func newResilientOptions() *resilientOptions {
	connTimeout := time.Minute
	return &resilientOptions{
		c:            &http.Client{Timeout: connTimeout},
		retryWaitMin: 1 * time.Second,
		retryWaitMax: 30 * time.Second,
		retryMax:     4,
		l:            log.New(io.Discard, "", log.LstdFlags),
		ipV6:         true,
	}
}

// ResilientOptions is a set of options for the ResilientClient.
type ResilientOptions func(o *resilientOptions)

// ResilientClientWithTracer wraps the http clients transport with a tracing instrumentation
func ResilientClientWithTracer(tracer trace.Tracer) ResilientOptions {
	return func(o *resilientOptions) {
		o.tracer = tracer
	}
}

// ResilientClientWithMaxRetry sets the maximum number of retries.
func ResilientClientWithMaxRetry(retryMax int) ResilientOptions {
	return func(o *resilientOptions) {
		o.retryMax = retryMax
	}
}

// ResilientClientWithConnectionTimeout sets the connection timeout for the client.
func ResilientClientWithConnectionTimeout(connTimeout time.Duration) ResilientOptions {
	return func(o *resilientOptions) {
		o.c.Timeout = connTimeout
	}
}

// ResilientClientWithLogger sets the logger to be used by the client.
func ResilientClientWithLogger(l *logrusx.Logger) ResilientOptions {
	return func(o *resilientOptions) {
		o.l = l
	}
}

// ResilientClientDisallowInternalIPs disallows internal IPs from being used.
func ResilientClientDisallowInternalIPs() ResilientOptions {
	return func(o *resilientOptions) {
		o.noInternalIPs = true
	}
}

// NewResilientClient creates a new ResilientClient.
func NewResilientClient(opts ...ResilientOptions) *retryablehttp.Client {
	o := newResilientOptions()
	for _, f := range opts {
		f(o)
	}

	if o.noInternalIPs {
		o.c.Transport = &noInternalIPRoundTripper{
			onWhitelist:          ifelse(o.ipV6, allowInternalAllowIPv6, allowInternalProhibitIPv6),
			notOnWhitelist:       ifelse(o.ipV6, prohibitInternalAllowIPv6, prohibitInternalProhibitIPv6),
			internalIPExceptions: o.internalIPExceptions,
		}
	} else {
		o.c.Transport = ifelse(o.ipV6, allowInternalAllowIPv6, allowInternalProhibitIPv6)
	}

	cl := retryablehttp.NewClient()
	cl.HTTPClient = o.c
	cl.Logger = o.l
	cl.RetryWaitMin = o.retryWaitMin
	cl.RetryWaitMax = o.retryWaitMax
	cl.RetryMax = o.retryMax
	cl.CheckRetry = retryablehttp.DefaultRetryPolicy
	cl.Backoff = retryablehttp.DefaultBackoff

	return cl
}

func ifelse[A any](b bool, x, y A) A {
	if b {
		return x
	}
	return y
}
