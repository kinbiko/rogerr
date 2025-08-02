package cmd

import (
	"context"

	"github.com/kinbiko/rogerr/internal/myapp/pkg/handler"
	"github.com/kinbiko/rogerr/internal/myapp/pkg/service"
)

type App struct {
	name    string
	service *service.ProcessingService
	handler *handler.RequestHandler
}

func NewApp() *App {
	svc := service.NewProcessingService()
	handler := handler.NewRequestHandler(svc)
	return &App{name: "demo-app", service: svc, handler: handler}
}

func (app *App) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return app.handler.HandleRequest(ctx, "default")
	}

	return app.handler.HandleRequest(ctx, args[0])
}
