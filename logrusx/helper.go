package logrusx

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/huanggze/x/errorsx"
)

type (
	Logger struct {
		*logrus.Entry
		leakSensitive bool
		redactionText string
		opts          []Option
		name          string
		version       string
	}
)

var opts = otelhttptrace.WithPropagators(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

func (l *Logger) Logrus() *logrus.Logger {
	return l.Entry.Logger
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	ll := *l
	ll.Entry = l.Logger.WithContext(ctx)
	return &ll
}

func (l *Logger) HTTPHeadersRedacted(h http.Header) map[string]interface{} {
	headers := map[string]interface{}{}

	for key, value := range h {
		switch keyLower := strings.ToLower(key); keyLower {
		case "authorization", "cookie", "set-cookie", "x-session-token":
			headers[keyLower] = l.maybeRedact(value)
		case "location":
			locationURL, err := url.Parse(h.Get("Location"))
			if err != nil {
				headers[keyLower] = l.maybeRedact(value)
				continue
			}
			if l.leakSensitive {
				headers[keyLower] = locationURL.String()
			} else {
				locationURL.RawQuery = ""
				locationURL.Fragment = ""
				headers[keyLower] = locationURL.Redacted()
			}
		default:
			headers[keyLower] = h.Get(key)
		}
	}

	return headers
}

func (l *Logger) WithRequest(r *http.Request) *Logger {
	headers := l.HTTPHeadersRedacted(r.Header)
	if ua := r.UserAgent(); len(ua) > 0 {
		headers["user-agent"] = ua
	}

	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}

	ll := l.WithField("http_request", map[string]interface{}{
		"remote":  r.RemoteAddr,
		"method":  r.Method,
		"path":    r.URL.EscapedPath(),
		"query":   l.maybeRedact(r.URL.RawQuery),
		"scheme":  scheme,
		"host":    r.Host,
		"headers": headers,
	})

	spanCtx := trace.SpanContextFromContext(r.Context())
	if !spanCtx.IsValid() {
		_, _, spanCtx = otelhttptrace.Extract(r.Context(), r, opts)
	}
	if spanCtx.IsValid() {
		traces := make(map[string]string, 2)
		if spanCtx.HasTraceID() {
			traces["trace_id"] = spanCtx.TraceID().String()
		}
		if spanCtx.HasSpanID() {
			traces["span_id"] = spanCtx.SpanID().String()
		}
		ll = ll.WithField("otel", traces)
	}
	return ll
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, format, args...)
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	ll := *l
	ll.Entry = l.Entry.WithField(key, value)
	return &ll
}

func (l *Logger) maybeRedact(value interface{}) interface{} {
	if fmt.Sprintf("%v", value) == "" || value == nil {
		return nil
	}
	if !l.leakSensitive {
		return l.redactionText
	}
	return value
}

func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}

	ctx := map[string]interface{}{"message": err.Error()}
	if l.Entry.Logger.IsLevelEnabled(logrus.DebugLevel) {
		if e, ok := err.(errorsx.StackTracer); ok {
			ctx["stack_trace"] = fmt.Sprintf("%+v", e.StackTrace())
		} else {
			ctx["stack_trace"] = fmt.Sprintf("stack trace could not be recovered from error type %s", reflect.TypeOf(err))
		}
	}

	if c := errorsx.ReasonCarrier(nil); errors.As(err, &c) {
		ctx["reason"] = c.Reason()
	}
	if c := errorsx.RequestIDCarrier(nil); errors.As(err, &c) && c.RequestID() != "" {
		ctx["request_id"] = c.RequestID()
	}
	if c := errorsx.DetailsCarrier(nil); errors.As(err, &c) && c.Details() != nil {
		ctx["details"] = c.Details()
	}
	if c := errorsx.StatusCarrier(nil); errors.As(err, &c) && c.Status() != "" {
		ctx["status"] = c.Status()
	}
	if c := errorsx.StatusCodeCarrier(nil); errors.As(err, &c) && c.StatusCode() != 0 {
		ctx["status_code"] = c.StatusCode()
	}
	if c := errorsx.DebugCarrier(nil); errors.As(err, &c) {
		ctx["debug"] = c.Debug()
	}

	return l.WithField("error", ctx)
}
