package grovelog_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/AlonMell/grovelog"
)

// TestNewLogger tests the creation of loggers with different formats
func TestNewLogger(t *testing.T) {
	tests := []struct {
		name        string
		format      grovelog.Format
		expectRegex string
	}{
		{
			name:        "JSONFormat",
			format:      grovelog.JSON,
			expectRegex: `\{"time":".*","level":"INFO","msg":"test message"\}`,
		},
		{
			name:        "PlainFormat",
			format:      grovelog.Plain,
			expectRegex: `time=.* level=INFO msg="test message"`,
		},
		{
			name:        "ColorFormat",
			format:      grovelog.Color,
			expectRegex: `\[\d{2}:\d{2}:\d{2}\.\d{3}\] INFO: test message`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			opts := grovelog.NewOptions(slog.LevelInfo, "", tt.format)
			logger := grovelog.NewLogger(&buf, opts)
			logger.Info("test message")

			logOutput := buf.String()
			matched, err := regexp.MatchString(tt.expectRegex, logOutput)
			if err != nil {
				t.Fatalf("Error matching regex: %v", err)
			}
			if !matched {
				t.Errorf("Log output did not match expected format.\nGot: %s\nExpected regex: %s", logOutput, tt.expectRegex)
			}
		})
	}
}

// TestLogLevels tests that log levels are properly filtered
func TestLogLevels(t *testing.T) {
	var buf bytes.Buffer
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(&buf, opts)

	// Debug should be filtered out
	logger.Debug("debug message")
	if buf.Len() > 0 {
		t.Errorf("Debug message should have been filtered out, but got: %s", buf.String())
	}
	buf.Reset()

	// Info should pass through
	logger.Info("info message")
	if buf.Len() == 0 {
		t.Error("Info message should have been logged, but buffer is empty")
	}
	buf.Reset()

	// Error should pass through
	logger.Error("error message")
	if buf.Len() == 0 {
		t.Error("Error message should have been logged, but buffer is empty")
	}
}

// TestWithAttrs tests the WithAttrs functionality
func TestWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(&buf, opts)

	// Log with attributes
	attrLogger := logger.With("key1", "value1", "key2", 42)
	attrLogger.Info("message with attributes")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "key1") || !strings.Contains(logOutput, "value1") ||
		!strings.Contains(logOutput, "key2") || !strings.Contains(logOutput, "42") {
		t.Errorf("Log output missing attributes. Got: %s", logOutput)
	}
}

// TestWithGroup tests the WithGroup functionality
func TestWithGroup(t *testing.T) {
	var buf bytes.Buffer
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(&buf, opts)

	// Log with group
	groupLogger := logger.WithGroup("group1")
	groupLogger.Info("message with group", "key1", "value1")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "group1.key1") {
		t.Errorf("Log output missing group prefix. Got: %s", logOutput)
	}
}

// TestNestedGroups tests nested groups
func TestNestedGroups(t *testing.T) {
	var buf bytes.Buffer
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(&buf, opts)

	// Create nested groups
	level1 := logger.WithGroup("level1")
	level2 := level1.WithGroup("level2")
	level3 := level2.WithGroup("level3")

	level3.Info("nested message", "key", "value")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "level1.level2.level3.key") {
		t.Errorf("Log output missing nested group prefixes. Got: %s", logOutput)
	}
}

// TestGroupWithAttrs tests groups with attributes
func TestGroupWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(&buf, opts)

	// Add group and attributes
	groupLogger := logger.WithGroup("group").With("attr", "value")
	groupLogger.Info("message")

	logOutput := buf.String()
	if !strings.Contains(logOutput, "group.attr") {
		t.Errorf("Log output missing group prefix with attribute. Got: %s", logOutput)
	}
}

