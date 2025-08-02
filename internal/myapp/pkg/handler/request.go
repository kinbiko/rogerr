package handler

import (
	"context"

	"github.com/kinbiko/rogerr/internal/myapp/pkg/service"
)

type RequestHandler struct {
	service *service.ProcessingService
}

func NewRequestHandler(svc *service.ProcessingService) *RequestHandler {
	return &RequestHandler{service: svc}
}

func (h *RequestHandler) HandleRequest(ctx context.Context, input string) error {
	return h.processRequest(ctx, input)
}

func (h *RequestHandler) processRequest(ctx context.Context, input string) error {
	validateFunc := func() error {
		return h.service.ProcessData(ctx, input)
	}

	return validateFunc()
}
