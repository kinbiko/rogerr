package rogerr

import (
	"context"
	"fmt"
)

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

// Wrap attaches ctx data and wraps the given error with message.
// ctx, err, and msgAndFmtArgs are all optional, but at least one must be given
// for this function to return a non-nil error.
// Any attached diagnostic data from this ctx will be preserved should you
// pass the returned error further up the stack.
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
