.PHONY: help install up up-dev down logs test test-api test-ui test-coverage test-coverage-unit test-coverage-integration test-coverage-html bench clean clean-volumes reset-db fix ui api api-bg api-stop worker dev-api dev-worker air infra-up infra-down test-watch coverage-ci check-deps lint build build-fast

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

install: ## Install all dependencies (Go + npm + golangci-lint)
	@echo "${YELLOW}Installing Go dependencies...${RESET}"
	@go mod download
	@echo "${YELLOW}Installing UI dependencies (npm)...${RESET}"
	@cd src/ui && npm install
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "${YELLOW}golangci-lint not found. Installing...${RESET}"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest; \
	fi
	@echo "${GREEN}✓ All dependencies installed!${RESET}"

up: ## Start all services in Docker (API + Worker + UI + infra)
	@echo "${YELLOW}Starting all services with optimized builds...${RESET}"
	@DOCKER_BUILDKIT=1 docker-compose build --parallel
	@docker-compose up -d
	@echo "${GREEN}✓ Services started!${RESET}"
	@echo "${BLUE}API:${RESET} http://localhost:8080"
	@echo "${BLUE}UI:${RESET}  http://localhost:3000"
	@echo "${BLUE}Asynqmon:${RESET} http://localhost:8081"
	@make logs

up-dev: ## Start all services in development mode (faster health checks)
	@echo "${YELLOW}Starting all services in development mode...${RESET}"
	@HEALTHCHECK_INTERVAL_POSTGRES=5s \
	HEALTHCHECK_TIMEOUT_POSTGRES=2s \
	HEALTHCHECK_START_PERIOD_POSTGRES=5s \
	HEALTHCHECK_INTERVAL_REDIS=5s \
	HEALTHCHECK_TIMEOUT_REDIS=2s \
	HEALTHCHECK_START_PERIOD_REDIS=3s \
	HEALTHCHECK_INTERVAL_API=10s \
	HEALTHCHECK_TIMEOUT_API=3s \
	HEALTHCHECK_START_PERIOD_API=8s \
	HEALTHCHECK_INTERVAL_WORKER=10s \
	HEALTHCHECK_TIMEOUT_WORKER=3s \
	HEALTHCHECK_START_PERIOD_WORKER=5s \
	HEALTHCHECK_INTERVAL_UI=10s \
	HEALTHCHECK_TIMEOUT_UI=3s \
	HEALTHCHECK_START_PERIOD_UI=8s \
	HEALTHCHECK_RETRIES=3 \
	DOCKER_BUILDKIT=1 docker-compose build --parallel
	@HEALTHCHECK_INTERVAL_POSTGRES=5s \
	HEALTHCHECK_TIMEOUT_POSTGRES=2s \
	HEALTHCHECK_START_PERIOD_POSTGRES=5s \
	HEALTHCHECK_INTERVAL_REDIS=5s \
	HEALTHCHECK_TIMEOUT_REDIS=2s \
	HEALTHCHECK_START_PERIOD_REDIS=3s \
	HEALTHCHECK_INTERVAL_API=10s \
	HEALTHCHECK_TIMEOUT_API=3s \
	HEALTHCHECK_START_PERIOD_API=8s \
	HEALTHCHECK_INTERVAL_WORKER=10s \
	HEALTHCHECK_TIMEOUT_WORKER=3s \
	HEALTHCHECK_START_PERIOD_WORKER=5s \
	HEALTHCHECK_INTERVAL_UI=10s \
	HEALTHCHECK_TIMEOUT_UI=3s \
	HEALTHCHECK_START_PERIOD_UI=8s \
	HEALTHCHECK_RETRIES=3 \
	docker-compose up -d
	@echo "${GREEN}✓ Services started in dev mode!${RESET}"
	@echo "${BLUE}API:${RESET} http://localhost:8080"
	@echo "${BLUE}UI:${RESET}  http://localhost:3000"
	@echo "${BLUE}Asynqmon:${RESET} http://localhost:8081"
	@echo "${YELLOW}ℹ Dev mode: 3x faster health checks (5-10s intervals)${RESET}"
	@make logs

build: ## Build all Docker images in parallel (optimized)
	@echo "${YELLOW}Building all images in parallel with BuildKit...${RESET}"
	@DOCKER_BUILDKIT=1 docker-compose build --parallel
	@echo "${GREEN}✓ Build completed!${RESET}"

