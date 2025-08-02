package rogerr

import (
	"context"
	"strings"
	"testing"
)

func TestStacktraceCapture(t *testing.T) {
	func1 := func() error { return NewErrorHandler().Wrap(context.Background(), nil, "test error") }
	func2 := func() error { return func1() }

	err := func2()
	if err == nil {
		t.Fatal("expected error")
	}

	for _, frame := range err.(*rError).stacktrace {
		if frame.File == "" {
			t.Error("expected file to be populated")
		}
		if frame.Line == 0 {
			t.Error("expected line number to be populated")
		}
		if frame.Function == "" {
			t.Error("expected function name to be populated")
		}
	}
}

func TestFrameFiltering(t *testing.T) {
	err := NewErrorHandler().Wrap(context.Background(), nil, "test error")
	if err == nil {
		t.Fatal("expected error")
	}

	for _, frame := range err.(*rError).stacktrace {
		if strings.Contains(frame.Function, "github.com/kinbiko/rogerr") {
			t.Errorf("found rogerr internal frame: %s", frame.Function)
		}
	}
}

func TestIsInApp(t *testing.T) {
	modulePath := "github.com/example/myapp"
	for name, tc := range map[string]struct {
		function string
		expected bool
	}{
		"main function":         {"main.main", true},
		"main package function": {"main.someFunction", true},
		"app module function":   {"github.com/example/myapp/pkg.DoSomething", true},
		"app module method":     {"github.com/example/myapp/internal/service.(*Service).Process", true},
		"external dependency":   {"github.com/external/lib.Function", false},
		"stdlib function":       {"fmt.Printf", false},
	} {
		t.Run(name, func(t *testing.T) {
			if result := isInApp(tc.function, modulePath); result != tc.expected {
				t.Errorf("isInApp(%q, %q) = %v, expected %v", tc.function, modulePath, result, tc.expected)
			}
		})
	}

	t.Run("empty module path", func(t *testing.T) {
		if isInApp("main.main", "") {
			t.Errorf(`isInApp("main.main", "") = true, expected false`)
		}
	})
}
