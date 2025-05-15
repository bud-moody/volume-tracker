# Makefile for volume-tracker project

# Variables
APP_NAME=volume-tracker

# Default target
.PHONY: all
all: run

# Build the application
.PHONY: build
build:
	go build -o $(APP_NAME)

# Run the application
.PHONY: run
run: build
	./$(APP_NAME)

# Run tests
.PHONY: test
test:
	go test ./...
