.PHONY: help ssl proto build up down logs clean restart server-logs client-logs

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

ssl: ## Generate SSL certificates
	@./scripts/generate-ssl.sh

proto: ## Generate protobuf code
	@echo "🔧 Generating protobuf code..."
	@PATH="$(shell go env GOPATH)/bin:$$PATH" protoc \
		-I./server/proto \
		--go_opt=module=github.com/abisalde/go-showcase \
		--go_out=. \
		--go-grpc_opt=module=github.com/abisalde/go-showcase \
		--go-grpc_out=. \
		./server/proto/church.proto
	@echo "✅ Protobuf code generated successfully!"

build: ## Build Docker images
	@echo "🐳 Building Docker images..."
	@docker-compose build
	@echo "✅ Build complete!"

up: ## Start all services
	@echo "🚀 Starting services..."
	@docker-compose up -d
	@echo "✅ Services started!"
	@echo ""
	@echo "📊 Traefik Dashboard: http://localhost:8081"
	@echo "🔌 gRPC Server: localhost:50052"

down: ## Stop all services
	@echo "🛑 Stopping services..."
	@docker-compose down
	@echo "✅ Services stopped!"

logs: ## Show logs from all services
	@docker-compose logs -f

server-logs: ## Show server logs
	@docker-compose logs -f grpc-server

client-logs: ## Show client logs
	@docker-compose logs -f grpc-client

restart: down up ## Restart all services

clean: down ## Clean up containers, volumes, and temporary files
	@echo "🧹 Cleaning up..."
	@docker-compose down -v
	@rm -rf tmp/
	@echo "✅ Cleanup complete!"

dev: ## Start development environment
	@echo "🔧 Starting development environment..."
	@make ssl
	@make proto
	@make up
	@make logs
