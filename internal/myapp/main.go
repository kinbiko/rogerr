package main

import (
	"context"
	"fmt"
	"os"

	"github.com/kinbiko/rogerr"
	"github.com/kinbiko/rogerr/internal/myapp/cmd"
)

func main() {
	args := os.Args[1:] // Skip program name
	if len(args) == 0 {
		args = []string{"demo-data"} // Default argument for demo
	}
	
	err := run(args)
	if err != nil {
		// Extract stacktrace using ErrorHandler
		handler := rogerr.NewErrorHandler()
		frames := handler.Stacktrace(err)
		
		// Print stacktrace to stdout
		for i, frame := range frames {
			fmt.Printf("[%d] %s (%s:%d) InApp=%v\n", i, frame.Function, frame.File, frame.Line, frame.InApp)
		}
	}
}

func run(args []string) error {
	ctx := context.Background()
	app := cmd.NewApp("demo-app")
	return app.Execute(ctx, args)
}