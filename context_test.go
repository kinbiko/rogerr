package rogerr_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/kinbiko/rogerr"
)

func TestMetadata(t *testing.T) {
	var (
		key      = "some key"
		value    = 4321
		err      = errors.New("ooi")
		metadata = map[string]interface{}{"some keys": 2134, "and": "values"}
	)

	t.Run("nil ctx should return nil without panicking even if used later", func(t *testing.T) {
		t.Run("WithMetadatum", func(t *testing.T) {
			ctx := rogerr.WithMetadatum(nil, key, value) //nolint:staticcheck // Testing that we don't do a dumb when users do a dumb
			if ctx != nil {
				t.Fatalf("expected nil ctx returned but got %+v", ctx)
			}
			err := rogerr.Wrap(ctx, err)
			md := rogerr.Metadata(err)
			if md != nil {
				t.Fatalf("expected nil metadata returned from nil ctx but got %+v", md)
			}
		})

		t.Run("WithMetadata", func(t *testing.T) {
			ctx := rogerr.WithMetadata(nil, metadata) //nolint:staticcheck // Testing that we don't do a dumb when users do a dumb
			if ctx != nil {
				t.Fatalf("expected nil ctx returned but got %+v", ctx)
			}
			err := rogerr.Wrap(ctx, err)
			md := rogerr.Metadata(err)
			if md != nil {
				t.Fatalf("expected nil metadata returned from nil ctx but got %+v", md)
			}
		})
	})

	t.Run("successfully extracts all metadata", func(t *testing.T) {
		t.Run("WithMetadata", func(t *testing.T) {
			ctx := context.Background()
			ctx = rogerr.WithMetadata(ctx, metadata)
			err := rogerr.Wrap(ctx, err)
			md := rogerr.Metadata(err)
			for k, v := range md {
				if got := metadata[k]; got != v {
					t.Errorf("expected to find extracted value %v under key %s in stored metadata, but got %v", v, k, got)
				}
			}
			for k, v := range metadata {
				if got := md[k]; got != v {
					t.Errorf("expected to find stored metadata %v under key %s in extracted metadata, but got %v", v, k, got)
				}
			}
		})

		t.Run("WithMetadatum", func(t *testing.T) {
			ctx := context.Background()
			ctx = rogerr.WithMetadatum(ctx, key, value)
			err := rogerr.Wrap(ctx, err)
			md := rogerr.Metadata(err)

			if got := len(md); got != 1 {
				t.Errorf("length of extracted metadata is different (%d) than the expected 1", got)
			}

			got, ok := md[key]
			if !ok {
				t.Fatalf("expected to find key '%s' in extracted metadata, but it was missing", key)
			}
			if got != value {
				t.Errorf("expected to find <%v> in extracted metadata under key '%s', but it was %v", value, key, got)
			}
		})
	})

	t.Run("integration test", func(t *testing.T) {
		/*
			This test case reflects the real world a bit more, in particular:
			- The rogerr.rError type might be wrapped multiple kinds
				- by another rogerr.rError
					- simple wrap-and-return
					- the programmer may choose to add context data *after* they know there's an error, and then wrap the error with a different context.
				- by other errors like fmt.Errorf or errors.New
		*/
		complicatedErrorFlow := func() error {
			ctx := context.Background()
			ctx = rogerr.WithMetadata(ctx, metadata)
			lowestErr := errors.New("some low level err")
			firstWrap := rogerr.Wrap(ctx, lowestErr, "first wrap args")
			wrapWithFmt := fmt.Errorf("wrap with fmt: %w", firstWrap)
			secondWrap := rogerr.Wrap(ctx, wrapWithFmt, "second wrap args")
			ctx = rogerr.WithMetadatum(ctx, key, value)
			thirdWrap := rogerr.Wrap(ctx, secondWrap, "third wrap args")
			ctx = rogerr.WithMetadatum(ctx, key, value+1) // should overwrite lowest context
			fourthWrap := rogerr.Wrap(ctx, thirdWrap, "fourth wrap args")
			return fmt.Errorf("ooi: %w", fourthWrap)
		}

		err := complicatedErrorFlow()

		md := rogerr.Metadata(err)
		if got := len(md); got != len(metadata)+1 /* metadatum adds one */ {
			t.Fatalf("expected metadata to have length %d but had length %d", len(metadata)+1, len(md))
		}

		for k, v := range map[string]interface{}{"some keys": 2134, "and": "values", key: value + 1} {
			if md[k] != v {
				t.Errorf("expected extracted metadata at key '%s' to have value <%v> but was <%v>", k, v, md[k])
			}
		}

		exp := "ooi: fourth wrap args: third wrap args: second wrap args: wrap with fmt: first wrap args: some low level err"
		if got := err.Error(); got != exp {
			t.Errorf("expected error string\n%s\nbut got\n%s\n", exp, got)
		}
	})
}
