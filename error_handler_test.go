package rogerr

import (
	"errors"
	"testing"
)

func TestNewErrorHandler(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		if !NewErrorHandler().stacktrace {
			t.Error("expected stacktrace to be enabled by default")
		}
	})

	t.Run("with stacktrace enabled", func(t *testing.T) {
		if !NewErrorHandler(WithStacktrace(true)).stacktrace {
			t.Error("expected stacktrace to be enabled")
		}
	})

	t.Run("with stacktrace disabled", func(t *testing.T) {
		if NewErrorHandler(WithStacktrace(false)).stacktrace {
			t.Error("expected stacktrace to be disabled")
		}
	})
}

func TestErrorHandlerWrap(t *testing.T) {
	handler := NewErrorHandler()
	ctx := t.Context()

	t.Run("basic error wrapping with stacktrace enabled", func(t *testing.T) {
		baseErr := errors.New("base error")

		err := handler.Wrap(ctx, baseErr, "wrapped error")
		if err == nil {
			t.Fatal("expected non-nil error")
		}

		rErr := err.(*rError)
		if rErr.err != baseErr {
			t.Error("expected wrapped error to contain base error")
		}
		if rErr.ctx != ctx {
			t.Error("expected wrapped error to contain context")
		}
		if rErr.msg != "wrapped error" {
			t.Errorf("expected message 'wrapped error', got '%s'", rErr.msg)
		}
		if len(rErr.stacktrace) == 0 {
			t.Error("expected stacktrace to be captured when enabled")
		}
	})

	t.Run("error wrapping with stacktrace disabled", func(t *testing.T) {
		baseErr := errors.New("base error")
		err := NewErrorHandler(WithStacktrace(false)).Wrap(ctx, baseErr, "wrapped error")
		if err == nil {
			t.Fatal("expected non-nil error")
		}

		rErr := err.(*rError)
		if len(rErr.stacktrace) != 0 {
			t.Error("expected no stacktrace when disabled")
		}
	})

	t.Run("nil inputs return nil", func(t *testing.T) {
		err := handler.Wrap(t.Context(), nil)
		if err == nil {
			t.Error("expected non-nil error when context is provided")
		}
	})

	t.Run("message formatting works", func(t *testing.T) {
		err := handler.Wrap(t.Context(), nil, "user %d failed: %s", 123, "timeout")

		if err == nil {
			t.Fatal("expected non-nil error")
		}

		expected := "user 123 failed: timeout"
		if err.Error() != expected {
			t.Errorf("expected '%s', got '%s'", expected, err.Error())
		}
	})
}

func TestErrorHandlerStacktrace(t *testing.T) {
	handler := NewErrorHandler()
	t.Run("extract stacktrace from error", func(t *testing.T) {
		if len(handler.Stacktrace(handler.Wrap(t.Context(), nil, "test error"))) == 0 {
			t.Error("expected stacktrace frames to be extracted")
		}
	})

	t.Run("return nil for non-rogerr error", func(t *testing.T) {
		if handler.Stacktrace(errors.New("regular error")) != nil {
			t.Error("expected nil stacktrace for non-rogerr error")
		}
	})

	t.Run("return nil for nil error", func(t *testing.T) {
		if handler.Stacktrace(nil) != nil {
			t.Error("expected nil stacktrace for nil error")
		}
	})
}
