package errorsx

import "github.com/pkg/errors"

// WithStack mirror pkg/errors.WithStack but does not wrap existing stack
// traces.
// Deprecated: you should probably use errors.WithStack instead and only annotate stacks when it makes sense.
func WithStack(err error) error {
	if e, ok := err.(StackTracer); ok && len(e.StackTrace()) > 0 {
		return err
	}

	return errors.WithStack(err)
}

// StatusCodeCarrier can be implemented by an error to support setting status codes in the error itself.
type StatusCodeCarrier interface {
	// StatusCode returns the status code of this error.
	StatusCode() int
}

// RequestIDCarrier can be implemented by an error to support error contexts.
type RequestIDCarrier interface {
	// RequestID returns the ID of the request that caused the error, if applicable.
	RequestID() string
}

// ReasonCarrier can be implemented by an error to support error contexts.
type ReasonCarrier interface {
	// Reason returns the reason for the error, if applicable.
	Reason() string
}

// DebugCarrier can be implemented by an error to support error contexts.
type DebugCarrier interface {
	// Debug returns debugging information for the error, if applicable.
	Debug() string
}

// StatusCarrier can be implemented by an error to support error contexts.
type StatusCarrier interface {
	// ID returns the error id, if applicable.
	Status() string
}

// DetailsCarrier can be implemented by an error to support error contexts.
type DetailsCarrier interface {
	// Details returns details on the error, if applicable.
	Details() map[string]interface{}
}

// IDCarrier can be implemented by an error to support error contexts.
type IDCarrier interface {
	// ID returns application error ID on the error, if applicable.
	ID() string
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}
