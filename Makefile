.PHONY: help install up down logs test bench clean

# Colors
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
BLUE   := $(shell tput -Txterm setaf 4)
RESET  := $(shell tput -Txterm sgr0)

help: ## Mostra esta mensagem de ajuda
	@echo '${BLUE}Comandos disponíveis:${RESET}'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-15s${RESET} %s\n", $$1, $$2}'

install: ## Instala todas as dependências (Go + npm)
	@echo "${YELLOW}Instalando dependências Go...${RESET}"
	@go mod download
	@echo "${YELLOW}Instalando dependências UI (npm)...${RESET}"
	@cd src/ui && npm install
	@echo "${GREEN}✓ Todas as dependências instaladas!${RESET}"

up: ## Sobe todos os serviços no Docker (API + Worker + UI + infra)
	@echo "${YELLOW}Iniciando todos os serviços...${RESET}"
	@docker-compose up -d --build
	@echo "${GREEN}✓ Serviços iniciados!${RESET}"
	@echo "${BLUE}API:${RESET} http://localhost:8080"
	@echo "${BLUE}UI:${RESET}  http://localhost:5173"

down: ## Para todos os serviços
	@echo "${YELLOW}Parando todos os serviços...${RESET}"
	@docker-compose down
	@echo "${GREEN}✓ Serviços parados!${RESET}"

logs: ## Visualiza logs dos serviços
	@docker-compose logs -f

test: ## Executa todos os testes em paralelo com race detector
	@echo "${YELLOW}Executando todos os testes com race detector...${RESET}"
	@GOEXPERIMENT=greenteagc,rangefunc go test -race -v -parallel 4 ./src/...
	@echo "${GREEN}✓ Testes concluídos!${RESET}"

bench: ## Executa benchmark do motor de regras
	@echo "${YELLOW}Executando benchmark do motor de regras...${RESET}"
	@GOEXPERIMENT=greenteagc,rangefunc go test -bench=. -benchmem -benchtime=5s -run=^$$ ./src/workers/internal/rules/...
	@echo "${GREEN}✓ Benchmark concluído!${RESET}"

clean: ## Remove artefatos de build e volumes Docker
	@echo "${YELLOW}Limpando artefatos...${RESET}"
	@rm -rf bin/ coverage.out coverage.html
	@rm -rf src/ui/node_modules src/ui/.svelte-kit
	@docker-compose down -v
	@go clean -testcache -cache
	@echo "${GREEN}✓ Limpeza concluída!${RESET}"
