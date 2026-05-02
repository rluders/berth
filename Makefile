GO           := go
APP_NAME     := berth
APP_PATH     := ./cmd/berth
DOCKER_IMAGE := berth-dev
DOCKER_RUN   := docker run --rm -v $(shell pwd):/app -w /app $(DOCKER_IMAGE)

.PHONY: all build run clean test lint help docker-image docker-build docker-test docker-lint

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

test: generate-mocks
	@echo "Running tests..."
	$(GO) test ./...

lint:
	@echo "Running golangci-lint..."
	# Install golangci-lint if you haven't already:
	# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
		golangci-lint run ./...

docker-image:
	@echo "Building Docker dev image (Go 1.25)..."
	docker build -t $(DOCKER_IMAGE) -f Dockerfile.dev .

docker-build: docker-image
	@echo "Building $(APP_NAME) in Docker..."
	$(DOCKER_RUN) go build -o $(APP_NAME) $(APP_PATH)

docker-test: docker-image
	@echo "Running tests in Docker..."
	$(DOCKER_RUN) sh -c "mockery --config .mockery.yaml && go test ./..."

docker-lint: docker-image
	@echo "Running lint in Docker..."
	$(DOCKER_RUN) sh -c "mockery --config .mockery.yaml && golangci-lint run --timeout=5m ./..."

help:
	@echo "Usage: make <command>"
	@echo "\nCommands:"
	@echo "  all          : Builds the application (default)"
	@echo "  build        : Compiles the application binary"
	@echo "  run          : Runs the application"
	@echo "  clean        : Removes build artifacts and the application binary"
	@echo "  test         : Runs all tests"
	@echo "  lint         : Runs golangci-lint"
	@echo "  docker-image : Builds Docker dev image (Go 1.25)"
	@echo "  docker-build : Builds the binary inside Docker"
	@echo "  docker-test  : Runs tests inside Docker"
	@echo "  docker-lint  : Runs lint inside Docker"
	@echo "  help         : Displays this help message"
