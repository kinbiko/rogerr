package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/kinbiko/rogerr"
	"github.com/kinbiko/rogerr/internal/myapp/cmd"
)

func main() {
	// Configure slog for JSON output
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	args := os.Args[1:] // Skip program name
	if len(args) == 0 {
		args = []string{"demo-data"} // Default argument for demo
	}

	err := run(args)
	if err != nil {
		// Extract stacktrace using ErrorHandler
		handler := rogerr.NewErrorHandler()
		frames := handler.Stacktrace(err)

		// Convert frames to structured log format
		frameData := make([]map[string]interface{}, len(frames))
		for i, frame := range frames {
			frameData[i] = map[string]interface{}{
				"function": frame.Function,
				"file":     frame.File,
				"line":     frame.Line,
				"in_app":   frame.InApp,
			}
		}

		// Log error with stacktrace as structured JSON
		slog.Error("application error",
			slog.String("error.message", err.Error()),
			slog.Any("error.stacktrace", frameData),
			slog.Int("stacktrace.frame_count", len(frames)),
		)
	}
}

func run(args []string) error {
	ctx := context.Background()
	app := cmd.NewApp("demo-app")
	return app.Execute(ctx, args)
}
