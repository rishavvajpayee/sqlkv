# Makefile for SQLKV project

# Variables
APP_NAME = sqlkv
GO_FILES = $(shell find . -name '*.go')
DB_FILE = sqlkv.db

# Default target
all: build

# Build the application
build:
	CGO_ENABLED=1 go build -o sqlkv
	codesign --options runtime --timestamp -s - ./sqlkv

# Run the application
run: build
	./$(APP_NAME)

# Run tests
test:
	go test ./...

# Clean up build artifacts
clean:
	rm -f $(APP_NAME) $(DB_FILE)

# Install dependencies
deps:
	go mod tidy

# Help command
help:
	@echo "Makefile commands:"
	@echo "  make build   - Build the application"
	@echo "  make run     - Run the application"
	@echo "  make test    - Run tests"
	@echo "  make clean   - Clean up build artifacts"
	@echo "  make deps    - Install dependencies"
	@echo "  make help    - Show this help message"