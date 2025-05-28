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

# Docker targets
.PHONY: docker-build-backend
.PHONY: docker-build-frontend
.PHONY: docker-compose-up

# Build backend Docker image
docker-build-backend:
	docker build -f Dockerfile.backend -t volume-tracker-backend .

# Build frontend Docker image
docker-build-frontend:
	docker build -f Dockerfile.frontend -t volume-tracker-frontend .

# Build both images
.PHONY: docker-build-all
docker-build-all: docker-build-backend docker-build-frontend

# Start the full stack
.PHONY: up
docker-compose-up:
	docker-compose up --build
