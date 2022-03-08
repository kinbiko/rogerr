package rogerr_test

import (
	"context"
	"testing"

	"github.com/kinbiko/rogerr"
)

func TestMetadata(t *testing.T) {
	var (
		key      = "some key"
		value    = "some value"
		metadata = map[string]interface{}{"some keys": 2134, "and": "values"}
	)

	t.Run("nil ctx should return nil without panicking even if used later", func(t *testing.T) {
		t.Run("WithMetadatum", func(t *testing.T) {
			ctx := rogerr.WithMetadatum(nil, key, value)
			if ctx != nil {
				t.Fatalf("expected nil ctx returned but got %+v", ctx)
			}
			md := rogerr.Metadata(ctx)
			if md != nil {
				t.Fatalf("expected nil metadata returned from nil ctx but got %+v", md)
			}
		})

		t.Run("WithMetadata", func(t *testing.T) {
			ctx := rogerr.WithMetadata(nil, metadata)
			if ctx != nil {
				t.Fatalf("expected nil ctx returned but got %+v", ctx)
			}
			md := rogerr.Metadata(ctx)
			if md != nil {
				t.Fatalf("expected nil metadata returned from nil ctx but got %+v", md)
			}
		})
	})

	t.Run("successfully extracts all metadata", func(t *testing.T) {
		t.Run("WithMetadata", func(t *testing.T) {
			ctx := context.Background()
			ctx = rogerr.WithMetadata(ctx, metadata)
			md := rogerr.Metadata(ctx)
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
			md := rogerr.Metadata(ctx)

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
}
