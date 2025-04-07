package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"log/slog"

	"github.com/AlonMell/grovelog"
)

func main() {
	// Example showing different formatting options
	fmt.Println("=== DIFFERENT LOGGER FORMATS ===")
	demonstrateFormats()

	// Example showing different log levels
	fmt.Println("\n=== LOG LEVELS ===")
	demonstrateLogLevels()

	// Example showing attributes and context values
	fmt.Println("\n=== ATTRIBUTES ===")
	demonstrateAttributes()

	// Example showing groups
	fmt.Println("\n=== GROUPS ===")
	demonstrateGroups()

	// Example showing nested groups
	fmt.Println("\n=== NESTED GROUPS ===")
	demonstrateNestedGroups()

	// Example showing multi-threading
	fmt.Println("\n=== CONCURRENT USAGE ===")
	demonstrateConcurrentUsage()

	// Example of practical structured logging
	fmt.Println("\n=== PRACTICAL EXAMPLE ===")
	practicalExample()
}

// demonstrateFormats shows different logger formats
func demonstrateFormats() {
	// JSON format
	jsonOpts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.JSON)
	jsonLogger := grovelog.NewLogger(os.Stdout, jsonOpts)
	jsonLogger.Info("This is a JSON formatted log")

	// Plain text format
	plainOpts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Plain)
	plainLogger := grovelog.NewLogger(os.Stdout, plainOpts)
	plainLogger.Info("This is a plain text formatted log")

	// Color format (default)
	colorOpts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	colorLogger := grovelog.NewLogger(os.Stdout, colorOpts)
	colorLogger.Info("This is a color formatted log")

	// Custom time format
	customTimeOpts := grovelog.NewOptions(slog.LevelInfo, "2006-01-02 15:04:05", grovelog.Color)
	customTimeLogger := grovelog.NewLogger(os.Stdout, customTimeOpts)
	customTimeLogger.Info("This log has a custom time format")
}

// demonstrateLogLevels shows different log levels
func demonstrateLogLevels() {
	// Create a logger with debug level
	opts := grovelog.NewOptions(slog.LevelDebug, "", grovelog.Color)
	logger := grovelog.NewLogger(os.Stdout, opts)

	// Log at different levels
	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// Log with formatted message
	logger.Info("User logged in", "user_id", 12345, "source", "api")
}

// demonstrateAttributes shows how to add attributes to logs
func demonstrateAttributes() {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(os.Stdout, opts)

	// Base logger
	logger.Info("Base logger message")

	// Logger with added attributes
	requestLogger := logger.With("request_id", "req-123", "client_ip", "192.168.1.1")
	requestLogger.Info("Processing request")

	// More attributes added
	requestLogger.Info("Request completed",
		"status", 200,
		"duration_ms", 342,
		"response_size", 2048)

	// With slog.Attr for more control
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Using LogAttrs method",
		slog.Int("user_id", 42),
		slog.Time("timestamp", time.Now()),
		slog.Group("metadata",
			slog.String("version", "v1.0.0"),
			slog.Bool("beta", true),
		),
	)
}

// demonstrateGroups shows how to use groups
func demonstrateGroups() {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(os.Stdout, opts)

	// Base logger
	logger.Info("Base logger")

	// Logger with a group
	httpLogger := logger.WithGroup("http")
	httpLogger.Info("HTTP request received",
		"method", "GET",
		"path", "/api/users",
		"status", 200)

	// Another group
	dbLogger := logger.WithGroup("database")
	dbLogger.Info("Database query",
		"query", "SELECT * FROM users",
		"rows", 10,
		"duration_ms", 25)

	// Using context and attributes with groups
	ctx := context.Background()
	httpLogger.InfoContext(ctx, "Request processed with context",
		"user_agent", "Mozilla/5.0",
		"referer", "https://example.com")
}

// demonstrateNestedGroups shows nested group usage
func demonstrateNestedGroups() {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(os.Stdout, opts)

	// Create a nested logger structure
	apiLogger := logger.WithGroup("api")
	userApiLogger := apiLogger.WithGroup("users")

	// Log with the deeply nested logger
	userApiLogger.Info("User created",
		"user_id", 1001,
		"email", "user@example.com")

	// Add more nesting
	adminApiLogger := userApiLogger.WithGroup("admin")
	adminApiLogger.Info("Admin action",
		"action", "user_delete",
		"target_id", 2002)

	// With attributes at each level
	authLogger := logger.WithGroup("auth").With("service", "oauth")
	tokenLogger := authLogger.WithGroup("token").With("type", "jwt")
	tokenLogger.Info("Token refreshed",
		"expires_in", 3600,
		"scope", "read write")
}

// demonstrateConcurrentUsage shows how the logger behaves in concurrent environment
func demonstrateConcurrentUsage() {
	opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
	logger := grovelog.NewLogger(os.Stdout, opts)

	// Create a common base logger
	baseLogger := logger.With("app", "demo")

	// Create a wait group to wait for all goroutines
	var wg sync.WaitGroup

	// Launch multiple goroutines with different logger configurations
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Each goroutine has its own logger with unique ID
			routineLogger := baseLogger.With(
				"goroutine_id", id,
				"timestamp", time.Now().UnixNano(),
			)

			// Simulate some work and logging
			time.Sleep(time.Duration(100*id) * time.Millisecond)

			// Create a group for this routine
			routineGroupLogger := routineLogger.WithGroup(fmt.Sprintf("routine_%d", id))

			// Log from each goroutine
			routineGroupLogger.Info("Goroutine starting work")

			// More work
			time.Sleep(time.Duration(50*id) * time.Millisecond)

			// Log completion
			routineGroupLogger.Info("Goroutine completed",
				"duration_ms", id*150,
				"status", "success")
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	logger.Info("All goroutines completed")
}

// simulated HTTP request handler
func handleRequest(logger *slog.Logger, requestID, method, path string) {
	// Create a request-scoped logger
	reqLogger := logger.With(
		"request_id", requestID,
		"method", method,
		"path", path,
	)

	reqLogger.Info("Request received")

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Log with database group
	dbLogger := reqLogger.WithGroup("db")
	dbLogger.Info("Executing database query",
		"table", "users",
		"operation", "select")

	// Simulate more processing
	time.Sleep(50 * time.Millisecond)

	// Log response
	reqLogger.Info("Request completed",
		"status", 200,
		"duration_ms", 150)
}

// A practical example showing structured logging in a typical application
func practicalExample() {
	// Create application logger with custom options
	opts := grovelog.NewOptions(slog.LevelInfo, "[2006-01-02 15:04:05.000]", grovelog.Color)
	appLogger := grovelog.NewLogger(os.Stdout, opts)

	// Add application-wide attributes
	logger := appLogger.With(
		"app", "api-server",
		"version", "1.2.3",
		"env", "development",
	)

	// Log application startup
	logger.Info("Application starting")

	// Create a server group
	serverLogger := logger.WithGroup("server")
	serverLogger.Info("Server listening", "address", ":8080")

	// Handle some simulated requests
	handleRequest(serverLogger, "req-001", "GET", "/api/users")
	handleRequest(serverLogger, "req-002", "POST", "/api/login")
	handleRequest(serverLogger, "req-003", "PUT", "/api/users/42")

	// Log application shutdown
	logger.Info("Application shutting down gracefully")
}
