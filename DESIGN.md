# Stacktrace Feature Design Document

## Overview

This document outlines the design for adding stacktrace functionality to the rogerr package through a new object-oriented API. The feature will capture the call stack at the point where errors are wrapped, providing valuable debugging information while maintaining the package's zero-dependency philosophy and backward compatibility.

## Goals

1. **Object-oriented design**: Introduce an `ErrorHandler` type with configurable behavior
2. **Capture meaningful stacktraces**: Record the call stack when errors are wrapped, excluding internal rogerr frames
3. **Distinguish application vs dependency code**: Mark each frame as belonging to the application or external dependencies
4. **Maintain zero dependencies**: Implement using only Go's standard library
5. **Preserve backward compatibility**: Keep existing package-level functions but deprecate in favor of new API
6. **Flexible configuration**: Use options pattern for ErrorHandler configuration

## Design Decisions

### Core Data Structures

#### Frame Structure

```go
type Frame struct {
    File     string // Full file path
    Line     int    // Line number
    Function string // Function or method name
    InApp    bool   // true if application code, false if dependency
}
```

#### ErrorHandler Type

```go
type ErrorHandler struct {
    skipStacktrace bool
    // Future configuration options can be added here
}

type Option func(*ErrorHandler)
```

#### Enhanced rError Structure

```go
type rError struct {
    err        error
    ctx        context.Context
    msg        string
    stacktrace []Frame // New field for stacktrace (nil when skipStacktrace is true)
}
```

### API Design

#### ErrorHandler Construction (Options Pattern)

```go
func WithStacktrace(b bool) Option {
    return func(h *ErrorHandler) {
        h.stacktrace = b
    }
}

// Constructor with reasonable defaults.
func NewErrorHandler(opts ...Option) *ErrorHandler {
    h := &ErrorHandler{ stacktrace: true } // stacktrace enabled by default
    for _, opt := range opts {
        opt(h)
    }
    return h
}

// Examples:
// handler := rogerr.NewErrorHandler() // stacktrace enabled
// handler := rogerr.NewErrorHandler(rogerr.WithStacktrace(false)) // disabled
```

#### ErrorHandler Methods

```go
// Primary error wrapping method
func (h *ErrorHandler) Wrap(ctx context.Context, err error, msgAndFmtArgs ...interface{}) error

// Metadata extraction (also available at package level for backward compatibility)
func (h *ErrorHandler) Metadata(err error) map[string]interface{}

// Stacktrace extraction
func (h *ErrorHandler) Stacktrace(err error) []Frame
```

#### Backward Compatibility

```go
// Deprecated: Use ErrorHandler.Wrap instead
func Wrap(ctx context.Context, err error, msgAndFmtArgs ...interface{}) error {
    // Implementation delegates to a default ErrorHandler with stacktrace enabled
}

// Keep existing package-level functions
func Metadata(err error) map[string]interface{} // unchanged
func WithMetadata(ctx context.Context, data map[string]interface{}) context.Context // unchanged
func WithMetadatum(ctx context.Context, key string, value interface{}) context.Context // unchanged
```

### Implementation Strategy

For reference, here's the implementation for extracting stackframes from a different package (github.com/kinbiko/bugsnag), that might be helpful (or it migh be outdated. Use your judgment):

```go
func makeStacktrace(module string) []*JSONStackframe {
	ptrs := [50]uintptr{}
	// Skip 0 frames as we strip this manually later by ignoring any frames
	// including github.com/kinbiko/bugsnag (or below).
	pcs := ptrs[0:runtime.Callers(0, ptrs[:])]

	stacktrace := make([]*JSONStackframe, len(pcs))
	for i, pc := range pcs { //nolint:varnamelen // indexes are conventionally i
		pc-- // pc - 1 is the *real* program counter, for reasons beyond me.

		file, lineNumber, method := "unknown", 0, "unknown"
		if fn := runtime.FuncForPC(pc); fn != nil {
			file, lineNumber = fn.FileLine(pc)
			method = fn.Name()
		}
		inProject := module != "" && strings.Contains(method, module) || strings.Contains(method, "main.main")
		if inProject {
			file = calculateSourcepathHeuristic(file)
		}

		stacktrace[i] = &JSONStackframe{File: file, LineNumber: lineNumber, Method: method, InProject: inProject}
	}

	// Drop any frames from this package, and further down, for example Go
	// stdlib packages. Rather than trying to guess how many frames to skip,
	// this approach will work better on multiple platforms
	lastBugsnagIndex := 0
	for i, sf := range stacktrace {
		if strings.Contains(sf.Method, "github.com/kinbiko/bugsnag.") {
			lastBugsnagIndex = i
		}
	}
	return stacktrace[lastBugsnagIndex+1:]
}

// This function attempst to rewrite the filepath value to be relative to the
// root of the repository. This allows correct filepaths in the Bugsnag
// dashboard.
// For now, this is limited to in-project files hosted on GitHub.
// There will be false positives, for most Go repos hosted on GitHub this
// should work out of the box.
func calculateSourcepathHeuristic(file string) string {
	if strings.HasPrefix(file, "github.com") {
		// Split
		// "github.com/kinbiko/bugsnag/examples/cmd/cli/main.go"
		// into
		// [
		//   "github.com",
		//   "kinbiko",
		//   "bugsnag",
		//   "examples/cmd/cli/main.go",
		// ]
		numSplits := 4
		split := strings.SplitN(file, "/", numSplits)
		if len(split) == numSplits {
			return split[numSplits-1]
		}
	}
	return file
}

// makeModulePath defines the root of the project that uses this package.
// Used to identify if a file is "in-project" or a third party library,
// which is in turn used by Bugsnag to group errors by the top stackframe
// that's "in project".
func makeModulePath() string {
	if bi, ok := debug.ReadBuildInfo(); ok {
		return bi.Main.Path
	}
	return ""
}
```

