package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/kinbiko/rogerr"
	"github.com/kinbiko/rogerr/internal/myapp/cmd"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	args := os.Args[1:] // Skip program name
	if len(args) == 0 {
		args = []string{"demo-data"}
	}

	err := run(args)
	if err == nil {
		return
	}

	frames := rogerr.NewErrorHandler().Stacktrace(err)

	// Convert frames to OTEL logging data model format
	frameData := make([]map[string]interface{}, len(frames))
	for i, frame := range frames {
		frameData[i] = map[string]interface{}{
			"code.function": frame.Function,
			"code.filepath": frame.File,
			"code.lineno":   frame.Line,
			"code.namespace": func() string {
				if frame.InApp {
					return "application"
				}
				return "dependency"
			}(),
		}
	}

	// Log using OTEL logging data model semantic conventions
	slog.Error("Exception occurred",
		slog.String("exception.type", "ApplicationError"),
		slog.String("exception.message", err.Error()),
		slog.Any("exception.stacktrace", frameData),
		slog.String("service.name", "demo-app"),
		slog.String("service.version", "1.0.0"),
	)
}

func run(args []string) error {
	ctx := context.Background()
	app := cmd.NewApp()
	return app.Execute(ctx, args)
}
