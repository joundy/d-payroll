# Dealls Payroll System Makefile

.PHONY: help run seed build clean test docker-up docker-down

# Default target
help:
	@echo "Available commands:"
	@echo "  run        - Run the main application"
	@echo "  seed       - Run database seeding"
	@echo "  build      - Build the application"
	@echo "  test       - Run tests"
	@echo "  docker-up  - Start Docker services"
	@echo "  docker-down- Stop Docker services"

# Run the main application
run:
	@echo "Starting Dealls Payroll System..."
	go run cmd/app/main.go

# Run database seeding
seed:
	@echo "Seeding database..."
	go run cmd/seed/main.go

# Build the application
build:
	@echo "Building application..."
	go build -o bin/app cmd/app/main.go
	go build -o bin/seed cmd/seed/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./tests/integration


# Start Docker services (database)
docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d

# Stop Docker services
docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy