package util

import (
	"context"
	"log/slog"
)

type MockHandler struct{} //nolint:revive, gocritic

func NewMockLogger() *slog.Logger { //nolint:revive, gocritic
	return slog.New(&MockHandler{})
}

func (m *MockHandler) Handle(_ context.Context, _ slog.Record) error { //nolint:revive, gocritic
	return nil
}

func (m *MockHandler) Enabled(_ context.Context, _ slog.Level) bool { //nolint:revive, gocritic
	return false
}

func (m *MockHandler) WithAttrs(_ []slog.Attr) slog.Handler { //nolint:revive, gocritic
	return m
}

func (m *MockHandler) WithGroup(_ string) slog.Handler { //nolint:revive, gocritic
	return m
}
