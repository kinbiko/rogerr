package mylib

import (
	"context"
	"errors"

	"github.com/kinbiko/rogerr"
)

type ErrorService struct {
	handler *rogerr.ErrorHandler
}

func NewErrorService() *ErrorService {
	return &ErrorService{handler: rogerr.NewErrorHandler()}
}

func (es *ErrorService) CreateError(ctx context.Context, msg string) error {
	return es.handler.Wrap(ctx, errors.New("library internal error"), msg)
}

func ProcessDataWithError(ctx context.Context, data string) error {
	service := NewErrorService()
	processFunc := func() error {
		return service.CreateError(ctx, "failed to process data: "+data)
	}
	return processFunc()
}

func ComplexOperation(ctx context.Context, input string) error {
	return intermediateFunction(ctx, input)
}

func intermediateFunction(ctx context.Context, input string) error {
	return ProcessDataWithError(ctx, input)
}
