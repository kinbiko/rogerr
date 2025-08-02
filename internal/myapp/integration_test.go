package main

import (
	"strings"
	"testing"

	"github.com/kinbiko/rogerr"
)

func TestStacktraceIntegration(t *testing.T) {
	err := run([]string{"testdata"})
	if err == nil {
		t.Fatal("expected error from run")
	}

	frames := rogerr.NewErrorHandler().Stacktrace(err)
	var myappFrames, mylibFrames, rogerrFrames int

	for _, frame := range frames {
		if frame.File == "" || frame.Function == "" || frame.Line <= 0 {
			t.Errorf("Invalid frame: %+v", frame)
		}

		switch {
		case strings.Contains(frame.File, "/internal/myapp/"):
			myappFrames++
			if !frame.InApp {
				t.Errorf("myapp frame should be InApp=true: %s", frame.Function)
			}
		case strings.Contains(frame.File, "/internal/mylib/"):
			mylibFrames++
			if frame.InApp {
				t.Errorf("mylib frame should be InApp=false: %s", frame.Function)
			}
		case strings.HasPrefix(frame.Function, "github.com/kinbiko/rogerr."):
			rogerrFrames++
			t.Errorf("Found rogerr main package frame: %s", frame.Function)
		}
	}

	if myappFrames != 9 {
		t.Errorf("expected 9 frames from myapp module, got %d", myappFrames)
	}
	if mylibFrames != 5 {
		t.Errorf("expected 5 frames from mylib module, got %d", mylibFrames)
	}
	if rogerrFrames != 0 {
		t.Errorf("expected 0 frames from rogerr module, got %d", rogerrFrames)
	}
}
