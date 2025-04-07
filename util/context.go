package util

import (
	"context"
	"log/slog"
	"maps"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
)

type logCtx map[string]any

// Example of a function that uses the logCtx
func WithLogOp(ctx context.Context, op string) context.Context {
	return UpdateLogCtx(ctx, "op", op)
}

func UpdateLogCtx(ctx context.Context, key string, value any) context.Context {
	return updateLogCtx(ctx, logCtx{key: value})
}

func ExtractLogAttrs(ctx context.Context) []slog.Attr {
	if lctx, ok := getLogCtx(ctx); ok {
		attrs := make([]slog.Attr, 0, len(lctx))
		for k, v := range lctx {
			attrs = append(attrs, KV(k, v))
		}
		return attrs
	}
	return nil
}

func updateLogCtx(ctx context.Context, newCtx logCtx) context.Context {
	if existingCtx, ok := getLogCtx(ctx); ok {
		maps.Copy(existingCtx, newCtx)
		return context.WithValue(ctx, logCtxKey, existingCtx)
	}
	return context.WithValue(ctx, logCtxKey, newCtx)
}

func getLogCtx(ctx context.Context) (logCtx, bool) {
	c, ok := ctx.Value(logCtxKey).(logCtx)
	return c, ok
}
