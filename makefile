BINARY_NAME=ping-dashboard
PACKAGE_DIR=./cmd/main.go

# Default target when you just run 'make'
.PHONY: all
all: build

## Build the Go application
.PHONY: build
build:
	@echo "Building Go application..."
	go build -o $(BINARY_NAME) $(PACKAGE_DIR)

## Run the Go application
.PHONY: run
run:
	@echo "Running the application..."
	sudo ./$(BINARY_NAME)

## Build and then run
.PHONY: run-dev
run-dev: build run

## Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up build artifacts..."
	go clean
	rm -f $(BINARY_NAME)

## Install project dependencies
.PHONY: install-deps
install-deps:
	@echo "Installing Go dependencies..."
	go mod tidy

## Format Go code
.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	go fmt ./...