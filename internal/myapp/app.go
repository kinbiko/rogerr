package main

import (
	"context"

	"github.com/kinbiko/rogerr/internal/mylib"
)

// Application represents our main application
type Application struct {
	name string
}

// NewApplication creates a new application instance
func NewApplication(name string) *Application {
	return &Application{name: name}
}

// Run executes the application logic with various call patterns
func (app *Application) Run(ctx context.Context, input string) error {
	return app.processInput(ctx, input)
}

// processInput is a method that calls library code
func (app *Application) processInput(ctx context.Context, input string) error {
	// Call library function through package-level function
	return callLibraryFunction(ctx, input)
}

// callLibraryFunction is a package-level function in myapp
func callLibraryFunction(ctx context.Context, input string) error {
	// Add anonymous function to call chain
	executeFunc := func() error {
		return mylib.ComplexOperation(ctx, input)
	}
	
	return executeFunc()
}

// StartApplication is the entry point function
func StartApplication(ctx context.Context, appName, input string) error {
	app := NewApplication(appName)
	return app.Run(ctx, input)
}