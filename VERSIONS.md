# Versões das Tecnologias - AlgoShield

Este documento lista as versões de todas as tecnologias utilizadas no projeto AlgoShield.

**Última atualização**: Dezembro 2025

## 🔧 Tecnologias Core

### Go
- **Versão**: `1.23.4`
- **Motivo**: Versão estável mais recente com suporte a `GOEXPERIMENT=rangefunc`
- **Onde**: `go.mod`, `Dockerfile.api`, `Dockerfile.worker`, `.github/workflows/ci.yml`

### PostgreSQL
- **Versão**: `17-alpine`
- **Motivo**: Versão estável mais recente com melhorias de performance
- **Onde**: `docker-compose.yml`, `.github/workflows/ci.yml`
- **Mudanças da v16 para v17**:
  - Melhorias significativas de performance
  - Novo sistema de vacuum mais eficiente
  - Melhor suporte para JSON/JSONB
  - Otimizações em queries complexas

### Redis
- **Versão**: `7.4-alpine`
- **Motivo**: Versão estável mais recente com melhorias de performance
- **Onde**: `docker-compose.yml`, `.github/workflows/ci.yml`
- **Mudanças da v7.0 para v7.4**:
  - Melhorias de performance em operações de lista
  - Otimizações de memória
  - Melhor suporte para clustering

### Node.js
- **Versão**: `22-alpine`
- **Motivo**: Versão LTS mais recente
- **Onde**: `Dockerfile.ui`
- **Mudanças da v20 para v22**:
  - Performance melhorada do V8
  - Suporte nativo a ESM melhorado
  - Melhorias no gerenciamento de memória

## 📦 Dependências Go

### Principais

| Pacote | Versão | Descrição |
|--------|--------|-----------|
| `github.com/gofiber/fiber/v2` | v2.52.0 | Framework web de alta performance |
| `github.com/jackc/pgx/v5` | v5.5.1 | Driver PostgreSQL otimizado |
| `github.com/redis/go-redis/v9` | v9.4.0 | Cliente Redis oficial |
| `github.com/google/uuid` | v1.5.0 | Geração de UUIDs |
| `github.com/joho/godotenv` | v1.5.1 | Gerenciamento de variáveis de ambiente |

### Dependências Indiretas

| Pacote | Versão |
|--------|--------|
| `github.com/valyala/fasthttp` | v1.51.0 |
| `golang.org/x/crypto` | v0.17.0 |
| `golang.org/x/sys` | v0.15.0 |
| `golang.org/x/text` | v0.14.0 |

## 🎨 Dependências Frontend (UI)

### SvelteKit
- **Versão**: `^2.0.0`
- **Framework**: Svelte 5

### Principais

| Pacote | Versão | Descrição |
|--------|--------|-----------|
| `@sveltejs/kit` | ^2.0.0 | Framework full-stack |
| `@sveltejs/vite-plugin-svelte` | ^4.0.0 | Plugin Vite para Svelte |
| `svelte` | ^5.0.0 | Framework reativo |
| `typescript` | ^5.0.0 | Type safety |
| `vite` | ^5.0.0 | Build tool |

## 🐳 Docker

### Imagens Base

| Serviço | Imagem | Tag |
|---------|--------|-----|
| API | `golang` | `1.23.4-alpine` |
| Worker | `golang` | `1.23.4-alpine` |
| UI | `node` | `22-alpine` |
| PostgreSQL | `postgres` | `17-alpine` |
| Redis | `redis` | `7.4-alpine` |
| Runtime (API/Worker) | `alpine` | `latest` |

### Docker Compose
- **Versão**: `3.8`

## 🔄 CI/CD

### GitHub Actions

| Action | Versão |
|--------|--------|
| `actions/checkout` | v4 |
| `actions/setup-go` | v5 |
| `actions/cache` | v3 |
| `codecov/codecov-action` | v3 |
| `golangci/golangci-lint-action` | v3 |

## 📊 Histórico de Atualizações

### Dezembro 2025
- ✅ Go: `1.23` → `1.23.4`
- ✅ PostgreSQL: `16-alpine` → `17-alpine`
- ✅ Redis: `7-alpine` → `7.4-alpine`
- ✅ Node.js: `20-alpine` → `22-alpine`

## 🔍 Como Verificar Versões

### Go
```bash
go version
# Esperado: go version go1.23.4 darwin/arm64 (ou seu OS)
```

### PostgreSQL (Docker)
```bash
docker-compose exec postgres psql --version
# Esperado: psql (PostgreSQL) 17.x
```

### Redis (Docker)
```bash
docker-compose exec redis redis-server --version
# Esperado: Redis server v=7.4.x
```

### Node.js (Docker)
```bash
docker-compose exec ui node --version
# Esperado: v22.x.x
```

## 🔄 Política de Atualizações

### Atualizações Automáticas
- ❌ **Não recomendado** para produção
- ✅ Revisar release notes antes de atualizar

### Frequência de Revisão
- **Go**: Atualizar para patches de segurança imediatamente
- **PostgreSQL**: Revisar a cada minor release
- **Redis**: Revisar a cada minor release
- **Node.js**: Manter na versão LTS mais recente
- **Dependências Go**: Revisar mensalmente

### Processo de Atualização

1. **Verificar Release Notes**
   - Ler changelog completo
   - Identificar breaking changes
   - Verificar deprecations

2. **Testar em Desenvolvimento**
   ```bash
   # Atualizar versões
   # Rodar testes
   make test
   
   # Testar localmente
   make dev-infra
   make run-api
   make run-worker
   make dev-ui
   ```

3. **Validar CI/CD**
   - Verificar se todos os testes passam
   - Verificar build das imagens Docker

4. **Deploy Gradual**
   - Staging primeiro
   - Monitorar métricas
   - Produção após validação

## 🔐 Versões de Segurança

### Alertas de Segurança
- Configurar GitHub Dependabot
- Revisar CVEs regularmente
- Aplicar patches de segurança imediatamente

### Comandos Úteis

```bash
# Verificar vulnerabilidades Go
go list -json -m all | nancy sleuth

# Atualizar dependências Go
go get -u ./...
go mod tidy

# Verificar vulnerabilidades npm
cd ui && npm audit

# Atualizar dependências npm
cd ui && npm update
```

## 📚 Referências

- [Go Releases](https://go.dev/dl/)
- [PostgreSQL Releases](https://www.postgresql.org/support/versioning/)
- [Redis Releases](https://redis.io/download)
- [Node.js Releases](https://nodejs.org/en/about/releases/)
- [Docker Hub - Golang](https://hub.docker.com/_/golang)
- [Docker Hub - PostgreSQL](https://hub.docker.com/_/postgres)
- [Docker Hub - Redis](https://hub.docker.com/_/redis)
- [Docker Hub - Node](https://hub.docker.com/_/node)

---

**AlgoShield** - Sempre atualizado, sempre seguro! 🛡️

