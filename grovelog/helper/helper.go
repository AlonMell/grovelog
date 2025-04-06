package helper

import (
	"context"
	"log/slog"
	"runtime"
	"strconv"
)

// Err создает атрибут slog.Attr из ошибки
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

// KV создает пару ключ-значение для логирования
func KV(key string, value interface{}) slog.Attr {
	return slog.Any(key, value)
}

// Caller возвращает файл и номер строки вызывающей функции
func Caller(skip int) slog.Attr {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return slog.String("caller", "unknown")
	}
	return slog.String("caller", file+":"+strconv.Itoa(line))
}

// WithContext возвращает логгер из контекста
func WithContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(loggerKey{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

// ContextWithLogger возвращает новый контекст с прикрепленным логгером
func ContextWithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

type loggerKey struct{}