#### 1. Stacktrace Capture

- Use `runtime.Callers()` to capture program counters starting from the caller of `Wrap`
- Use `runtime.CallersFrames()` to convert PCs to frame information
- Filter out rogerr internal frames by checking package paths
- Skip capture entirely when `skipStacktrace` is true for performance

#### 2. Application vs Dependency Detection

**Module-based approach**:

- Use `debug.ReadBuildInfo()` to get the main module path
- Compare each frame's package path against the main module path
- Mark as `InApp: true` if the package path starts with the main module path
- Special handling for:
  - Standard library packages (mark as dependency)
  - Test packages (mark as application if within main module)
  - Vendor directories (mark as dependency)

#### 3. Frame Filtering

Exclude frames from:

- The rogerr package itself (`github.com/kinbiko/rogerr`)
- Runtime package internals that aren't useful for debugging
- Any frame before the actual caller of `ErrorHandler.Wrap`

#### 4. Performance Considerations

- When `skipStacktrace` is true, avoid any stacktrace-related work entirely
- Limit stacktrace depth to prevent excessive memory usage (e.g., 64 frames max)
- Store frames immediately rather than lazy evaluation to avoid GC pressure

### Integration Points

#### Error Interface Compliance

- The enhanced `rError` maintains full compatibility with Go's error interface
- `Error()`, `Unwrap()`, and `errors.As()` behavior remain unchanged
- Stacktrace information doesn't affect error equality or wrapping behavior

#### Default ErrorHandler

- Package-level functions will use a singleton default ErrorHandler with stacktrace enabled
- This ensures backward compatibility while providing the new functionality

### Example Usage

#### New Object-Oriented API

```go
// Create handler with stacktrace enabled (default)
handler := rogerr.NewErrorHandler()

// Create handler with stacktrace disabled
handlerNoStack := rogerr.NewErrorHandler(rogerr.WithStacktrace(false))

// Use handler to wrap errors
ctx = rogerr.WithMetadatum(ctx, "user_id", 12345)
err = handler.Wrap(ctx, err, "failed to process user request")

// Extract metadata (available on both handler and package level)
metadata := handler.Metadata(err)
// OR
metadata = rogerr.Metadata(err) // backward compatibility

// Extract stacktrace (only available on handler)
if frames := handler.Stacktrace(err); len(frames) > 0 {
    for _, frame := range frames {
        log.Printf("%s:%d in %s (app: %v)",
            frame.File, frame.Line, frame.Function, frame.InApp)
    }
}
```

#### Backward Compatible Usage

```go
// Existing code continues to work unchanged
ctx = rogerr.WithMetadatum(ctx, "user_id", 12345)
err = rogerr.Wrap(ctx, err, "failed to process user request") // deprecated but functional

// Existing metadata extraction unchanged
metadata := rogerr.Metadata(err)
```

## Implementation Plan

### Phase 1: Core Structure

1. Define `Frame` struct
2. Create `ErrorHandler` type with options pattern
3. Implement `NewErrorHandler` constructor
4. Add `skipStacktrace` option

### Phase 2: Stacktrace Capture

1. Implement stacktrace capture logic using `runtime` package
2. Add frame filtering to exclude rogerr internals
3. Integrate stacktrace into `rError` struct
4. Implement `ErrorHandler.Wrap` method

### Phase 3: Application Detection

