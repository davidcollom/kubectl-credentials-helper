GO111MODULE=on
TAGS=unit

.PHONY: all build test test-verbose test-coverage clean deps lint fmt vet

all: test build

# Build the binary
build:
	go build -o bin/kubectl-credentials-keychain ./

# Run tests
test:
	go test -tags=$(TAGS) ./...

# Run tests with verbose output
test-verbose:
	go test -v -tags=$(TAGS) ./...

# Run tests with coverage
test-coverage:
	go test -v -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean test cache and built binaries
clean:
	go clean -testcache
	rm -rf bin/
	rm -f coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Run linter (requires golangci-lint to be installed)
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run all checks
check: fmt vet lint test

# Install the binary
install: build
	cp bin/kubectl-credentials-helper $(GOPATH)/bin/

# Development helpers
dev-setup: deps
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
