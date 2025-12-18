.PHONY: help install up down logs test bench clean ui api api-bg api-stop worker infra-up infra-down

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
	@echo "${YELLOW}Starting all services...${RESET}"
	@docker-compose up -d --build
	@echo "${GREEN}✓ Services started!${RESET}"
	@echo "${BLUE}API:${RESET} http://localhost:8080"
	@echo "${BLUE}UI:${RESET}  http://localhost:3000"

down: ## Stop all services
	@echo "${YELLOW}Stopping all services...${RESET}"
	@docker-compose down
	@echo "${GREEN}✓ Services stopped!${RESET}"

logs: ## View service logs
	@docker-compose logs -f

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
	@rm -rf src/ui/node_modules src/ui/.svelte-kit
	@docker-compose down -v
	@go clean -testcache -cache
	@echo "${GREEN}✓ Cleanup completed!${RESET}"
