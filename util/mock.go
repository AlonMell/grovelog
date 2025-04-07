package util

import (
	"context"
	"log/slog"
)

type MockHandler struct{}

func NewMockLogger() *slog.Logger {
	return slog.New(&MockHandler{})
}

func (m *MockHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (m *MockHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (m *MockHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return m
}

func (m *MockHandler) WithGroup(_ string) slog.Handler {
	return m
}
