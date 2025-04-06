package grovelog

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"
	"os"
	"slices"

	"github.com/fatih/color"
)

// GroveHandler - обработчик логов с поддержкой цветного вывода
type GroveHandler struct {
	opts        Options
	attrs       []slog.Attr
	groups      []string
	logger      *stdLog.Logger
	jsonHandler slog.Handler
	textHandler slog.Handler
}

// NewGroveHandler создает новый обработчик с заданными опциями
func NewGroveHandler(opts Options) *GroveHandler {
	levelVar := new(slog.LevelVar)
	levelVar.Set(opts.Level)

	slogOpts := &slog.HandlerOptions{
		Level:     levelVar,
		AddSource: opts.AddSource,
	}

	jsonHandler := slog.NewJSONHandler(opts.Output, slogOpts)
	textHandler := slog.NewTextHandler(opts.Output, slogOpts)

	return &GroveHandler{
		opts:        opts,
		logger:      stdLog.New(opts.Output, "", 0),
		jsonHandler: jsonHandler,
		textHandler: textHandler,
	}
}

// Enabled сообщает, обрабатывает ли обработчик записи на данном уровне
func (h *GroveHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level
}

// Handle обрабатывает запись лога
func (h *GroveHandler) Handle(ctx context.Context, r slog.Record) error {
	// Если формат не цветной, используем стандартные обработчики
	switch h.opts.Format {
	case JSONFormat:
		return h.jsonHandler.Handle(ctx, r)
	case TextFormat:
		return h.textHandler.Handle(ctx, r)
	}

	// Форматирование времени
	timeStr := ""
	if h.opts.TimeFormat != "" {
		timeStr = r.Time.Format(h.opts.TimeFormat)
	}

	// Форматирование уровня с цветом
	levelStr := h.formatLevel(r.Level)

	// Форматирование сообщения
	msg := color.CyanString(r.Message)

	// Сбор атрибутов
	attrs := make(map[string]any)
	r.Attrs(func(a slog.Attr) bool {
		attrs[a.Key] = a.Value.Any()
		return true
	})

	// Добавление атрибутов обработчика
	for _, a := range h.attrs {
		attrs[a.Key] = a.Value.Any()
	}

	// Форматирование атрибутов как JSON, если они есть
	attrsStr := ""
	if len(attrs) > 0 {
		b, err := json.MarshalIndent(attrs, "", "  ")
		if err == nil {
			attrsStr = color.WhiteString(string(b))
		}
	}

	// Вывод строки лога
	h.logger.Println(timeStr, levelStr, msg, attrsStr)
	return nil
}

// WithAttrs возвращает новый обработчик с добавленными атрибутами
func (h *GroveHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := &GroveHandler{
		opts:        h.opts,
		attrs:       slices.Concat(h.attrs, attrs),
		groups:      h.groups,
		logger:      h.logger,
		jsonHandler: h.jsonHandler.WithAttrs(attrs),
		textHandler: h.textHandler.WithAttrs(attrs),
	}
	return newHandler
}

// WithGroup возвращает новый обработчик с добавленной группой
func (h *GroveHandler) WithGroup(name string) slog.Handler {
	return &GroveHandler{
		opts:        h.opts,
		attrs:       h.attrs,
		groups:      append(slices.Clone(h.groups), name),
		logger:      h.logger,
		jsonHandler: h.jsonHandler.WithGroup(name),
		textHandler: h.textHandler.WithGroup(name),
	}
}

// formatLevel возвращает цветную строку для уровня лога
func (h *GroveHandler) formatLevel(level slog.Level) string {
	var levelColorFunc func(format string, a ...any) string
	switch {
	case level >= slog.LevelError:
		levelColorFunc = color.RedString
	case level >= slog.LevelWarn:
		levelColorFunc = color.YellowString
	case level >= slog.LevelInfo:
		levelColorFunc = color.BlueString
	default:
		levelColorFunc = color.MagentaString
	}

	return levelColorFunc("%s:", level.String())
}

// FileHandler создает обработчик логов для файла
func FileHandler(path string, opts Options) (io.WriteCloser, slog.Handler, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}

	// Для файлов обычно используем формат JSON или текст
	fileOpts := opts
	if fileOpts.Format == ColorFormat {
		fileOpts.Format = JSONFormat
	}
	fileOpts.Output = f

	return f, NewGroveHandler(fileOpts), nil
}
