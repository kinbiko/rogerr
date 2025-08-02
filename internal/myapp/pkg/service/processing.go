package service

import (
	"context"

	"github.com/kinbiko/rogerr/internal/mylib"
)

type ProcessingService struct {
	name string
}

func NewProcessingService() *ProcessingService {
	return &ProcessingService{
		name: "data-processor",
	}
}

func (s *ProcessingService) ProcessData(ctx context.Context, data string) error {
	return s.executeBusinessLogic(ctx, data)
}

func (s *ProcessingService) executeBusinessLogic(ctx context.Context, data string) error {
	callLibrary := func() error {
		return mylib.ComplexOperation(ctx, data)
	}

	return callLibrary()
}
