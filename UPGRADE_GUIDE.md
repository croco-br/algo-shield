# Guia de Atualização - AlgoShield

Este guia ajuda você a atualizar o AlgoShield para as versões mais recentes das tecnologias.

## 📋 Resumo das Atualizações (Dezembro 2025)

### Versões Atualizadas

| Tecnologia | Versão Anterior | Nova Versão | Tipo |
|------------|----------------|-------------|------|
| Go | 1.23 | 1.23.4 | Patch |
| PostgreSQL | 16-alpine | 17-alpine | Major |
| Redis | 7-alpine | 7.4-alpine | Minor |
| Node.js | 20-alpine | 22-alpine | Major |

### Novas Funcionalidades

✅ **Git Hooks Automatizados**
- Pre-commit: Testes unitários + verificação de formatação
- Commit-msg: Validação de Conventional Commits
- Pre-push: Suite completa de testes

## 🚀 Como Atualizar

### Opção 1: Usando Docker (Recomendado)

Se você está usando Docker, a atualização é automática:

```bash
# 1. Parar serviços
docker-compose down

# 2. Atualizar código
git pull origin main

# 3. Rebuild imagens
docker-compose build --no-cache

# 4. Iniciar serviços
docker-compose up -d

# 5. Verificar saúde
curl http://localhost:8080/health
```

### Opção 2: Desenvolvimento Local

Se você está rodando localmente, siga estes passos:

#### 1. Atualizar Go para 1.23.4

**macOS (Homebrew)**:
```bash
brew update
brew upgrade go
go version  # Verificar: go1.23.4
```

**Linux**:
```bash
# Download
wget https://go.dev/dl/go1.23.4.linux-amd64.tar.gz

# Remover versão antiga
sudo rm -rf /usr/local/go

# Instalar nova versão
sudo tar -C /usr/local -xzf go1.23.4.linux-amd64.tar.gz

# Verificar
go version
```

**Windows**:
- Baixar instalador em: https://go.dev/dl/go1.23.4.windows-amd64.msi
- Executar instalador
- Verificar: `go version`

#### 2. Atualizar PostgreSQL para 17

**Via Docker (Recomendado)**:
```bash
# Backup dos dados
docker-compose exec postgres pg_dump -U algoshield algoshield > backup.sql

# Parar e remover container antigo
docker-compose down postgres

# Atualizar docker-compose.yml (já atualizado)
# Iniciar novo container
docker-compose up -d postgres

# Restaurar dados (se necessário)
docker-compose exec -T postgres psql -U algoshield algoshield < backup.sql
```

**Instalação Local** (não recomendado para dev):
- Consulte: https://www.postgresql.org/download/

#### 3. Atualizar Redis para 7.4

**Via Docker (Recomendado)**:
```bash
# Redis não tem dados persistentes críticos neste projeto
docker-compose down redis
docker-compose up -d redis
```

#### 4. Atualizar Node.js para 22

**macOS (Homebrew)**:
```bash
brew update
brew upgrade node
node --version  # Verificar: v22.x.x
```

**Linux (nvm)**:
```bash
nvm install 22
nvm use 22
nvm alias default 22
node --version
```

**Windows**:
- Baixar instalador em: https://nodejs.org/
- Executar instalador
- Verificar: `node --version`

#### 5. Atualizar Dependências do Projeto

```bash
# Go
go get -u ./...
go mod tidy

# UI
cd ui
npm update
cd ..
```

## 🔧 Instalar Git Hooks

**Novo recurso!** Instale os hooks para testes automatizados:

```bash
# Instalar hooks
./scripts/install-hooks.sh

# Ou via Makefile
make install-hooks
```

Isso irá:
- ✅ Rodar testes antes de cada commit
- ✅ Validar mensagens de commit
- ✅ Rodar suite completa antes de push

## ✅ Verificação Pós-Atualização

Execute estes comandos para verificar se tudo está funcionando:

### 1. Verificar Versões

```bash
# Go
go version
# Esperado: go version go1.23.4

# PostgreSQL (Docker)
docker-compose exec postgres psql --version
# Esperado: psql (PostgreSQL) 17.x

# Redis (Docker)
docker-compose exec redis redis-server --version
# Esperado: Redis server v=7.4.x

# Node (Docker)
docker-compose exec ui node --version
# Esperado: v22.x.x
```

### 2. Rodar Testes

```bash
# Testes Go
make test

# Ou manualmente
GOEXPERIMENT=rangefunc go test -v ./...
```

### 3. Testar Aplicação

```bash
# Iniciar serviços
make docker-up

# Verificar health
curl http://localhost:8080/health

# Testar transação
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "external_id": "test_upgrade",
    "amount": 1000,
    "currency": "BRL",
    "from_account": "ACC001",
    "to_account": "ACC002",
    "type": "transfer",
    "timestamp": "2024-12-05T10:00:00Z"
  }'

# Verificar UI
open http://localhost:5173
```

## 🐛 Troubleshooting

### Erro: "go: cannot find main module"

**Solução**:
```bash
cd /caminho/para/algo-shield
go mod download
```

### Erro: PostgreSQL não inicia

**Solução**:
```bash
# Verificar logs
docker-compose logs postgres

# Remover volumes e reiniciar
docker-compose down -v
docker-compose up -d postgres
```

### Erro: "permission denied" nos hooks

**Solução**:
```bash
chmod +x .githooks/*
chmod +x scripts/install-hooks.sh
./scripts/install-hooks.sh
```

### Erro: Testes falhando após atualização

**Solução**:
```bash
# Limpar cache
go clean -testcache
go clean -modcache
go mod download

# Rodar novamente
make test
```

### Erro: UI não compila

**Solução**:
```bash
cd ui
rm -rf node_modules package-lock.json
npm install
npm run dev
```

## 📊 Breaking Changes

### PostgreSQL 16 → 17

**Mudanças que podem afetar você**:
- ✅ Sintaxe SQL permanece compatível
- ✅ Queries existentes continuam funcionando
- ⚠️ Performance pode melhorar (requer reindex)

**Ação recomendada**:
```sql
-- Conectar ao banco
docker-compose exec postgres psql -U algoshield algoshield

-- Reindexar para melhor performance
REINDEX DATABASE algoshield;
```

### Node.js 20 → 22

**Mudanças que podem afetar você**:
- ✅ Código SvelteKit continua compatível
- ✅ Dependências npm continuam funcionando
- ⚠️ Algumas dependências podem ter warnings

**Ação recomendada**:
```bash
cd ui
npm audit fix
```

## 🔄 Rollback

Se algo der errado, você pode voltar para as versões anteriores:

### Via Git

```bash
# Voltar para commit anterior
git checkout <commit-hash-anterior>

# Rebuild
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

### Versões Específicas

Edite `docker-compose.yml`:

```yaml
# Versões antigas
postgres:
  image: postgres:16-alpine

redis:
  image: redis:7-alpine
```

Edite `go.mod`:

```go
go 1.23
```

Depois:
```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

## 📚 Recursos Adicionais

- [VERSIONS.md](./VERSIONS.md) - Lista completa de versões
- [CHANGELOG.md](./CHANGELOG.md) - Histórico de mudanças
- [docs/GIT_HOOKS.md](./docs/GIT_HOOKS.md) - Documentação dos hooks

## 💬 Suporte

Se encontrar problemas:

1. Verifique os logs: `docker-compose logs`
2. Consulte [Issues](https://github.com/yourusername/algo-shield/issues)
3. Abra uma nova issue se necessário

---

**AlgoShield** - Sempre atualizado! 🛡️

