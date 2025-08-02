package rogerr_test

import (
	"context"
	"fmt"

	"github.com/kinbiko/rogerr"
)

func ExampleWrap() {
	handler := rogerr.NewErrorHandler()

	someFuncWithAProblem := func(_ context.Context) error {
		return fmt.Errorf("some low level err")
	}

	someFuncThatWrapsWithRogerr := func(ctx context.Context) error {
		// Attach some projectID to the context as structured metadata
		ctx = rogerr.WithMetadatum(ctx, "projectID", 123)

		err := someFuncWithAProblem(ctx)
		if err != nil {
			return handler.Wrap(ctx, err, "wrap args")
		}
		return nil
	}

	someFuncThatWrapsARogerrError := func(ctx context.Context) error {
		err := someFuncThatWrapsWithRogerr(ctx)
		if err != nil {
			return fmt.Errorf("wrap with fmt: %w", err)
		}
		return nil
	}

	err := someFuncThatWrapsARogerrError(context.Background())
	md := rogerr.Metadata(err)
	fmt.Println(err.Error())     // error message should be cleanly wrapped
	fmt.Println(md["projectID"]) // structured metadata should be available
	//output:
	// wrap with fmt: wrap args: some low level err
	// 123
}