build-fast: ## Build only changed services (smart incremental build)
	@echo "${YELLOW}Building changed services only...${RESET}"
	@CHANGED_FILES=$$(git diff --name-only HEAD~1 2>/dev/null || echo "all"); \
	if echo "$$CHANGED_FILES" | grep -q "all\|go.mod\|go.sum"; then \
		echo "${YELLOW}→ Detected significant changes, building all services...${RESET}"; \
		DOCKER_BUILDKIT=1 docker-compose build --parallel; \
	else \
		if echo "$$CHANGED_FILES" | grep -q "src/api/\|Dockerfile.api"; then \
			echo "${YELLOW}→ Building API service...${RESET}"; \
			DOCKER_BUILDKIT=1 docker-compose build api; \
		fi; \
		if echo "$$CHANGED_FILES" | grep -q "src/workers/\|Dockerfile.worker"; then \
			echo "${YELLOW}→ Building Worker service...${RESET}"; \
			DOCKER_BUILDKIT=1 docker-compose build worker; \
		fi; \
		if echo "$$CHANGED_FILES" | grep -q "src/ui/\|Dockerfile.ui"; then \
			echo "${YELLOW}→ Building UI service...${RESET}"; \
			DOCKER_BUILDKIT=1 docker-compose build ui; \
		fi; \
	fi
	@echo "${GREEN}✓ Incremental build completed!${RESET}"

down: ## Stop all services
	@echo "${YELLOW}Stopping all services...${RESET}"
	@docker-compose down
	@echo "${GREEN}✓ Services stopped!${RESET}"

logs: ## View service logs
	@docker-compose logs -f

lint: install ## Run linters (golangci-lint)
	@echo "${YELLOW}Running linters...${RESET}"
	@golangci-lint run ./src/...
	@echo "${GREEN}✓ Lint completed!${RESET}"

test: install test-api test-ui ## Run all tests (API + UI)
	@echo "${GREEN}✓ All tests completed!${RESET}"

test-api: gotestsum ## Run API tests with race detector
	@echo "${YELLOW}Running API tests with race detector...${RESET}"
	@gotestsum --format testdox -- -race -parallel 8 ./src/...
	@echo "${GREEN}✓ API tests completed!${RESET}"

test-ui: ## Run UI tests with vitest
	@echo "${YELLOW}Running UI tests...${RESET}"
	@cd src/ui && npm test
	@echo "${GREEN}✓ UI tests completed!${RESET}"

test-coverage: test-coverage-unit test-coverage-integration gocovmerge ## Run all tests with combined coverage (unit + integration)
	@echo "${YELLOW}Combining coverage reports...${RESET}"
	@gocovmerge coverage-unit.out coverage-integration.out > coverage-combined.out
	@echo "${YELLOW}Coverage summary:${RESET}"
	@go tool cover -func=coverage-combined.out | tail -1
	@echo "${GREEN}✓ Combined coverage report generated: coverage-combined.out${RESET}"

test-coverage-unit: gotestsum ## Run unit tests with coverage
	@echo "${YELLOW}Running unit tests with coverage...${RESET}"
	@rm -f coverage-unit.out
	@gotestsum --format testdox -- -race -parallel 4 -coverprofile=coverage-unit.out -covermode=atomic ./src/api/... ./src/workers/... ./src/pkg/...
	@echo "${YELLOW}Unit tests coverage:${RESET}"
	@go tool cover -func=coverage-unit.out | tail -1
	@echo "${GREEN}✓ Unit tests coverage report generated: coverage-unit.out${RESET}"

test-coverage-integration: gotestsum ## Run integration tests with coverage
	@echo "${YELLOW}Running integration tests with coverage...${RESET}"
	@echo "${YELLOW}Note: This requires Docker containers (postgres, redis)${RESET}"
	@rm -f coverage-integration.out
	@gotestsum --format testdox -- -tags=integration -race -parallel 4 -coverprofile=coverage-integration.out -covermode=atomic ./src/api/... ./src/workers/...
	@echo "${YELLOW}Integration tests coverage:${RESET}"
	@go tool cover -func=coverage-integration.out | tail -1
	@echo "${GREEN}✓ Integration tests coverage report generated: coverage-integration.out${RESET}"

test-coverage-html: test-coverage ## Generate HTML coverage report
	@echo "${YELLOW}Generating HTML coverage report...${RESET}"
	@go tool cover -html=coverage-combined.out -o coverage.html
	@echo "${GREEN}✓ HTML coverage report generated: coverage.html${RESET}"
	@echo "${BLUE}Open coverage.html in your browser to view the report${RESET}"

gocovmerge: ## Install gocovmerge if not present
	@which gocovmerge >/dev/null 2>&1 || (echo "${YELLOW}gocovmerge not found. Installing...${RESET}" && go install github.com/wadey/gocovmerge@latest)

gotestsum: ## Install gotestsum if not present
	@which gotestsum >/dev/null 2>&1 || (echo "${YELLOW}gotestsum not found. Installing...${RESET}" && go install gotest.tools/gotestsum@latest)

bench: ## Run rules engine benchmark
	@echo "${YELLOW}Running rules engine benchmark...${RESET}"
	@go test -bench=. -benchmem -benchtime=5s -run=^$$ ./src/workers/internal/rules/...
	@echo "${GREEN}✓ Benchmark completed!${RESET}"

clean: ## Remove build artifacts and Docker volumes
	@echo "${YELLOW}Cleaning artifacts...${RESET}"
	@rm -rf bin/ coverage.out coverage.html coverage-unit.out coverage-integration.out coverage-combined.out
	@rm -rf src/ui/node_modules src/ui/.next
	@docker-compose down -v
	@echo "${YELLOW}Removing stale Docker containers...${RESET}"
	@docker container prune -f
	@go clean -testcache -cache
	@echo "${GREEN}✓ Cleanup completed!${RESET}"

