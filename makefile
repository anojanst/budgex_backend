# Budgex Backend Makefile

.PHONY: help run build test clean swagger dev stop logs

# Default target
help:
	@echo "Available commands:"
	@echo "  run        - Run the server in development mode"
	@echo "  build      - Build the server binary"
	@echo "  test       - Run tests"
	@echo "  test-cover - Run tests with coverage"
	@echo "  swagger    - Generate Swagger documentation"
	@echo "  dev        - Run with hot reload (requires air)"
	@echo "  stop       - Stop running server instances"
	@echo "  clean      - Clean build artifacts"
	@echo "  tidy       - Tidy Go modules"
	@echo "  lint       - Run linter"
	@echo "  logs       - View server logs"

# Development
run:
	@echo "Starting Budgex backend server..."
	go run ./cmd/server

dev:
	@echo "Starting with hot reload..."
	air

# Build
build:
	@echo "Building server binary..."
	go build -o bin/server cmd/server/main.go
	@echo "Binary built: bin/server"

build-linux:
	@echo "Building Linux binary..."
	GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go
	@echo "Linux binary built: bin/server-linux"

# Testing
test:
	@echo "Running tests..."
	go test ./...

test-cover:
	@echo "Running tests with coverage..."
	go test -cover ./...

test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

# Documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go -o internal/docs
	@echo "Swagger docs generated in internal/docs/"

swagger-serve:
	@echo "Starting server and opening Swagger UI..."
	@echo "Swagger UI will be available at: http://localhost:8080/swagger/index.html"
	go run ./cmd/server

# Database
db-reset:
	@echo "Resetting database..."
	dropdb budgex || true
	createdb budgex
	@echo "Database reset complete"

db-migrate:
	@echo "Running database migrations..."
	go run ./cmd/server &
	@sleep 3
	@echo "Migrations complete"
	pkill -f "go run ./cmd/server" || true

# Utilities
stop:
	@echo "Stopping server instances..."
	pkill -f "go run cmd/server/main.go" || true
	pkill -f "go run ./cmd/server" || true
	@echo "Server stopped"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf internal/docs/
	@echo "Clean complete"

tidy:
	@echo "Tidying Go modules..."
	go mod tidy
	@echo "Modules tidied"

lint:
	@echo "Running linter..."
	golangci-lint run || true

# Health check
health:
	@echo "Checking server health..."
	@curl -s http://localhost:8080/api/healthz | jq . || echo "Server not responding"

# API testing
test-auth:
	@echo "Testing authentication endpoint..."
	@curl -s http://localhost:8080/api/me -H "Authorization: Bearer YOUR_TOKEN_HERE" | jq . || echo "Authentication test failed"

# Development setup
setup:
	@echo "Setting up development environment..."
	@echo "Installing dependencies..."
	go mod tidy
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Generating Swagger docs..."
	$(MAKE) swagger
	@echo "Setup complete!"
	@echo "Run 'make run' to start the server"
