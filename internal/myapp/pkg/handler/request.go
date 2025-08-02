package handler

import (
	"context"

	"github.com/kinbiko/rogerr/internal/myapp/pkg/service"
)

// RequestHandler handles incoming requests
type RequestHandler struct {
	service *service.ProcessingService
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(svc *service.ProcessingService) *RequestHandler {
	return &RequestHandler{
		service: svc,
	}
}

// HandleRequest processes a request with the given input
func (h *RequestHandler) HandleRequest(ctx context.Context, input string) error {
	return h.processRequest(ctx, input)
}

// processRequest is an internal method that validates and processes the request
func (h *RequestHandler) processRequest(ctx context.Context, input string) error {
	// Add validation layer
	validateFunc := func() error {
		return h.service.ProcessData(ctx, input)
	}
	
	return validateFunc()
}