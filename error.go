package rogerr

import (
	"context"
	"fmt"
)

// Frame represents a single frame in a stacktrace.
type Frame struct {
	File     string // Full file path
	Line     int    // Line number
	Function string // Function or method name
	InApp    bool   // true if application code, false if dependency
}

// ErrorHandler provides configurable error handling with optional stacktrace capture.
type ErrorHandler struct {
	stacktrace bool
}

// Option is a function that configures an ErrorHandler.
type Option func(*ErrorHandler)

type rError struct {
	err error
	ctx context.Context
	msg string
}

// Error returns the message of the rError, along with any wrapped error messages.
func (e *rError) Error() string {
	if e.err == nil && e.msg == "" {
		return "unknown error"
	}
	if e.err == nil {
		return e.msg
	}
	if e.msg == "" {
		return e.err.Error()
	}
	return fmt.Sprintf("%s: %s", e.msg, e.err)
}

// Unwrap is the conventional method for getting the underlying error of an error.
func (e *rError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

// WithStacktrace configures whether stacktraces should be captured.
func WithStacktrace(enabled bool) Option {
	return func(h *ErrorHandler) {
		h.stacktrace = enabled
	}
}

// NewErrorHandler creates a new ErrorHandler with the given options.
// By default, stacktrace capture is enabled.
func NewErrorHandler(opts ...Option) *ErrorHandler {
	h := &ErrorHandler{
		stacktrace: true, // stacktrace enabled by default
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Wrap attaches ctx data and wraps the given error with message.
// ctx, err, and msgAndFmtArgs are all optional, but at least one must be given
// for this function to return a non-nil error.
// Any attached diagnostic data from this ctx will be preserved should you
// pass the returned error further up the stack.
// Deprecated: Use ErrorHandler.Wrap instead.
func Wrap(ctx context.Context, err error, msgAndFmtArgs ...interface{}) error {
	if ctx == nil && err == nil && msgAndFmtArgs == nil {
		return nil
	}
	e := &rError{err: err, ctx: ctx}

	if l := len(msgAndFmtArgs); l > 0 {
		if msg, ok := msgAndFmtArgs[0].(string); ok {
			e.msg = fmt.Sprintf(msg, msgAndFmtArgs[1:]...)
		}
	}
	return e
}