1. Implement module-based detection using `debug.ReadBuildInfo()`
2. Add logic to classify frames as application vs dependency
3. Handle edge cases (standard library, vendor, etc.)

### Phase 4: API Integration

1. Implement `ErrorHandler.Stacktrace` method
2. Implement `ErrorHandler.Metadata` method
3. Update package-level functions to use default ErrorHandler
4. Add deprecation notices to existing functions

### Phase 5: Testing & Documentation

1. Comprehensive unit tests for all new functionality
2. Integration tests with realistic call stacks
3. Performance benchmarks comparing with/without stacktrace
4. Update documentation and examples

## Backward Compatibility

This design maintains complete backward compatibility:

- All existing package-level functions remain unchanged in behavior
- Existing error types and interfaces are unmodified
- No breaking changes to public API
- Migration path is opt-in through new ErrorHandler API

The only change is deprecation warnings on existing functions, encouraging migration to the new object-oriented API.

## Performance Impact

- **With stacktrace disabled**: Zero performance impact through `skipStacktrace` option
- **With stacktrace enabled**: Minimal overhead only when errors are actually wrapped
- **Memory usage**: Additional memory only for Frame slices when stacktraces are captured
- **Backward compatibility**: Existing code gets stacktrace functionality automatically but can opt out if needed

## Documentation Updates Required

### Package Documentation (doc.go)
- Update package-level documentation to introduce the new `ErrorHandler` API
- Add examples showing both old and new usage patterns
- Document the `WithStacktrace` option and its default behavior
- Include migration guidance from package-level functions to `ErrorHandler`

### Function Documentation
- Add comprehensive documentation for `ErrorHandler.Stacktrace()` method:
  - Explain what stacktraces are captured (call site of `Wrap`)
  - Document the `Frame` struct fields and their meanings
  - Clarify the `InApp` field logic for application vs dependency detection
  - Provide examples of iterating through frames
  - Document behavior when no stacktrace is available (returns empty slice)

### Structured Logging Integration

#### log/slog Integration
Document how to use rogerr with Go's structured logging:

```go
// Basic metadata logging
metadata := handler.Metadata(err)
slog.Error("request failed", 
    "error", err.Error(),
    "metadata", metadata)

// Stacktrace logging for debugging
frames := handler.Stacktrace(err)
if len(frames) > 0 {
    slog.Error("request failed",
        "error", err.Error(),
        "metadata", metadata,
        "stacktrace", frames)
}
```

#### Datadog-Compatible Logging
For logging engines like Datadog that can index stacktrace data:

```go
// Format stacktrace for Datadog ingestion
func formatStacktraceForDatadog(frames []rogerr.Frame) []map[string]interface{} {
    result := make([]map[string]interface{}, len(frames))
    for i, frame := range frames {
        result[i] = map[string]interface{}{
            "filename":    frame.File,
            "lineno":      frame.Line,
            "function":    frame.Function,
            "in_app":      frame.InApp,
            "abs_path":    frame.File,
        }
    }
    return result
}

// Usage with slog for Datadog
metadata := handler.Metadata(err)
frames := handler.Stacktrace(err)

slog.Error("application error",
    slog.String("error.message", err.Error()),
    slog.Any("error.metadata", metadata),
    slog.Any("error.stack", formatStacktraceForDatadog(frames)),
    slog.Bool("error.handled", true))
```

#### General Structured Logging Best Practices
- **Consistent field naming**: Use standardized field names (`error.message`, `error.stack`, `error.metadata`)
- **Searchable metadata**: Ensure metadata keys are consistent across the application for better querying
- **Stack filtering**: Log full stacktraces in development, consider filtering in production to reduce noise
- **Error categorization**: Use metadata to categorize errors (e.g., `error_type: "validation"`, `component: "user_service"`)

```go
// Recommended structured logging pattern
func logError(ctx context.Context, handler *rogerr.ErrorHandler, err error, msg string) {
    metadata := handler.Metadata(err)
    frames := handler.Stacktrace(err)
    
    logEntry := slog.With(
        slog.String("error.message", err.Error()),
        slog.String("error.description", msg),
        slog.Any("error.metadata", metadata),
    )
    
    // Include stacktrace for application errors
    if len(frames) > 0 {
        // Filter to only application frames for cleaner logs
        appFrames := make([]rogerr.Frame, 0, len(frames))
        for _, frame := range frames {
            if frame.InApp {
                appFrames = append(appFrames, frame)
            }
        }
        if len(appFrames) > 0 {
            logEntry = logEntry.With(slog.Any("error.stack", appFrames))
        }
    }
    
    logEntry.ErrorContext(ctx, "operation failed")
}
```

