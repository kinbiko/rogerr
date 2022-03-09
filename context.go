package rogerr

import (
	"context"
	"errors"
)

type ctxKey int

const ctxDataKey ctxKey = 0

// WithMetadatum attaches the given key and value to the rogerr metadata
// associated with this context.
// Returns a new context with the metadatum attached, or nil if the given ctx was nil.
func WithMetadatum(ctx context.Context, key string, value interface{}) context.Context {
	return WithMetadata(ctx, map[string]interface{}{key: value})
}

// WithMetadata attaches the given keys and values to the rogerr metadata
// associated with this context.
// Returns a new context with the metadata attached, or nil if the given ctx was nil.
func WithMetadata(ctx context.Context, data map[string]interface{}) context.Context {
	if ctx == nil {
		return nil
	}
	md := getOrInitializeMetadata(ctx)
	for k, v := range data {
		md[k] = v
	}
	return context.WithValue(ctx, ctxDataKey, md)
}

// Metadata pulls out all the metadata known by this package as a
// map[key]value from the given error.
func Metadata(err error) map[string]interface{} {
	rErr := &rError{}

	// Yes, that's a double pointer. The error type is a struct pointer, and
	// errors.As requires a pointer to a type that implements the error
	// interface, which is *Error, hence passing **Error here.
	errors.As(err, &rErr)
	return getOrInitializeMetadata(rErr.ctx)
}

func getOrInitializeMetadata(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return nil
	}

	m := map[string]interface{}{}
	if val := ctx.Value(ctxDataKey); val != nil {
		m = val.(map[string]interface{}) // this package owns the ctx key type so this cast is safe.
	}
	return m
}
