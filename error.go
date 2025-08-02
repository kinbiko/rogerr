package rogerr

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
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

// getModulePath returns the main module path for determining if frames are in-app.
func getModulePath() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		return bi.Main.Path
	}
	return ""
}

// captureStacktrace captures the current call stack, excluding rogerr internal frames.
func captureStacktrace(modulePath string) []Frame {
	const maxFrames = 64
	ptrs := [maxFrames]uintptr{}

	// Skip 0 frames as we'll filter manually
	pcs := ptrs[0:runtime.Callers(0, ptrs[:])]

	allFrames := make([]Frame, 0, len(pcs))
	iter := runtime.CallersFrames(pcs)

	// First, collect all frames
	for {
		frame, more := iter.Next()

		// Determine if this is application code
		inApp := isInApp(frame.Function, modulePath)

		allFrames = append(allFrames, Frame{
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
			InApp:    inApp,
		})

		if !more {
			break
		}
	}

	// Now filter out rogerr frames, but keep everything after the last rogerr frame
	lastRogerrIndex := -1
	for i, frame := range allFrames {
		// Only filter out the main rogerr package, not internal modules
		if strings.HasPrefix(frame.Function, "github.com/kinbiko/rogerr.") {
			lastRogerrIndex = i
		}
	}

	// Return frames after the last rogerr frame
	if lastRogerrIndex >= 0 && lastRogerrIndex+1 < len(allFrames) {
		return allFrames[lastRogerrIndex+1:]
	}

	// If no rogerr frames found, return all frames (shouldn't happen)
	return allFrames
}

// isInApp determines if a function belongs to the application or a dependency.
func isInApp(function, modulePath string) bool {
	if modulePath == "" {
		return false
	}

	// Special case for main.main
	if strings.Contains(function, "main.main") {
		return true
	}

	// Check if function belongs to the main module
	return strings.HasPrefix(function, modulePath)
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

	// Capture stacktrace if enabled
	if h.stacktrace {
		e.stacktrace = captureStacktrace(getModulePath())
	}

	return e
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
