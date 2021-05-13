package ctxutil

import (
	"context"

	"github.com/google/uuid"
)

var (
	ContextRequestIDKey ContextKey = "request_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextRequestIDKey, requestID)
}

func ExtractRequestID(ctx context.Context) (context.Context, string) {
	ctx, val := ExtractField(ctx, ContextRequestIDKey)
	if requestID, ok := val.(string); ok && len(requestID) > 0 {
		return ctx, requestID
	}
	requestID := uuid.NewString()
	ctx = WithRequestID(ctx, requestID)
	return ctx, requestID
}
