package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/AlonMell/grovelog"
	"github.com/AlonMell/grovelog/util"
)

// WithLogOp adds an operation name to the context for logging
// It's a convenience wrapper around UpdateLogCtx
// It's example how to use util.UpdateLogCtx
func WithLogOp(ctx context.Context, op string) context.Context {
	return util.UpdateLogCtx(ctx, "op", op)
}

func main() {
	fmt.Println("=== GROVELOG EXAMPLE ===")

	// 1. Different formats demo
	fmt.Println("\n== Logger Formats ==")

	// JSON format
	jsonOpts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.JSON)
	jsonLogger := grovelog.NewLogger(os.Stdout, jsonOpts)
	jsonLogger.Info("JSON formatted log")

	// Plain text format
	plainOpts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Plain)
	plainLogger := grovelog.NewLogger(os.Stdout, plainOpts)
	plainLogger.Info("Plain text formatted log")

	// Color format
	colorOpts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(os.Stdout, colorOpts)
	logger.Info("Color formatted log")

	// 2. Log levels demo
	fmt.Println("\n== Log Levels ==")
	debugOpts := grovelog.NewOptions(slog.LevelDebug, "", grovelog.Color)
	debugLogger := grovelog.NewLogger(os.Stdout, debugOpts)
	debugLogger.Debug("This debug message will be displayed")
	debugLogger.Info("Info message")
	debugLogger.Warn("Warning message")
	debugLogger.Error("Error message")

	// 3. Attributes demo
	fmt.Println("\n== Attributes ==")
	logger.Info("Log with attributes",
		"user_id", 1234,
		"action", "login",
		"timestamp", time.Now())

	// 4. With attributes
	requestLogger := logger.With("request_id", "req-123", "client_ip", "192.168.1.1")
	requestLogger.Info("Processing request with preset attributes")

	// 5. Groups demo
	fmt.Println("\n== Groups ==")
	apiLogger := logger.WithGroup("api")
	apiLogger.Info("API request received",
		"method", "GET",
		"path", "/users",
		"duration_ms", 42)

	// 6. Nested groups
	userApiLogger := apiLogger.WithGroup("users")
	userApiLogger.Info("User API call",
		"user_id", 42,
		"action", "get_profile")

	// 7. Context usage
	fmt.Println("\n== Context ==")
	ctx := context.Background()
	ctx = util.UpdateLogCtx(ctx, "trace_id", "trace-xyz-123")
	ctx = util.UpdateLogCtx(ctx, "session_id", "sess-abc-456")

	// Log with context attributes
	logger.InfoContext(ctx, "Log with context attributes")

	// 8. Error wrapping with context
	fmt.Println("\n== Error Context ==")
	// Create a context with attributes
	ctxWithAttrs := util.UpdateLogCtx(context.Background(), "operation", "data_processing")
	ctxWithAttrs = util.UpdateLogCtx(ctxWithAttrs, "component", "processor")

	// Simulate an error
	err := fmt.Errorf("operation failed: database connection timeout")

	// Wrap the error with context
	wrappedErr := util.WrapCtx(ctxWithAttrs, err)

	// Create a new context and extract attributes from the error
	newCtx := util.ErrorCtx(context.Background(), wrappedErr)

	// Log with the extracted context
	logger.InfoContext(newCtx, "Handling error",
		"error", err.Error(),
		"status", "failed")

	// 9. Group and attributes combination
	fmt.Println("\n== Combined Features ==")
	dbLogger := logger.WithGroup("database").With("db_name", "users_db")
	dbLogger.Info("Executing query",
		"query", "SELECT * FROM users WHERE id = ?",
		"params", []int{42},
		"duration_ms", 10)
}