ui: ## Start UI service only
	@echo "${YELLOW}Building and starting UI service...${RESET}"
	@DOCKER_BUILDKIT=1 docker-compose build ui
	@docker-compose up -d ui
	@echo "${GREEN}✓ UI service started!${RESET}"
	@echo "${BLUE}UI:${RESET}  http://localhost:3000"

api: ## Start API service with infrastructure (postgres + redis)
	@echo "${YELLOW}Starting infrastructure services (postgres + redis)...${RESET}"
	@docker-compose up -d postgres redis
	@echo "${YELLOW}Building API service...${RESET}"
	@DOCKER_BUILDKIT=1 docker-compose build api
	@echo "${YELLOW}Waiting for infrastructure to be healthy...${RESET}"
	@docker-compose up -d api
	@echo "${GREEN}✓ API service with infrastructure started!${RESET}"
	@echo "${BLUE}API:${RESET} http://localhost:8080"
	@make logs

worker: ## Start Worker service with infrastructure (postgres + redis)
	@echo "${YELLOW}Starting infrastructure services (postgres + redis)...${RESET}"
	@docker-compose up -d postgres redis
	@echo "${YELLOW}Building Worker service...${RESET}"
	@DOCKER_BUILDKIT=1 docker-compose build worker
	@echo "${YELLOW}Waiting for infrastructure to be healthy...${RESET}"
	@docker-compose up -d worker
	@echo "${GREEN}✓ Worker service with infrastructure started!${RESET}"
	@make logs

dev-api: air infra-up ## Start API with hot reload (requires air)
	@echo "${YELLOW}Starting API with hot reload...${RESET}"
	@echo "${YELLOW}Waiting for infrastructure to be ready...${RESET}"
	@sleep 3
	@air -c .air.api.toml
	@echo "${GREEN}✓ API with hot reload started!${RESET}"

dev-worker: air infra-up ## Start Worker with hot reload (requires air)
	@echo "${YELLOW}Starting Worker with hot reload...${RESET}"
	@echo "${YELLOW}Waiting for infrastructure to be ready...${RESET}"
	@sleep 3
	@air -c .air.worker.toml
	@echo "${GREEN}✓ Worker with hot reload started!${RESET}"

air: ## Install air if not present
	@which air >/dev/null 2>&1 || (echo "${YELLOW}air not found. Installing...${RESET}" && go install github.com/cosmtrek/air@latest)

infra-up: ## Start only infrastructure (postgres + redis)
	@echo "${YELLOW}Starting infrastructure services...${RESET}"
	@docker-compose up -d postgres redis
	@echo "${GREEN}✓ Infrastructure started!${RESET}"

infra-down: ## Stop only infrastructure
	@echo "${YELLOW}Stopping infrastructure services...${RESET}"
	@docker-compose stop postgres redis
	@echo "${GREEN}✓ Infrastructure stopped!${RESET}"

test-watch: gotestsum ## Watch mode for tests (re-runs on file changes)
	@echo "${YELLOW}Starting test watch mode...${RESET}"
	@gotestsum --watch --format testdox -- -short ./...

coverage-ci: gotestsum ## Run coverage check matching CI requirements (80% minimum)
	@echo "${YELLOW}Running coverage check (CI mode)...${RESET}"
	@rm -f coverage-ci.out
	@gotestsum --format testdox -- -race -parallel 4 -coverprofile=coverage-ci.out -covermode=atomic ./src/...
	@echo "${YELLOW}Coverage report:${RESET}"
	@go tool cover -func=coverage-ci.out | tail -1
	@COVERAGE=$$(go tool cover -func=coverage-ci.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	MIN_COVERAGE=80.0; \
	echo "Coverage: $${COVERAGE}%"; \
	echo "Minimum required: $${MIN_COVERAGE}%"; \
	if [ $$(echo "$${COVERAGE} < $${MIN_COVERAGE}" | bc -l) -eq 1 ]; then \
		echo "${RED}✗ Coverage $${COVERAGE}% is below minimum threshold of $${MIN_COVERAGE}%${RESET}"; \
		exit 1; \
	fi; \
	echo "${GREEN}✓ Coverage meets minimum threshold${RESET}"

check-deps: ## Run dependency security scans
	@echo "${YELLOW}Running security scans...${RESET}"
	@echo "${YELLOW}→ govulncheck (Go vulnerabilities)...${RESET}"
	@which govulncheck >/dev/null 2>&1 || (echo "${YELLOW}govulncheck not found. Installing...${RESET}" && go install golang.org/x/vuln/cmd/govulncheck@latest)
	@govulncheck ./...
	@echo "${YELLOW}→ npm audit (UI dependencies)...${RESET}"
	@cd src/ui && npm audit
	@echo "${GREEN}✓ Security scans completed!${RESET}"
