package grovelog

import (
	"context"
	"encoding/json"
	"io"
	stdLog "log"
	"log/slog"
	"strings"
	"sync"
	"time"

	"slices"

	"github.com/AlonMell/grovelog/util"
	"github.com/fatih/color"
)

// Format defines log output format
type Format int

const (
	// JSON format outputs logs in JSON format
	JSON Format = iota
	// Plain format outputs logs in plain text format
	Plain
	// Color format outputs logs with color highlighting
	Color
)

// DefaultTimeFormat is the default time format
const DefaultTimeFormat = "[15:05:05.000]"

// Options holds configuration options for the logger
type Options struct {
	SlogOpts   *slog.HandlerOptions
	TimeFormat string
	Format     Format
}

// Handler implements the slog.Handler interface with custom formatting
type Handler struct {
	opts Options
	l    *stdLog.Logger

	groups []string // Stores the group hierarchy
	attrs  []slog.Attr

	bufferPool *sync.Pool
	mu         sync.RWMutex
}

// Message represents a formatted log message
type Message struct {
	Time  string
	Level string
	Msg   string
	Atrs  string
}

// NewOptions creates Options with the specified level, time format, and output format
func NewOptions(level slog.Level, timeFormat string, format Format) Options {
	if timeFormat == "" {
		timeFormat = DefaultTimeFormat
	}

	return Options{
		SlogOpts:   &slog.HandlerOptions{Level: level},
		TimeFormat: timeFormat,
		Format:     format,
	}
}

// NewLogger creates a new slog.Logger with the specified options
func NewLogger(out io.Writer, opts Options) *slog.Logger {
	if out == nil {
		out = io.Discard
	}
	h := NewHandler(out, opts)
	return slog.New(h)
}

// NewHandler creates a new slog.Handler
func NewHandler(out io.Writer, opts Options) slog.Handler {
	if out == nil {
		out = io.Discard
	}

	if opts.SlogOpts == nil {
		opts.SlogOpts = &slog.HandlerOptions{Level: slog.LevelInfo}
	}
	if opts.TimeFormat == "" {
		opts.TimeFormat = DefaultTimeFormat
	}

	switch opts.Format {
	case JSON:
		return slog.NewJSONHandler(out, opts.SlogOpts)
	case Plain:
		return slog.NewTextHandler(out, opts.SlogOpts)
	default:
		h := &Handler{
			l:    stdLog.New(out, "", 0),
			opts: opts,
			bufferPool: &sync.Pool{
				New: func() any {
					return new([]byte)
				},
			},
		}
		return h
	}
}

// Handle processes a log record
func (h *Handler) Handle(ctx context.Context, r slog.Record) error { //nolint:gocritic
	ctxAttrs := util.ExtractLogAttrs(ctx)
	if len(ctxAttrs) > 0 {
		for _, attr := range ctxAttrs {
			r.AddAttrs(attr)
		}
	}

	timeStr := h.formatTime(r.Time)

	logMsg := r.Message
	formatLevel := r.Level.String() + ":"

	fields := h.collectFields(r)

	var output string
	if len(fields) > 0 {
		jsonOutput, err := h.marshalFields(fields)
		if err != nil {
			return err
		}
		output = string(jsonOutput)
	}

	type colorFn func(format string, a ...any) string
	levelColorMap := map[slog.Level]colorFn{
		slog.LevelDebug: color.BlueString,
		slog.LevelInfo:  color.GreenString,
		slog.LevelWarn:  color.YellowString,
		slog.LevelError: color.RedString,
	}

	levelColorFunc, ok := levelColorMap[r.Level]
	if !ok {
		levelColorFunc = color.WhiteString // Default color for unknown levels
	}

	level := levelColorFunc(formatLevel)
	msg := Message{
		Time:  timeStr,
		Level: level,
		Msg:   color.CyanString(logMsg),
		Atrs:  color.WhiteString(output),
	}

	h.l.Println(msg.Time, msg.Level, msg.Msg, msg.Atrs)
	return nil
}

type jsonWriter struct {
	buf *[]byte
}

func (w *jsonWriter) Write(p []byte) (n int, err error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}

// marshalFields optimizes JSON serialization of fields
func (h *Handler) marshalFields(fields map[string]any) ([]byte, error) {
	if h.bufferPool != nil {
		bufPtr, ok := h.bufferPool.Get().(*[]byte)
		if !ok || bufPtr == nil {
			return json.MarshalIndent(fields, "", "  ")
		}

		*bufPtr = (*bufPtr)[:0] // Clear buffer

		encoder := json.NewEncoder(io.MultiWriter(io.Discard, &jsonWriter{buf: bufPtr}))
		encoder.SetIndent("", "  ")

		err := encoder.Encode(fields)
		jsonData := *bufPtr
		h.bufferPool.Put(bufPtr) // Return buffer to pool

		if err != nil {
			return nil, err
		}

		// Remove trailing newline added by json.Encoder
		if len(jsonData) > 0 && jsonData[len(jsonData)-1] == '\n' {
			jsonData = jsonData[:len(jsonData)-1]
		}

		result := make([]byte, len(jsonData))
		copy(result, jsonData)
		return result, nil
	}

	return json.MarshalIndent(fields, "", "  ")
}

func (h *Handler) formatTime(t time.Time) string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	format := h.opts.TimeFormat
	if format == "" {
		format = DefaultTimeFormat
	}

	return t.Format(format)
}

func (h *Handler) collectFields(r slog.Record) map[string]any { //nolint:gocritic
	fields := make(map[string]any, r.NumAttrs()+len(h.attrs))

	h.mu.RLock()
	groupPrefix := ""
	if len(h.groups) > 0 {
		groupPrefix = strings.Join(h.groups, ".") + "."
	}

	var processAttr func(a slog.Attr, prefix string)
	processAttr = func(a slog.Attr, prefix string) {
		if a.Key == "" {
			return
		}

		fullKey := prefix + a.Key

		if a.Value.Kind() == slog.KindGroup {
			group := a.Value.Group()
			for _, groupAttr := range group {
				if groupAttr.Key != "" {
					processAttr(groupAttr, fullKey+".")
				}
			}
		} else {
			fields[fullKey] = a.Value.Any()
		}
	}

	r.Attrs(func(a slog.Attr) bool {
		processAttr(a, groupPrefix)
		return true
	})

	for _, a := range h.attrs {
		processAttr(a, groupPrefix)
	}
	h.mu.RUnlock()

	return fields
}

// Enabled determines if this level should be logged
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	minLevel := slog.LevelInfo
	if h.opts.SlogOpts != nil && h.opts.SlogOpts.Level != nil {
		minLevel = h.opts.SlogOpts.Level.Level()
	}
	return level >= minLevel
}

// WithAttrs returns a new Handler with the given attributes added
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	validAttrs := make([]slog.Attr, 0, len(attrs))
	for _, attr := range attrs {
		if attr.Key != "" {
			validAttrs = append(validAttrs, attr)
		}
	}

	if len(validAttrs) == 0 {
		return h
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	return &Handler{
		l:          h.l,
		opts:       h.opts,
		groups:     slices.Clone(h.groups),
		bufferPool: h.bufferPool,
		attrs:      slices.Concat(slices.Clone(h.attrs), validAttrs),
	}
}

// WithGroup returns a new Handler with the given group name added
func (h *Handler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	// Create a new handler with the same attributes but a new group
	newHandler := &Handler{
		l:          h.l,
		opts:       h.opts,
		attrs:      slices.Clone(h.attrs),
		groups:     append(slices.Clone(h.groups), name),
		bufferPool: h.bufferPool,
	}

	return newHandler
}
