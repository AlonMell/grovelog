package util

import (
	"context"
	"errors"
)

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

func WrapCtx(ctx context.Context, err error) error {
	c, _ := getLogCtx(ctx)
	return &errorWithLogCtx{err: err, ctx: c}
}

func ErrorCtx(ctx context.Context, err error) context.Context {
	var errCtx *errorWithLogCtx
	if errors.As(err, &errCtx) {
		return updateLogCtx(ctx, errCtx.ctx)
	}
	return ctx
}
