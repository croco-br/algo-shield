# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AlgoShield is an open-source, high-performance fraud detection and AML transaction analysis system. It processes transactions with ultra-low latency (<50ms) using a custom expression-based rules engine powered by expr-lang.

**Key Capabilities:** Real-time transaction processing, custom rules engine with hot-reload, event schema management, risk scoring, JWT-based auth with RBAC, white-label branding, synthetic event generation.

## Tech Stack

**Backend (Go 1.25.5):** Fiber v2, pgx v5, go-redis v9, Asynq, golang-jwt v5, expr-lang v1.17.7, OpenTelemetry
**Frontend (Vue.js 3.5):** Vue 3.5.26 (`<script setup>`), TypeScript 5.9.3, Pinia 3.0, Vuetify 3.11.6, Tailwind CSS 4.1.18, Vite 7.3.1, vue-i18n 11.2.8, Font Awesome 7.1.0
**Infrastructure:** PostgreSQL 16, Redis 7, Docker + Docker Compose

## Common Commands

```bash
make install                 # Install all deps (Go + npm + golangci-lint)
./scripts/install-hooks.sh   # Git hooks (pre-commit: lint, pre-push: tests)
make up                      # Start all services (Docker Compose)
make up-dev                  # Start with faster health checks
make down                    # Stop all services
make logs                    # Tail service logs

# Local dev (requires postgres + redis running)
cd src/api/cmd && go run main.go       # API on :8080
cd src/workers/cmd && go run main.go   # Worker
cd src/ui && npm run dev               # UI on :5173

# Testing
make test                    # All (API + UI)
make test-api                # Go with race detector
make test-ui                 # Vitest
go test -race -run TestName ./src/api/internal/auth/...  # Specific
go test -count=50 ./src/...  # Flaky detection
make bench                   # Rules engine benchmark

# Quality
make lint                    # golangci-lint
cd src/ui && npm run type-check

# Coverage
make test-coverage           # Combined unit + integration
make coverage-ci             # 80% minimum check

# Security
make check-deps              # govulncheck + npm audit
```

### Database Migrations

SQL files in `scripts/migrations/` — run in order: `001_schema.sql` → `002_indexes.sql` → `003_test_data.sql`.

## Architecture

Vertical slice architecture — code organized by feature, not by technical layer.

```
src/
├── api/              # RESTful API (Fiber v2)
│   ├── cmd/          # Entry point
│   └── internal/     # Feature modules (auth, transactions, rules, schemas, etc.)
├── workers/          # Background processing (Asynq)
│   ├── cmd/          # Entry point
│   └── internal/     # Rule engine, transaction processing
├── pkg/              # Shared packages (config, database, errors, models, queue, etc.)
└── ui/               # Vue.js frontend
    └── src/          # components, composables, lib, locales, stores, views
```

**Dependency flow:** Handler → Service (interfaces) → Repository → DB
**Wiring:** `src/api/internal/routes/routes.go`

### Key Subsystems

**Message Queue (Asynq):** API enqueues → Redis → Worker pops (BRPOP) → Evaluate rules → Save result. 3 priority levels, batch processing with semaphore concurrency.

**Rules Engine:** expr-lang expressions, hot-reload every 10s, cached per schema. Helpers: `velocityCount`, `velocitySum`, `pointInPolygon`. Field paths validated (`^[a-zA-Z0-9_.]+$`).

**Auth:** JWT access (24h) + refresh (168h), Redis revocation, CSRF for state-changing ops, RBAC (admin | rule_editor | viewer).

**Transactions:** `pending` → `approved` / `rejected` / `in_review`. Synthetic mode uses separate tables via `X-Synthetic-Mode` header.

**Caching:** Rules (Redis, 5min TTL), Schema (in-memory + Redis pub/sub), Dashboard (Redis, 30s TTL), CSRF (Redis, 24h), Token revocation (Redis sets).

## Key Conventions

### Go

- Standard conventions (gofmt, golangci-lint), `internal/` for private code
- Error wrapping: `fmt.Errorf("context: %w", err)`
- **All handlers MUST validate requests**
- `context.Context` for all DB/Redis ops
- Interfaces in the consumer package; mocks use `Mock*` prefix in `_test.go`

