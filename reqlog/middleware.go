package reqlog

import (
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/huanggze/x/logrusx"
)

type timer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

// Middleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type Middleware struct {
	// Logger is the log.Logger instance used to log messages with the Logger middleware
	Logger *logrusx.Logger
	// Name is the name of the application as recorded in latency metrics
	Name   string
	Before func(*logrusx.Logger, *http.Request, string) *logrusx.Logger
	After  func(*logrusx.Logger, *http.Request, negroni.ResponseWriter, time.Duration, string) *logrusx.Logger

	logStarting bool

	clock timer

	logLevel logrus.Level

	// Silence log for specific URL paths
	silencePaths map[string]bool

	sync.RWMutex
}

type realClock struct{}

func (rc *realClock) Now() time.Time {
	return time.Now()
}

func (rc *realClock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// NewMiddlewareFromLogger returns a new *Middleware which writes to a given logrus logger.
func NewMiddlewareFromLogger(logger *logrusx.Logger, name string) *Middleware {
	return &Middleware{
		Logger: logger,
		Name:   name,
		Before: DefaultBefore,
		After:  DefaultAfter,

		logLevel:     logrus.InfoLevel,
		logStarting:  true,
		clock:        &realClock{},
		silencePaths: map[string]bool{},
	}
}

// DefaultBefore is the default func assigned to *Middleware.Before
func DefaultBefore(entry *logrusx.Logger, req *http.Request, remoteAddr string) *logrusx.Logger {
	return entry.WithRequest(req)
}

// DefaultAfter is the default func assigned to *Middleware.After
func DefaultAfter(entry *logrusx.Logger, req *http.Request, res negroni.ResponseWriter, latency time.Duration, name string) *logrusx.Logger {
	e := entry.WithRequest(req).WithField("http_response", map[string]any{
		"status":      res.Status(),
		"size":        res.Size(),
		"text_status": http.StatusText(res.Status()),
		"took":        latency,
		"headers":     entry.HTTPHeadersRedacted(res.Header()),
	})
	if el := totalExternalLatency(req.Context()); el > 0 {
		e = e.WithFields(map[string]any{
			"took_internal": latency - el,
			"took_external": el,
		})
	}
	return e
}
