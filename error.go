package rogerr

import (
	"context"
	"fmt"
)

type rError struct {
	err        error
	ctx        context.Context
	msg        string
	stacktrace []Frame
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

// Wrap wraps errors with the default error handler settings.
// See ErrorHandler.Wrap for more details.
// Deprecated: Use ErrorHandler.Wrap instead.
func Wrap(ctx context.Context, err error, msgAndFmtArgs ...any) error {
	return NewErrorHandler().Wrap(ctx, err, msgAndFmtArgs...)
}