### TypeScript / Vue

- `<script setup lang="ts">`, strict mode, PascalCase components
- Pinia for state, Tailwind for styling, base components (BaseButton, BaseInput, etc.)
- Type-safe API calls via `src/ui/src/lib/api.ts`
- **CSRF**: Always use `api.post/put/patch/delete` — never raw `fetch()`

### i18n

**CRITICAL: No hardcoded user-facing strings.** Use `$t('key')` in templates, `useI18n()` in scripts.
Add to BOTH `pt-BR.json` and `en-US.json`. Key convention: `common.*`, `auth.*`, `views.<name>.*`, `errors.*`.
Details: `src/ui/src/locales/README.md`

### Database

- **ALWAYS parameterized queries** (`$1, $2`) — NEVER concatenate user input
- `pgxpool` for pooling, `pgx.Tx` for atomic ops

### Environment Variables

- All config via env vars (see `.env.example`)
- **MANDATORY**: Update `.env.example` when adding new vars
- **ALL timeouts MUST be configurable** — no hardcoded timeout values

## Security

**OWASP Top 10 compliance required.** Key rules:

- **SQL injection**: Parameterized queries only
- **XSS**: Vue auto-escapes; never `v-html` with user input; no `eval()`
- **CSRF**: Tokens in Redis (24h TTL), auto-added by API client
- **Auth**: bcrypt (cost 12+), JWT with Redis revocation, rate-limit login
- **Tokens**: Pinia (memory) only — NEVER localStorage for access/CSRF tokens
- **Headers**: `X-Content-Type-Options`, `X-Frame-Options`, `Strict-Transport-Security`, `CSP`
- **CORS**: Never `*` in production — use `CORS_ALLOWED_ORIGINS`
- **Sensitive data**: Never log passwords, tokens, PII
- **Scanning**: `govulncheck`, `npm audit`, `gitleaks`, `semgrep`

Full security guidelines: `openspec/project.md` § "Code Security Guidelines"

## Testing

**Coverage minimum: 80%. Race detection: mandatory.**

- **Unit tests**: Mock dependencies, table-driven, AAA pattern, `-race` flag
- **Integration tests**: `//go:build integration`, testcontainers-go, <500ms per test
- **UI tests**: Vitest + Vue Test Utils, `flushPromises()` for async, <300ms per test
- **Flaky detection**: `go test -count=50` / `npm run test -- --repeat=20`

**Anti-patterns** — do NOT generate:
- Trivial constructor tests (no logic = no test)
- Interface implementation tests (compiler checks this)
- Duplicate validation across layers (test once at service level)
- `exists()` without verifying actual values
- 3+ similar tests → consolidate to table-driven

Full testing guides: `docs/agents/unit-test.md`, `docs/agents/integration-test.md`

## Regression Prevention

Follow when adding new features. Full checklists: `openspec/project.md` § "Regression Prevention Rules"

1. **Before**: Impact analysis, identify affected modules, plan tests
2. **During**: Write tests alongside code, run frequently, check for race conditions
3. **After**: Full suite passes, coverage ≥80%, linters clean, docs updated
4. **Performance**: <50ms per transaction, no N+1 queries, connection pool limits respected

## Development Workflow

1. Install Git hooks: `./scripts/install-hooks.sh`
2. Pre-commit runs formatting + linting; pre-push runs full test suite
3. Conventional commits: `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, etc.
4. For significant changes, use OpenSpec proposals (`openspec/AGENTS.md`)

## Known Issues

- **Layout bugs**: Some remain from Vue.js migration — needs investigation
- **Database migrations**: Script-based process needs a proper migration library
- **OpenTelemetry**: Infrastructure defined but not fully implemented

## Canonical Documentation

| Topic | Location |
|-------|----------|
| Project conventions & security | `openspec/project.md` |
| OpenSpec workflow | `openspec/AGENTS.md` |
| Unit testing guide | `docs/agents/unit-test.md` |
| Integration testing guide | `docs/agents/integration-test.md` |
| i18n guidelines | `src/ui/src/locales/README.md` |
| API documentation | `README.md` § "API Documentation" |
| Environment variables | `.env.example` |
| DI wiring | `src/api/internal/routes/routes.go` |