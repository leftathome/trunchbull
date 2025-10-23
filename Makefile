.PHONY: help build run test clean docker-build docker-up docker-down docker-logs setup

# Default target
help:
	@echo "Trunchbull - Student Dashboard"
	@echo ""
	@echo "Available commands:"
	@echo "  make setup         - Initial setup (create directories, copy config)"
	@echo "  make build         - Build the Go binary"
	@echo "  make run           - Run the application locally"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start Docker containers"
	@echo "  make docker-down   - Stop Docker containers"
	@echo "  make docker-logs   - View Docker logs"
	@echo "  make frontend      - Setup frontend (future)"

# Initial setup
setup:
	@echo "Setting up Trunchbull..."
	@mkdir -p data config logs
	@if [ ! -f .env ]; then cp .env.example .env; echo "Created .env file - please edit with your credentials"; fi
	@if [ ! -f config/config.yaml ]; then cp config/config.example.yaml config/config.yaml; echo "Created config.yaml"; fi
	@echo "Setup complete! Please edit .env with your API credentials."

# Build the Go binary
build:
	@echo "Building Trunchbull..."
	@CGO_ENABLED=1 go build -o trunchbull ./cmd/server
	@echo "Build complete: ./trunchbull"

# Run locally
run:
	@echo "Running Trunchbull..."
	@go run ./cmd/server/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -f trunchbull
	@rm -f cmd/server/server
	@rm -rf data/*.db
	@echo "Clean complete"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Trunchbull is running at http://localhost:8080"

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f

# Development
dev:
	@echo "Starting development mode with hot reload..."
	@# TODO: Add air or similar for hot reload
	@go run ./cmd/server/main.go

# Database commands
db-reset:
	@echo "Resetting database..."
	@rm -f data/trunchbull.db
	@echo "Database reset complete"

# Lint and format
lint:
	@echo "Running linters..."
	@go fmt ./...
	@go vet ./...
	@echo "Linting complete"

# Dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

# Frontend (placeholder for future)
frontend:
	@echo "Frontend setup coming soon..."
