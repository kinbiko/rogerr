package cmd

import (
	"context"

	"github.com/kinbiko/rogerr/internal/myapp/pkg/handler"
	"github.com/kinbiko/rogerr/internal/myapp/pkg/service"
)

// App represents the main application
type App struct {
	name    string
	service *service.ProcessingService
	handler *handler.RequestHandler
}

// NewApp creates a new application instance
func NewApp(name string) *App {
	svc := service.NewProcessingService()
	hdl := handler.NewRequestHandler(svc)
	
	return &App{
		name:    name,
		service: svc,
		handler: hdl,
	}
}

// Execute runs the application with the given arguments
func (app *App) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return app.handler.HandleRequest(ctx, "default")
	}
	
	return app.handler.HandleRequest(ctx, args[0])
}