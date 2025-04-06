package grovelog

import (
	"context"
	"io"
	"log/slog"
)

// Logger - обертка вокруг slog.Logger
type Logger struct {
	*slog.Logger
	opts Options
}

// New создает новый логгер с заданными опциями
func New(opts Options) *Logger {
	handler := NewGroveHandler(opts)
	logger := slog.New(handler)

	return &Logger{
		Logger: logger,
		opts:   opts,
	}
}

// Default возвращает логгер с настройками по умолчанию
func Default() *Logger {
	return New(DefaultOptions())
}

// Development возвращает логгер для разработки
func Development() *Logger {
	return New(DevelopmentOptions())
}

// Production возвращает логгер для продакшена
func Production() *Logger {
	return New(ProductionOptions())
}

// With возвращает новый логгер с добавленными атрибутами
func (l *Logger) With(attrs ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(attrs...),
		opts:   l.opts,
	}
}

// WithGroup возвращает новый логгер с добавленной группой
func (l *Logger) WithGroup(name string) *Logger {
	return &Logger{
		Logger: slog.New(l.Logger.Handler().WithGroup(name)),
		opts:   l.opts,
	}
}

// NewWithFile создает логгер, который также пишет в файл
func NewWithFile(path string, opts Options) (*Logger, io.Closer, error) {
	// Создаем обработчик для файла
	fileCloser, fileHandler, err := FileHandler(path, opts)
	if err != nil {
		return nil, nil, err
	}

	// Создаем мультиобработчик, который пишет и в консоль, и в файл
	stdHandler := NewGroveHandler(opts)

	// Создаем мультиобработчик
	multiHandler := &MultiHandler{
		handlers: []slog.Handler{stdHandler, fileHandler},
	}

	logger := slog.New(multiHandler)

	return &Logger{
		Logger: logger,
		opts:   opts,
	}, fileCloser, nil
}

// MultiHandler - обработчик логов, пишущий в несколько мест
type MultiHandler struct {
	handlers []slog.Handler
}

// Enabled проверяет, активен ли обработчик для данного уровня
func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle обрабатывает запись лога для всех обработчиков
func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// WithAttrs возвращает новый обработчик с добавленными атрибутами
func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &MultiHandler{
		handlers: handlers,
	}
}

// WithGroup возвращает новый обработчик с добавленной группой
func (h *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &MultiHandler{
		handlers: handlers,
	}
}
