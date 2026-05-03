GO           := go
APP_NAME     := berth
APP_PATH     := ./cmd/berth
PODMAN_IMAGE := berth-dev
PODMAN_RUN   := podman run --rm -v $(shell pwd):/app:Z -w /app $(PODMAN_IMAGE)

.PHONY: all build run clean test lint help podman-image podman-build podman-test podman-lint docker-image docker-build docker-test docker-lint

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

podman-image:
	@echo "Building Podman dev image (Go 1.26)..."
	podman build -t $(PODMAN_IMAGE) -f Dockerfile.dev .

podman-build: podman-image
	@echo "Building $(APP_NAME) in Podman..."
	$(PODMAN_RUN) go build -o $(APP_NAME) $(APP_PATH)

podman-test: podman-image
	@echo "Running tests in Podman..."
	$(PODMAN_RUN) sh -c "mockery --config .mockery.yaml && go test ./..."

podman-lint: podman-image
	@echo "Running lint in Podman..."
	$(PODMAN_RUN) sh -c "mockery --config .mockery.yaml && golangci-lint run --timeout=5m ./..."

docker-image: podman-image
docker-build: podman-build
docker-test: podman-test
docker-lint: podman-lint

help:
	@echo "Usage: make <command>"
	@echo "\nCommands:"
	@echo "  all          : Builds the application (default)"
	@echo "  build        : Compiles the application binary"
	@echo "  run          : Runs the application"
	@echo "  clean        : Removes build artifacts and the application binary"
	@echo "  test         : Runs all tests"
	@echo "  lint         : Runs golangci-lint"
	@echo "  podman-image : Builds Podman dev image (Go 1.26)"
	@echo "  podman-build : Builds the binary inside Podman"
	@echo "  podman-test  : Runs tests inside Podman"
	@echo "  podman-lint  : Runs lint inside Podman"
	@echo "  help         : Displays this help message"
