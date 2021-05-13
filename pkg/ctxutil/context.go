package ctxutil

import (
	"context"
)

type ContextKey string

func (ck ContextKey) String() string {
	return string(ck)
}

func WithField(ctx context.Context, key ContextKey, value interface{}) context.Context {
	return context.WithValue(ctx, key, value)
}

func ExtractField(ctx context.Context, key ContextKey) (context.Context, interface{}) {
	return ctx, ctx.Value(key)
}
