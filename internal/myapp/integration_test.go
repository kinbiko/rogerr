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
	
	// Use ErrorHandler to extract stacktrace
	handler := rogerr.NewErrorHandler()
	frames := handler.Stacktrace(err)
	if len(frames) == 0 {
		t.Fatal("expected stacktrace frames")
	}
	
	// Log all frames for debugging
	t.Log("Captured stacktrace frames:")
	for i, frame := range frames {
		t.Logf("  [%d] %s (%s:%d) InApp=%v", i, frame.Function, frame.File, frame.Line, frame.InApp)
	}
	
	// Define expected frames in reverse order (closest to error first)
	expectedFrames := []struct {
		functionPattern string
		filePattern     string
		inApp           bool
		description     string
	}{
		{
			functionPattern: "executeBusinessLogic.func1",
			filePattern:     "/internal/myapp/pkg/service/processing.go",
			inApp:           true,
			description:     "anonymous function in service package",
		},
		{
			functionPattern: "executeBusinessLogic",
			filePattern:     "/internal/myapp/pkg/service/processing.go",
			inApp:           true,
			description:     "method on ProcessingService",
		},
		{
			functionPattern: "ProcessData",
			filePattern:     "/internal/myapp/pkg/service/processing.go",
			inApp:           true,
			description:     "method on ProcessingService",
		},
		{
			functionPattern: "processRequest.func1",
			filePattern:     "/internal/myapp/pkg/handler/request.go",
			inApp:           true,
			description:     "anonymous function in handler package",
		},
		{
			functionPattern: "processRequest",
			filePattern:     "/internal/myapp/pkg/handler/request.go",
			inApp:           true,
			description:     "method on RequestHandler",
		},
		{
			functionPattern: "HandleRequest",
			filePattern:     "/internal/myapp/pkg/handler/request.go",
			inApp:           true,
			description:     "method on RequestHandler",
		},
		{
			functionPattern: "Execute",
			filePattern:     "/internal/myapp/cmd/app.go",
			inApp:           true,
			description:     "method on App struct",
		},
		{
			functionPattern: "run",
			filePattern:     "/internal/myapp/main.go",
			inApp:           true,
			description:     "main package function",
		},
	}
	
	// Validate that myapp frames are marked as InApp=true
	myappFrameCount := 0
	for _, frame := range frames {
		if strings.Contains(frame.File, "/internal/myapp/") {
			myappFrameCount++
			if !frame.InApp {
				t.Errorf("Frame from myapp should be InApp=true: %s (%s:%d)", 
					frame.Function, frame.File, frame.Line)
			}
		}
	}
	
	if myappFrameCount == 0 {
		t.Error("expected to find frames from myapp module")
	}
	
	// Validate that mylib frames are marked as InApp=false
	mylibFrameCount := 0
	for _, frame := range frames {
		if strings.Contains(frame.File, "/internal/mylib/") {
			mylibFrameCount++
			if frame.InApp {
				t.Errorf("Frame from mylib should be InApp=false: %s (%s:%d)", 
					frame.Function, frame.File, frame.Line)
			}
		}
	}
	
	if mylibFrameCount == 0 {
		t.Error("expected to find frames from mylib module")
	}
	
	// Validate specific function patterns exist (only for myapp frames)
	foundFunctions := make(map[string]bool)
	for _, frame := range frames {
		// Skip testing framework frames and mylib frames for our validation
		if strings.Contains(frame.Function, "testing.") || 
		   strings.Contains(frame.Function, "runtime.") ||
		   strings.Contains(frame.Function, "github.com/kinbiko/rogerr/internal/mylib") {
			continue
		}
		
		for _, expected := range expectedFrames {
			if strings.Contains(frame.Function, expected.functionPattern) {
				foundFunctions[expected.functionPattern] = true
				
				// Validate file path
				if !strings.Contains(frame.File, expected.filePattern) {
					t.Errorf("Frame %s should be in file containing %s, got %s", 
						frame.Function, expected.filePattern, frame.File)
				}
				
				// Validate InApp status
				if frame.InApp != expected.inApp {
					t.Errorf("Frame %s should have InApp=%v, got %v", 
						frame.Function, expected.inApp, frame.InApp)
				}
				
				// Validate line number is reasonable
				if frame.Line <= 0 {
					t.Errorf("Frame %s should have positive line number, got %d", 
						frame.Function, frame.Line)
				}
			}
		}
	}
	
	// Check that we found all expected function patterns
	for _, expected := range expectedFrames {
		if !foundFunctions[expected.functionPattern] {
			t.Errorf("Expected to find function containing '%s' (%s) in stacktrace", 
				expected.functionPattern, expected.description)
		}
	}
	
	// Validate that no rogerr main package frames are present
	for _, frame := range frames {
		if strings.HasPrefix(frame.Function, "github.com/kinbiko/rogerr.") {
			t.Errorf("Found rogerr main package frame in stacktrace: %s", frame.Function)
		}
	}
}