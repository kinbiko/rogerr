package rogerr

import "context"

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
	md := Metadata(ctx)
	for k, v := range data {
		md[k] = v
	}
	return context.WithValue(ctx, ctxDataKey, md)
}

// Metadata pulls out all the metadata known by this package as a
// map[key]value from the given context.
func Metadata(ctx context.Context) map[string]interface{} {
	if ctx == nil {
		return nil
	}

	m := map[string]interface{}{}
	if val := ctx.Value(ctxDataKey); val != nil {
		m = val.(map[string]interface{}) // this package owns the ctx key type so this cast is safe.
	}
	return m
}
