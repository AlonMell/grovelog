# GroveLog

GroveLog is a flexible logging library for Go, built on top of the standard `log/slog` package. It provides enhanced formatting options, color support, group handling, context-aware logging, and optimized performance.

[![Go Version](https://img.shields.io/github/go-mod/go-version/AlonMell/grovelog)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/AlonMell/grovelog.svg)](https://pkg.go.dev/github.com/AlonMell/grovelog)
[![Go Report Card](https://goreportcard.com/badge/github.com/AlonMell/grovelog)](https://goreportcard.com/report/github.com/AlonMell/grovelog)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- Multiple output formats: JSON, Plain Text, and Colored Text
- Thread-safe for concurrent use
- Efficient memory usage with buffer pooling
- Support for structured logging with attributes
- Advanced grouping of attributes with nesting
- Context-aware logging with attribute propagation
- Error wrapping with context preservation
- Compatible with the standard `log/slog` interface
- Customizable time formats

## Installation

```bash
go get github.com/AlonMell/grovelog
```

## Quick Start

```go
package main

import (
    "os"
    "log/slog"

    "github.com/AlonMell/grovelog"
)

func main() {
    // Create a color logger with INFO level
    opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
    logger := grovelog.NewLogger(os.Stdout, opts)

    // Simple logging
    logger.Info("Hello, GroveLog!")

    // With attributes
    logger.Info("User logged in",
        "user_id", 1234,
        "source", "api")

    // With groups
    dbLogger := logger.WithGroup("database")
    dbLogger.Info("Query executed",
        "query", "SELECT * FROM users",
        "duration_ms", 42)
}
```

## Output Formats

GroveLog supports three output formats:

### JSON Format

```go
opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.JSON)
logger := grovelog.NewLogger(os.Stdout, opts)
logger.Info("Hello JSON", "key", "value")
```

Output:
```json
{"time":"2025-04-07T10:30:45.123456789Z","level":"INFO","msg":"Hello JSON","key":"value"}
```

### Plain Text Format

```go
opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Plain)
logger := grovelog.NewLogger(os.Stdout, opts)
logger.Info("Hello Plain", "key", "value")
```

Output:
```
time=2025-04-07T10:30:45.123456789Z level=INFO msg="Hello Plain" key=value
```

### Color Format

```go
opts := grovelog.NewOptions(slog.LevelInfo, "", grovelog.Color)
logger := grovelog.NewLogger(os.Stdout, opts)
logger.Info("Hello Color", "key", "value")
```

Output:
```
[10:30:45.123] INFO: Hello Color {"key":"value"}
```

## Advanced Usage

### Custom Time Format

```go
opts := grovelog.NewOptions(slog.LevelInfo, "2006-01-02 15:04:05", grovelog.Color)
logger := grovelog.NewLogger(os.Stdout, opts)
```

### Nested Groups

```go
apiLogger := logger.WithGroup("api")
userLogger := apiLogger.WithGroup("users")
userLogger.Info("User created", "id", 1001, "email", "user@example.com")
```

Output:
```
[10:30:45.123] INFO: User created {"api.users.id":1001,"api.users.email":"user@example.com"}
```

### Group and Attributes

```go
authLogger := logger.WithGroup("auth").With("service", "oauth")
authLogger.Info("Token generated", "expires_in", 3600)
```

## Performance

GroveLog is designed with performance in mind:

- Buffer pooling to reduce memory allocations
- Efficient JSON serialization
- Minimal lock contention for thread safety

## Development

### Testing, Linting, Coverage, Benchmarks

```bash
make help
```

## License

MIT License - see the [LICENSE](LICENSE) file for details.