.PHONY: help install up down logs test bench clean clean-volumes reset-db fix ui api api-bg api-stop worker infra-up infra-down lint build build-fast

# Enable BuildKit for faster builds with better caching
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1

# Colors
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
BLUE   := $(shell tput -Txterm setaf 4)
RESET  := $(shell tput -Txterm sgr0)

help: ## Show this help message
	@echo '${BLUE}Available commands:${RESET}'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-15s${RESET} %s\n", $$1, $$2}'

install: ## Install all dependencies (Go + npm)
	@echo "${YELLOW}Installing Go dependencies...${RESET}"
	@go mod download
	@echo "${YELLOW}Installing UI dependencies (npm)...${RESET}"
	@cd src/ui && npm install
	@echo "${GREEN}✓ All dependencies installed!${RESET}"

up: ## Start all services in Docker (API + Worker + UI + infra)
	@echo "${YELLOW}Starting all services with optimized builds...${RESET}"
	@docker-compose build --parallel
	@docker-compose up -d
	@echo "${GREEN}✓ Services started!${RESET}"
	@echo "${BLUE}API:${RESET} http://localhost:8080"
	@echo "${BLUE}UI:${RESET}  http://localhost:3000"
	@make logs

build: ## Build all Docker images in parallel (optimized)
	@echo "${YELLOW}Building all images in parallel with BuildKit...${RESET}"
	@docker-compose build --parallel
	@echo "${GREEN}✓ Build completed!${RESET}"

build-fast: ## Build only changed services (fast incremental build)
	@echo "${YELLOW}Building changed services only...${RESET}"
	@docker-compose build
	@echo "${GREEN}✓ Fast build completed!${RESET}"

down: ## Stop all services
	@echo "${YELLOW}Stopping all services...${RESET}"
	@docker-compose down
	@echo "${GREEN}✓ Services stopped!${RESET}"

logs: ## View service logs
	@docker-compose logs -f

lint: ## Run linters (golangci-lint)
	@echo "${YELLOW}Running linters...${RESET}"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "${YELLOW}golangci-lint not found. Installing...${RESET}"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest; \
	fi
	@golangci-lint run ./src/...
	@echo "${GREEN}✓ Lint completed!${RESET}"

test: ## Run all tests in parallel with race detector
	@echo "${YELLOW}Running all tests with race detector...${RESET}"
	@go test -race -v -parallel 4 ./src/...
	@echo "${GREEN}✓ Tests completed!${RESET}"

bench: ## Run rules engine benchmark
	@echo "${YELLOW}Running rules engine benchmark...${RESET}"
	@go test -bench=. -benchmem -benchtime=5s -run=^$$ ./src/workers/internal/rules/...
	@echo "${GREEN}✓ Benchmark completed!${RESET}"

clean: ## Remove build artifacts and Docker volumes
	@echo "${YELLOW}Cleaning artifacts...${RESET}"
	@rm -rf bin/ coverage.out coverage.html
	@rm -rf src/ui/node_modules src/ui/.next
	@docker-compose down -v
	@go clean -testcache -cache
	@echo "${GREEN}✓ Cleanup completed!${RESET}"

ui: ## Start UI service only
	@echo "${YELLOW}Building and starting UI service...${RESET}"
	@docker-compose build ui
	@docker-compose up -d ui
	@echo "${GREEN}✓ UI service started!${RESET}"
	@echo "${BLUE}UI:${RESET}  http://localhost:3000"

api: ## Start API service with infrastructure (postgres + redis)
	@echo "${YELLOW}Starting infrastructure services (postgres + redis)...${RESET}"
	@docker-compose up -d postgres redis
	@echo "${YELLOW}Building API service...${RESET}"
	@docker-compose build api
	@echo "${YELLOW}Waiting for infrastructure to be healthy...${RESET}"
	@docker-compose up -d api
	@echo "${GREEN}✓ API service with infrastructure started!${RESET}"
	@echo "${BLUE}API:${RESET} http://localhost:8080"
	@make logs

worker: ## Start Worker service with infrastructure (postgres + redis)
	@echo "${YELLOW}Starting infrastructure services (postgres + redis)...${RESET}"
	@docker-compose up -d postgres redis
	@echo "${YELLOW}Building Worker service...${RESET}"
	@docker-compose build worker
	@echo "${YELLOW}Waiting for infrastructure to be healthy...${RESET}"
	@docker-compose up -d worker
	@echo "${GREEN}✓ Worker service with infrastructure started!${RESET}"
	@make logs
