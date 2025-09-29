.PHONY: build run clean test deps examples help docker-build docker-run docker-stop docker-logs docker-dev docker-prod docker-clean

# Build the DiscordBotForge framework
build:
	go build -o bin/discord-bot-forge ./cmd/main.go

# Run the simple example bot
run-simple:
	cd examples/simple_bot && go run main.go

# Run the advanced example bot
run-advanced:
	cd examples/advanced_bot && go run main.go

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f *.log
	rm -f discord-bot-forge.log

# Run tests
test:
	go test ./...

# Install dependencies
deps:
	go mod tidy
	go mod download

# Build all examples
examples:
	cd examples/simple_bot && go build -o ../../bin/simple-bot
	cd examples/advanced_bot && go build -o ../../bin/advanced-bot

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Setup development environment
setup: deps
	@echo "🔥 DiscordBotForge development environment setup complete!"
	@echo "📝 Don't forget to:"
	@echo "   1. Copy env.example to .env"
	@echo "   2. Add your Discord bot token to .env"
	@echo "   3. Run 'make run-simple' to test your bot"

# Docker commands
docker-build:
	docker build -t discord-bot-forge:latest .

docker-build-dev:
	docker build -f Dockerfile.dev -t discord-bot-forge:dev .

docker-run:
	docker-compose up -d

docker-run-dev:
	docker-compose -f docker-compose.dev.yml up -d

docker-stop:
	docker-compose down

docker-stop-dev:
	docker-compose -f docker-compose.dev.yml down

docker-logs:
	docker-compose logs -f discord-bot-forge

docker-logs-dev:
	docker-compose -f docker-compose.dev.yml logs -f discord-bot-forge-dev

docker-shell:
	docker-compose exec discord-bot-forge sh

docker-shell-dev:
	docker-compose -f docker-compose.dev.yml exec discord-bot-forge-dev sh

docker-clean:
	docker-compose down -v
	docker system prune -f

docker-clean-dev:
	docker-compose -f docker-compose.dev.yml down -v
	docker system prune -f

# Production deployment
docker-prod: docker-build
	docker-compose up -d

# Development deployment
docker-dev: docker-build-dev
	docker-compose -f docker-compose.dev.yml up -d

# Quick start with Docker
docker-quickstart:
	@echo "🚀 Starting DiscordBotForge with Docker..."
	@echo "📝 Make sure to:"
	@echo "   1. Copy env.docker to .env"
	@echo "   2. Add your Discord bot token to .env"
	@echo "   3. Run 'make docker-prod'"
	@if [ ! -f .env ]; then \
		echo "⚠️  .env file not found. Copying env.docker to .env..."; \
		cp env.docker .env; \
		echo "✅ Please edit .env with your Discord bot token"; \
	fi

# Help
help:
	@echo "🔥 DiscordBotForge - Available targets:"
	@echo ""
	@echo "📦 Build & Run:"
	@echo "  build        - Build the framework"
	@echo "  run-simple   - Run the simple example bot"
	@echo "  run-advanced - Run the advanced example bot"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  deps         - Install dependencies"
	@echo "  examples     - Build all examples"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  setup        - Setup development environment"
	@echo ""
	@echo "🐳 Docker Commands:"
	@echo "  docker-build     - Build production Docker image"
	@echo "  docker-build-dev - Build development Docker image"
	@echo "  docker-run       - Run production stack with Docker Compose"
	@echo "  docker-run-dev   - Run development stack with Docker Compose"
	@echo "  docker-stop      - Stop production stack"
	@echo "  docker-stop-dev  - Stop development stack"
	@echo "  docker-logs      - View production logs"
	@echo "  docker-logs-dev  - View development logs"
	@echo "  docker-shell     - Access production container shell"
	@echo "  docker-shell-dev - Access development container shell"
	@echo "  docker-clean     - Clean production containers and volumes"
	@echo "  docker-clean-dev - Clean development containers and volumes"
	@echo "  docker-prod      - Deploy production stack"
	@echo "  docker-dev       - Deploy development stack"
	@echo "  docker-quickstart - Quick start with Docker setup"
