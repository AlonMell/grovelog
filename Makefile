.PHONY: build test bench clean lint coverage example

# Build settings
GO=go
GOFLAGS=-trimpath
LDFLAGS=-s -w
MODULE=github.com/AlonMell/grovelog
EXAMPLE=./example

build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" ./...

test:
	$(GO) test -v ./...

bench:
	$(GO) test -bench=. -benchmem ./...

clean:
	$(GO) clean
	rm -f coverage.out

lint:
	golangci-lint run

coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

example:
	$(GO) run $(EXAMPLE)/main.go

all: clean lint test coverage build

# Help
help:
	@echo "Available targets:"
	@echo "  build      - Build the package"
	@echo "  test       - Run tests"
	@echo "  bench      - Run benchmarks"
	@echo "  clean      - Clean build artifacts"
	@echo "  lint       - Run linter"
	@echo "  coverage   - Generate test coverage report"
	@echo "  example    - Run example application"
	@echo "  all        - Run clean, lint, test, coverage, and build"
	@echo "  help       - Show this help message"