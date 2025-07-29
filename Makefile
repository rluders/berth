GO := go
APP_NAME := berth
APP_PATH := ./cmd/berth

.PHONY: all build run clean test lint help

all: build

build:
	@echo "Building $(APP_NAME)..."
	$(GO) build -o $(APP_NAME) $(APP_PATH)
	@echo "$(APP_NAME) built successfully."

run:
	@echo "Running $(APP_NAME)..."
	$(GO) run $(APP_PATH)

generate-mocks:
	go install github.com/vektra/mockery/v3@latest
	export MOCKERY_IN_PACKAGE=true
	mockery --config .mockery.yaml

clean:
	@echo "Cleaning up..."
	$(GO) clean
	rm -f $(APP_NAME)
	@echo "Cleanup complete."

test:
	@echo "Running tests..."
	$(GO) test ./...

lint:
	@echo "Running golangci-lint..."
	# Install golangci-lint if you haven't already:
	# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
		golangci-lint run ./...

help:
	@echo "Usage: make <command>"
	@echo "\nCommands:"
	@echo "  all    : Builds the application (default)"
	@echo "  build  : Compiles the application binary"
	@echo "  run    : Runs the application"
	@echo "  clean  : Removes build artifacts and the application binary"
	@echo "  test   : Runs all tests"
	@echo "  lint   : Runs go fmt and go vet"
	@echo "  help   : Displays this help message"
