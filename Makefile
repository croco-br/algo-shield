.PHONY: help build run stop clean test dev docker-build docker-up docker-down deps-go deps-ui check-deps

# Colors for output
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

help: ## Show this help
	@echo '${YELLOW}Available commands:${RESET}'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-20s${RESET} %s\n", $$1, $$2}'

# Check and install Go dependencies
deps-go:
	@if [ ! -d "vendor" ] && [ ! -f "go.sum" ]; then \
		echo "${YELLOW}Installing Go dependencies...${RESET}"; \
		go mod download; \
	fi

# Check and install UI dependencies
deps-ui:
	@if [ ! -d "ui/node_modules" ]; then \
		echo "${YELLOW}Installing UI dependencies (npm install)...${RESET}"; \
		cd ui && npm install; \
	else \
		echo "${GREEN}UI dependencies already installed${RESET}"; \
	fi

check-deps: deps-go ## Check and install all dependencies
	@echo "${GREEN}All dependencies checked!${RESET}"

build: deps-go ## Build all binaries
	@echo "${GREEN}Building API...${RESET}"
	@GOEXPERIMENT=rangefunc go build -ldflags="-w -s" -o bin/api ./api/cmd/main.go
	@echo "${GREEN}Building Worker...${RESET}"
	@GOEXPERIMENT=rangefunc go build -ldflags="-w -s" -o bin/worker ./workers/cmd/main.go
	@echo "${GREEN}Build completed!${RESET}"

run-api: build ## Run API locally
	@echo "${GREEN}Starting API...${RESET}"
	@./bin/api

run-worker: build ## Run Worker locally
	@echo "${GREEN}Starting Worker...${RESET}"
	@./bin/worker

dev-ui: deps-ui ## Run UI in development mode
	@echo "${GREEN}Starting UI development server...${RESET}"
	@cd ui && npm run dev

test: deps-go ## Run tests
	@echo "${GREEN}Running tests...${RESET}"
	@GOEXPERIMENT=rangefunc go test -v ./...

docker-build: ## Build Docker images
	@echo "${GREEN}Building Docker images...${RESET}"
	@docker-compose build

docker-up: ## Start all services with Docker Compose
	@echo "${GREEN}Starting all services...${RESET}"
	@docker-compose up -d
	@echo "${GREEN}Services started!${RESET}"
	@echo "${YELLOW}API:    http://localhost:8080${RESET}"
	@echo "${YELLOW}UI:     http://localhost:5173${RESET}"

docker-down: ## Stop all services
	@echo "${GREEN}Stopping all services...${RESET}"
	@docker-compose down

docker-logs: ## View logs
	@docker-compose logs -f

dev-infra: ## Start infrastructure (postgres & redis) for local development
	@echo "${GREEN}Starting PostgreSQL and Redis...${RESET}"
	@docker-compose up -d postgres redis
	@echo "${GREEN}Infrastructure started!${RESET}"
	@echo "${YELLOW}PostgreSQL: localhost:5432${RESET}"
	@echo "${YELLOW}Redis:      localhost:6379${RESET}"

dev-all: deps ## Install dependencies and show instructions to run all services
	@echo "${GREEN}Dependencies installed!${RESET}"
	@echo ""
	@echo "${YELLOW}To run all services locally, open 4 terminals:${RESET}"
	@echo ""
	@echo "Terminal 1 - Infrastructure:"
	@echo "  ${GREEN}make dev-infra${RESET}"
	@echo ""
	@echo "Terminal 2 - API:"
	@echo "  ${GREEN}make run-api${RESET}"
	@echo ""
	@echo "Terminal 3 - Worker:"
	@echo "  ${GREEN}make run-worker${RESET}"
	@echo ""
	@echo "Terminal 4 - UI:"
	@echo "  ${GREEN}make dev-ui${RESET}"
	@echo ""

clean: ## Clean build artifacts
	@echo "${GREEN}Cleaning...${RESET}"
	@rm -rf bin/
	@rm -rf ui/node_modules
	@rm -rf ui/.svelte-kit
	@docker-compose down -v

clean-deps: ## Clean all dependencies
	@echo "${GREEN}Cleaning dependencies...${RESET}"
	@rm -rf ui/node_modules
	@rm -rf ui/.svelte-kit
	@go clean -modcache

deps: deps-go deps-ui ## Download and install all dependencies
	@echo "${GREEN}All dependencies installed!${RESET}"

deps-force: ## Force reinstall all dependencies
	@echo "${YELLOW}Force reinstalling dependencies...${RESET}"
	@go mod download
	@cd ui && rm -rf node_modules && npm install
	@echo "${GREEN}Dependencies reinstalled!${RESET}"

install-hooks: ## Install Git hooks for automated testing
	@echo "${GREEN}Installing Git hooks...${RESET}"
	@./scripts/install-hooks.sh

