package rogerr

import (
	"context"
	"errors"
	"fmt"
)

// ErrorHandler provides configurable error handling with optional stacktrace capture.
type ErrorHandler struct {
	stacktrace bool
}

// Option is a function that configures an ErrorHandler.
type Option func(*ErrorHandler)

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

// Stacktrace extracts the stacktrace from an error if it was created with ErrorHandler.
func (h *ErrorHandler) Stacktrace(err error) []Frame {
	rErr := &rError{}
	if errors.As(err, &rErr) {
		return rErr.stacktrace
	}
	return nil
}

// Wrap attaches ctx data and wraps the given error with message, optionally capturing stacktrace.
// ctx, err, and msgAndFmtArgs are all optional, but at least one must be given
// for this function to return a non-nil error.
// Any attached diagnostic data from this ctx will be preserved should you
// pass the returned error further up the stack.
func (h *ErrorHandler) Wrap(ctx context.Context, err error, msgAndFmtArgs ...interface{}) error {
	if ctx == nil && err == nil && msgAndFmtArgs == nil {
		return nil
	}
	e := &rError{err: err, ctx: ctx}

	if l := len(msgAndFmtArgs); l > 0 {
		if msg, ok := msgAndFmtArgs[0].(string); ok {
			e.msg = fmt.Sprintf(msg, msgAndFmtArgs[1:]...)
		}
	}
	if h.stacktrace {
		e.stacktrace = captureStacktrace(getModulePath())
	}

	return e
}
