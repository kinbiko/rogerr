package rogerr

import (
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

// getModulePath returns the application module path for determining if frames are in-app.
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

	for {
		frame, more := iter.Next()
		allFrames = append(allFrames, Frame{
			File:     frame.File,
			Line:     frame.Line,
			Function: frame.Function,
			InApp:    isInApp(frame.Function, modulePath), // Determine if this is application code
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

	// Handle case where the binary is built from a module and function names
	// start with "main." - check if the module path contains the current module
	if strings.HasPrefix(function, "main.") {
		return true
	}

	// Check if function belongs to the app module
	return strings.HasPrefix(function, modulePath)
}
