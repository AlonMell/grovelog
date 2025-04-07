package util

import (
	"context"
	"errors"
)

// errorWithLogCtx is an error type that carries a logging context
type errorWithLogCtx struct {
	err error
	ctx logCtx
}

func (e *errorWithLogCtx) Error() string {
	return e.err.Error()
}

func (e *errorWithLogCtx) Unwrap() error {
	return e.err
}

// WrapCtx wraps an error with the logging context from the provided context
// This allows context information to propagate along with errors
func WrapCtx(ctx context.Context, err error) error {
	c, _ := getLogCtx(ctx)
	return &errorWithLogCtx{err: err, ctx: c}
}

// ErrorCtx extracts logging context from an error (if it was wrapped with WrapCtx)
// and adds it to the provided context
func ErrorCtx(ctx context.Context, err error) context.Context {
	var errCtx *errorWithLogCtx
	if errors.As(err, &errCtx) {
		return updateLogCtx(ctx, errCtx.ctx)
	}
	return ctx
}
