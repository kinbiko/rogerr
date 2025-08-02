package rogerr

import (
	"context"
	"errors"
	"testing"
)

/*
Expected (and tested) behaviour:
| ctx | err | msg | args | test                                                                  |
|-----|-----|-----|------|-----------------------------------------------------------------------|
|     |     |     |      | return nil                                                            |
|  o  |     |     |      | unknown error & attach ctx                                            |
|     |  o  |     |      | wrap error with no further message                                    |
|  o  |  o  |     |      | wrap error with no further message & attach ctx                       |
|     |     |  o  |      | new error with the given message                                      |
|  o  |     |  o  |      | new error with the given message & attach ctx                         |
|     |  o  |  o  |      | wrap error with the given message                                     |
|  o  |  o  |  o  |      | wrap error with the given message & attach ctx                        |
|     |     |     |  o   | attempt to fill in an error based on the args as message              |
|  o  |     |     |  o   | attempt to fill in an error based on the args as message & attach ctx |
|     |  o  |     |  o   | attempt to wrap the error based on the args as message                |
|  o  |  o  |     |  o   | attempt to wrap the error based on the args as message & attach ctx   |
|     |     |  o  |  o   | new error with the given message and args                             |
|  o  |     |  o  |  o   | new error with the given message and args & attach ctx                |
|     |  o  |  o  |  o   | wrap error with message + args                                        |
|  o  |  o  |  o  |  o   | wrap error with message + args & attach ctx                           |
*/
// Tests that the way Wrap works is relatively intuitive.
func TestWrap(t *testing.T) {
	t.Run("no real input", func(t *testing.T) {
		if wrappedErr := Wrap(nil, nil); wrappedErr != nil { //nolint:staticcheck // Testing that we don't do a dumb when users do a dumb
			t.Errorf("expected no error returned but got: %s", wrappedErr)
		}
	})

	var (
		ctx           = context.Background()
		err           = errors.New("something bad happened")
		msg           = "unable to foo the bar"
		args          = []interface{}{1, "arg"}
		wrappedMsg    = msg + ": " + err.Error()
		fmtBase       = "I got %d and %s"
		fmtMsg        = "I got 1 and arg"
		wrappedFmtMsg = fmtMsg + ": " + err.Error()

		// Note: consider using a custom format instead of this format, as this
		// means we're relying on Go's representation, and if they change
		// across versions this could get awkward
		fmtFluff       = "%!(EXTRA int=1, string=arg)"
		errAndFmtFluff = fmtFluff + ": " + err.Error()
	)

	for _, tc := range []struct {
		name         string
		ctx          context.Context
		err          error
		msg          string
		args         []interface{}
		expErrString string
	}{
		{ctx: ctx, err: nil, msg: "", args: nil, expErrString: "unknown error", name: "only ctx"},
		{ctx: nil, err: err, msg: "", args: nil, expErrString: err.Error(), name: "only err"},
		{ctx: ctx, err: err, msg: "", args: nil, expErrString: err.Error(), name: "ctx and err"},
		{ctx: nil, err: nil, msg: msg, args: nil, expErrString: msg, name: "only msg"},
		{ctx: ctx, err: nil, msg: msg, args: nil, expErrString: msg, name: "msg and ctx"},
		{ctx: nil, err: err, msg: msg, args: nil, expErrString: wrappedMsg, name: "msg and err"},
		{ctx: ctx, err: err, msg: msg, args: nil, expErrString: wrappedMsg, name: "msg, err and ctx"},
		{ctx: nil, err: nil, msg: "", args: args, expErrString: fmtFluff, name: "just args"},
		{ctx: ctx, err: nil, msg: "", args: args, expErrString: fmtFluff, name: "args and ctx"},
		{ctx: nil, err: err, msg: "", args: args, expErrString: errAndFmtFluff, name: "args and err"},
		{ctx: ctx, err: err, msg: "", args: args, expErrString: errAndFmtFluff, name: "args, err, and ctx"},
		{ctx: nil, err: nil, msg: fmtBase, args: args, expErrString: fmtMsg, name: "message and args only"},
		{ctx: ctx, err: nil, msg: fmtBase, args: args, expErrString: fmtMsg, name: "message, args and ctx"},
		{ctx: nil, err: err, msg: fmtBase, args: args, expErrString: wrappedFmtMsg, name: "all args but ctx"},
		{ctx: ctx, err: err, msg: fmtBase, args: args, expErrString: wrappedFmtMsg, name: "all args"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.msg != "" || tc.args != nil {
				tc.args = append([]interface{}{tc.msg}, tc.args...)
			}
			wrapped := Wrap(tc.ctx, tc.err, tc.args...).(*rError)
			if wrapped == nil {
				t.Fatalf("got nil err")
			}
			if exp, got := tc.ctx, wrapped.ctx; got != exp {
				t.Errorf("expected ctx to be identical but differed:\nexp: %v\ngot: %v", exp, got)
			}
			if exp, got := tc.expErrString, wrapped.Error(); exp != got {
				t.Errorf("unexpected error message,\nexp: %s\ngot: %s", exp, got)
			}
		})
	}
}

func TestNewErrorHandler(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		handler := NewErrorHandler()
		if !handler.stacktrace {
			t.Error("expected stacktrace to be enabled by default")
		}
	})

	t.Run("with stacktrace enabled", func(t *testing.T) {
		handler := NewErrorHandler(WithStacktrace(true))
		if !handler.stacktrace {
			t.Error("expected stacktrace to be enabled")
		}
	})

	t.Run("with stacktrace disabled", func(t *testing.T) {
		handler := NewErrorHandler(WithStacktrace(false))
		if handler.stacktrace {
			t.Error("expected stacktrace to be disabled")
		}
	})
}
