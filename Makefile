# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Build parameters
BINARY_NAME=fiber-gorm-starter
BINARY_UNIX=$(BINARY_NAME)_unix
MAIN_PATH=./cmd/server/main.go

# Docker parameters
DOCKER_IMAGE=fiber-gorm-starter
DOCKER_TAG=latest

# Database parameters
DB_URL_MYSQL=mysql://user:password@tcp(localhost:3306)/fiber_gorm_starter
DB_URL_POSTGRES=postgresql://user:password@localhost:5432/fiber_gorm_starter?sslmode=disable

.PHONY: all build clean test coverage deps tidy fmt vet lint run help
.PHONY: docker-build docker-run docker-up docker-down docker-up-pg
.PHONY: migrate-up migrate-down migrate-status migrate-create seed swag
.PHONY: dev prod check install-tools

# Default target
all: check build

## Build targets

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_PATH)

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(MAIN_PATH)

# Build with optimizations for production
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-w -s" -o $(BINARY_NAME) -v $(MAIN_PATH)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out

## Development targets

# Run the application
run:
	@echo "Running application..."
	$(GOCMD) run $(MAIN_PATH)

# Run in development mode with auto-reload (requires air)
dev:
	@echo "Starting development server..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to normal run..."; \
		$(MAKE) run; \
	fi

## Testing targets

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

## Code quality targets

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -w .

# Vet code
vet:
	@echo "Vetting code..."
	$(GOVET) ./...

# Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run --timeout=5m; \
	else \
		echo "golangci-lint not found. Please install it first."; \
		echo "Installation: https://golangci-lint.run/usage/install/"; \
		exit 1; \
	fi

# Run all checks (format, vet, lint, test, build)
check: fmt vet lint test build
	@echo "All checks passed!"

## Dependency management

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

# Update dependencies
update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(MAKE) tidy

## Docker targets

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Start services with MySQL
docker-up:
	@echo "Starting services with MySQL..."
	docker-compose --profile mysql --profile app up -d

# Start services with PostgreSQL
docker-up-pg:
	@echo "Starting services with PostgreSQL..."
	docker-compose --profile postgres --profile app up -d

# Stop all services
docker-down:
	@echo "Stopping all services..."
	docker-compose down

# View logs
docker-logs:
	@echo "Viewing application logs..."
	docker-compose logs -f app

## Database targets

# Run database migrations up
migrate-up:
	@echo "Running database migrations up..."
	./scripts/migrate.sh up

# Run database migrations down
migrate-down:
	@echo "Running database migrations down..."
	./scripts/migrate.sh down 1

# Show migration status
migrate-status:
	@echo "Checking migration status..."
	./scripts/migrate.sh status

# Create new migration
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	./scripts/migrate.sh create $(name)

# Seed database
seed:
	@echo "Seeding database..."
	$(GOCMD) run scripts/seed.go

## Documentation targets

# Generate Swagger docs
swag:
	@echo "Generating Swagger documentation..."
	@if command -v swag > /dev/null; then \
		swag init -g $(MAIN_PATH) -o ./docs; \
		echo "Swagger docs generated in ./docs/"; \
	else \
		echo "swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"; \
		exit 1; \
	fi

## Performance testing

# Run k6 smoke test
k6-smoke:
	@echo "Running k6 smoke test..."
	@if command -v k6 > /dev/null; then \
		k6 run scripts/k6/users-smoke.js; \
	else \
		echo "k6 not found. Please install it first."; \
		echo "Installation: https://k6.io/docs/getting-started/installation/"; \
		exit 1; \
	fi

# Run k6 load test
k6-load:
	@echo "Running k6 load test..."
	@if command -v k6 > /dev/null; then \
		k6 run --vus 100 --duration 2m scripts/k6/users-smoke.js; \
	else \
		echo "k6 not found. Please install it first."; \
		exit 1; \
	fi

## Tool installation

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOCMD) install github.com/cosmtrek/air@latest
	$(GOCMD) install github.com/swaggo/swag/cmd/swag@latest
	@echo "Installing golangci-lint..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2
	@echo "Tools installed successfully!"

## Environment setup

# Setup for MySQL development
setup-mysql: docker-up migrate-up seed
	@echo "MySQL environment setup complete!"
	@echo "Application: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/docs/index.html"

# Setup for PostgreSQL development
setup-postgres: docker-up-pg migrate-up seed
	@echo "PostgreSQL environment setup complete!"
	@echo "Application: http://localhost:8080"
	@echo "Swagger: http://localhost:8080/docs/index.html"

## Production targets

# Production build with all optimizations
prod: clean fmt vet lint test build-prod
	@echo "Production build completed!"

# Deploy (placeholder for actual deployment script)
deploy:
	@echo "Deploy target - implement your deployment logic here"
	@echo "Example: deploy to Kubernetes, AWS, etc."

## Help

# Show help
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  build        - Build the application"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-prod   - Build with production optimizations"
	@echo "  clean        - Clean build artifacts"
	@echo ""
	@echo "Development targets:"
	@echo "  run          - Run the application"
	@echo "  dev          - Run with auto-reload (requires air)"
	@echo ""
	@echo "Testing targets:"
	@echo "  test         - Run tests"
	@echo "  coverage     - Run tests with coverage report"
	@echo "  bench        - Run benchmarks"
	@echo ""
	@echo "Code quality targets:"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Run linter"
	@echo "  check        - Run all checks (fmt, vet, lint, test, build)"
	@echo ""
	@echo "Dependency management:"
	@echo "  deps         - Download dependencies"
	@echo "  tidy         - Tidy dependencies"
	@echo "  update       - Update dependencies"
	@echo ""
	@echo "Docker targets:"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-up    - Start services with MySQL"
	@echo "  docker-up-pg - Start services with PostgreSQL"
	@echo "  docker-down  - Stop all services"
	@echo "  docker-logs  - View application logs"
	@echo ""
	@echo "Database targets:"
	@echo "  migrate-up   - Run database migrations up"
	@echo "  migrate-down - Run database migrations down"
	@echo "  migrate-status - Show migration status"
	@echo "  migrate-create name=<name> - Create new migration"
	@echo "  seed         - Seed database"
	@echo ""
	@echo "Documentation targets:"
	@echo "  swag         - Generate Swagger documentation"
	@echo ""
	@echo "Performance testing:"
	@echo "  k6-smoke     - Run k6 smoke test"
	@echo "  k6-load      - Run k6 load test"
	@echo ""
	@echo "Environment setup:"
	@echo "  setup-mysql  - Complete MySQL development setup"
	@echo "  setup-postgres - Complete PostgreSQL development setup"
	@echo "  install-tools - Install development tools"
	@echo ""
	@echo "Production targets:"
	@echo "  prod         - Production build with all checks"
	@echo "  deploy       - Deploy (customize for your environment)"
	@echo ""
	@echo "Examples:"
	@echo "  make setup-mysql     # Quick MySQL setup"
	@echo "  make dev             # Start development server"
	@echo "  make check           # Run all quality checks"
	@echo "  make migrate-create name=add_user_profile"
	@echo "  make k6-smoke        # Run performance test"