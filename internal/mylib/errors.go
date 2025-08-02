package mylib

import (
	"context"
	"errors"

	"github.com/kinbiko/rogerr"
)

// ErrorService handles error creation for the library
type ErrorService struct {
	handler *rogerr.ErrorHandler
}

// NewErrorService creates a new error service with stacktrace enabled
func NewErrorService() *ErrorService {
	return &ErrorService{
		handler: rogerr.NewErrorHandler(),
	}
}

// CreateError creates an error with stacktrace through method call
func (es *ErrorService) CreateError(ctx context.Context, msg string) error {
	baseErr := errors.New("library internal error")
	return es.handler.Wrap(ctx, baseErr, msg)
}

// ProcessDataWithError is a package-level function that creates an error
func ProcessDataWithError(ctx context.Context, data string) error {
	service := NewErrorService()
	
	// Call through anonymous function to add another stack frame
	processFunc := func() error {
		return service.CreateError(ctx, "failed to process data: "+data)
	}
	
	return processFunc()
}

// ComplexOperation demonstrates multiple call levels
func ComplexOperation(ctx context.Context, input string) error {
	return intermediateFunction(ctx, input)
}

// intermediateFunction is a helper function in the call chain
func intermediateFunction(ctx context.Context, input string) error {
	return ProcessDataWithError(ctx, input)
}