// TestTimeFormat tests custom time formats
func TestTimeFormat(t *testing.T) {
	var buf bytes.Buffer
	customFormat := "2006-01-02"
	opts := grovelog.NewOptions(slog.LevelInfo, customFormat, grovelog.Color)
	logger := grovelog.NewLogger(&buf, opts)

	logger.Info("custom time format")

	// Get today's date in the expected format
	today := time.Now().Format(customFormat)

	logOutput := buf.String()
	if !strings.Contains(logOutput, today) {
		t.Errorf("Log output has wrong time format. Got: %s, Expected to contain: %s", logOutput, today)
	}
}

// TestLogAttr tests the LogAttrs method with nested groups
func TestLogAttrs(t *testing.T) {
	var buf bytes.Buffer
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(&buf, opts)

	// Log with a group attribute
	logger.LogAttrs(context.Background(), slog.LevelInfo, "grouped attrs",
		slog.Group("metadata",
			slog.String("version", "1.0.0"),
			slog.Bool("active", true),
		),
	)

	logOutput := buf.String()
	if !strings.Contains(logOutput, "metadata.version") || !strings.Contains(logOutput, "metadata.active") {
		t.Errorf("Log output missing grouped attributes. Got: %s", logOutput)
	}
}

// TestConcurrentUsage tests thread safety
func TestConcurrentUsage(t *testing.T) {
	var buf bytes.Buffer
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(&buf, opts)

	// Number of goroutines to use for testing
	const goroutines = 100
	// Number of logs per goroutine
	const logsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Start multiple goroutines that all use the logger
	for i := range goroutines {
		go func(id int) {
			defer wg.Done()

			// Each goroutine creates its own derived logger
			threadLogger := logger.With("goroutine", id)

			for j := range logsPerGoroutine {
				// Mix of operations to test concurrency
				switch j % 3 {
				case 0:
					groupLogger := threadLogger.WithGroup("group")
					groupLogger.Info("grouped log", "count", j)
				case 1:
					attrLogger := threadLogger.With("count", j)
					attrLogger.Info("attr log")
				default:
					threadLogger.Info("simple log", "count", j)
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// If we get here without panics or deadlocks, the test passes
	// We can also check that some logging actually happened
	if buf.Len() == 0 {
		t.Error("Expected log output, but buffer is empty")
	}
}

// TestBigPayload tests logging with large amounts of data
func TestBigPayload(t *testing.T) {
	var buf bytes.Buffer
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(&buf, opts)

	// Create a large payload
	payload := make([]string, 1000)
	for i := range 1000 {
		payload[i] = "test data entry with some content to make it larger than just an index"
	}

	// Log with the large payload
	logger.Info("big payload", "data", payload)

	// Check that something was logged
	if buf.Len() == 0 {
		t.Error("Expected log output for big payload, but buffer is empty")
	}
}

// CapturingHandler is a test helper that captures logged records
type CapturingHandler struct {
	mu      sync.Mutex
	Records []slog.Record
}

func (h *CapturingHandler) Handle(ctx context.Context, r slog.Record) error { //nolint:gocritic
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Records = append(h.Records, r)
	return nil
}

func (h *CapturingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *CapturingHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *CapturingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

// TestFormatValid tests handling of valid format options
func TestFormatValid(t *testing.T) {
	var buf bytes.Buffer

	// Test each valid format option
	formats := []grovelog.Format{grovelog.JSON, grovelog.Plain, grovelog.Color}
	for _, format := range formats {
		opts := grovelog.NewOptions(slog.LevelInfo, "", format)
		logger := grovelog.NewLogger(&buf, opts)

		// If we get here without panics, the test passes
		logger.Info("test message")
		buf.Reset()
	}
}

// BenchmarkHandleBasic benchmarks basic logging
func BenchmarkHandleBasic(b *testing.B) {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(io.Discard, opts)

	for b.Loop() {
		logger.Info("benchmark message")
	}
}

// BenchmarkHandleWithAttrs benchmarks logging with attributes
func BenchmarkHandleWithAttrs(b *testing.B) {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(io.Discard, opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message",
			"string", "value",
			"int", 42,
			"bool", true,
			"float", 3.14)
	}
}

// BenchmarkHandleWithGroups benchmarks logging with groups
func BenchmarkHandleWithGroups(b *testing.B) {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(io.Discard, opts)
	groupedLogger := logger.WithGroup("group1").WithGroup("group2")

	for b.Loop() {
		groupedLogger.Info("benchmark message", "key", "value")
	}
}

// BenchmarkWithAttrs benchmarks the WithAttrs method
func BenchmarkWithAttrs(b *testing.B) {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(io.Discard, opts)

	for b.Loop() {
		_ = logger.With("key1", "value1", "key2", 42)
	}
}

// BenchmarkWithGroup benchmarks the WithGroup method
func BenchmarkWithGroup(b *testing.B) {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(io.Discard, opts)

	for b.Loop() {
		_ = logger.WithGroup("group")
	}
}

// BenchmarkNestedGroups benchmarks creating nested groups
func BenchmarkNestedGroups(b *testing.B) {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(io.Discard, opts)

	for b.Loop() {
		_ = logger.WithGroup("group1").WithGroup("group2").WithGroup("group3")
	}
}

// BenchmarkHandleJSON benchmarks JSON format logging
func BenchmarkHandleJSON(b *testing.B) {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.JSON)
	logger := grovelog.NewLogger(io.Discard, opts)

	for b.Loop() {
		logger.Info("benchmark message", "key", "value")
	}
}

// BenchmarkConcurrentLogging benchmarks concurrent logging
func BenchmarkConcurrentLogging(b *testing.B) {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(io.Discard, opts)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine gets its own logger with unique attributes
		threadLogger := logger.With("goroutine", "id")

		for pb.Next() {
			threadLogger.Info("concurrent benchmark")
		}
	})
}

// BenchmarkIndirectMarshalFields benchmarks the marshaling of fields indirectly
func BenchmarkIndirectMarshalFields(b *testing.B) {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(io.Discard, opts)

	// Create fields to log

	for b.Loop() {
		logger.Info("benchmark",
			"string", "value",
			"int", 42,
			"bool", true,
			"float", 3.14,
			"array", []string{"one", "two", "three"},
			"nested", map[string]any{
				"key1": "value1",
				"key2": 2,
			},
		)
	}
}

// BenchmarkCompareToStandardLogger benchmarks against the standard slog
func BenchmarkCompareToStandardLogger(b *testing.B) {
	b.Run("StandardJSONLogger", func(b *testing.B) {
		handler := slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger := slog.New(handler)

		b.ResetTimer()
		for b.Loop() {
			logger.Info("benchmark message", "key", "value")
		}
	})

	b.Run("GroveJSONLogger", func(b *testing.B) {
		opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.JSON)
		logger := grovelog.NewLogger(io.Discard, opts)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("benchmark message", "key", "value")
		}
	})

	b.Run("StandardTextLogger", func(b *testing.B) {
		handler := slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		logger := slog.New(handler)

		b.ResetTimer()
		for b.Loop() {
			logger.Info("benchmark message", "key", "value")
		}
	})

	b.Run("GrovePlainLogger", func(b *testing.B) {
		opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Plain)
		logger := grovelog.NewLogger(io.Discard, opts)

		b.ResetTimer()
		for b.Loop() {
			logger.Info("benchmark message", "key", "value")
		}
	})

	b.Run("GroveColorLogger", func(b *testing.B) {
		opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
		logger := grovelog.NewLogger(io.Discard, opts)

		b.ResetTimer()
		for b.Loop() {
			logger.Info("benchmark message", "key", "value")
		}
	})
}

// TestJSONFormat verifies JSON output can be properly parsed
func TestJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.JSON)
	logger := grovelog.NewLogger(&buf, opts)

	logger.Info("json test", "key", "value")

	var jsonMap map[string]any
	err := json.Unmarshal(buf.Bytes(), &jsonMap)
	if err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
	}

	// Verify expected fields exist
	if jsonMap["msg"] != "json test" {
		t.Errorf("Expected msg field to be 'json test', got %v", jsonMap["msg"])
	}

	if jsonMap["key"] != "value" {
		t.Errorf("Expected key field to be 'value', got %v", jsonMap["key"])
	}
}
