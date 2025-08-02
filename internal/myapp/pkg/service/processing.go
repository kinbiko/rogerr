package service

import (
	"context"

	"github.com/kinbiko/rogerr/internal/mylib"
)

// ProcessingService handles business logic
type ProcessingService struct {
	name string
}

// NewProcessingService creates a new processing service
func NewProcessingService() *ProcessingService {
	return &ProcessingService{
		name: "data-processor",
	}
}

// ProcessData processes the given data through the business logic layer
func (s *ProcessingService) ProcessData(ctx context.Context, data string) error {
	return s.executeBusinessLogic(ctx, data)
}

// executeBusinessLogic runs the core business logic
func (s *ProcessingService) executeBusinessLogic(ctx context.Context, data string) error {
	// Call library function through service layer
	callLibrary := func() error {
		return mylib.ComplexOperation(ctx, data)
	}
	
	return callLibrary()
}