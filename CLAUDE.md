# CLAUDE.md

AlgoShield — high-performance fraud detection & AML transaction analysis. Ultra-low latency (<50ms), expr-lang rules engine with hot-reload.

## Tech Stack

**Backend (Go 1.25.5):** Fiber v2, pgx v5, go-redis v9, Asynq, golang-jwt v5, expr-lang v1.17.7, OpenTelemetry
**Frontend (Vue.js 3.5):** Vue 3.5.26 (`<script setup>`), TypeScript 5.9.3, Pinia 3.0, Vuetify 3.11.6, Tailwind CSS 4.1.18, Vite 7.3.1, vue-i18n 11.2.8, Font Awesome 7.1.0
**Infrastructure:** PostgreSQL 16, Redis 7, Docker + Docker Compose

## Commands   

```bash
make install                 # Install all deps (Go + npm + golangci-lint)
make up                      # Start all services (Docker Compose)
make up-dev                  # Start with faster health checks
make down / make logs        # Stop services / tail logs

# Local dev (requires postgres + redis running)
cd src/api/cmd && go run main.go       # API on :8080
cd src/workers/cmd && go run main.go   # Worker
cd src/ui && npm run dev               # UI on :5173

# Testing & quality
make test                    # All (API + UI)
make test-api                # Go with race detector
make test-ui                 # Vitest
make lint                    # golangci-lint
cd src/ui && npm run type-check
make test-coverage           # Combined unit + integration (80% minimum)
make check-deps              # govulncheck + npm audit
```

Migrations: `scripts/migrations/` — run in order: `001_schema.sql` → `002_indexes.sql` → `003_test_data.sql`

## Architecture

Vertical slice — organized by feature, not by technical layer.

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

- **Message Queue:** API → Redis (Asynq) → Worker → rules evaluation → DB. 3 priority levels, batch + semaphore concurrency.
- **Rules Engine:** expr-lang, hot-reload 10s, cached per schema. Helpers: `velocityCount`, `velocitySum`, `pointInPolygon`.
- **Auth:** JWT access (24h) + refresh (168h), Redis revocation, CSRF for state-changing ops, RBAC (admin | rule_editor | viewer).
- **Transactions:** `pending` → `approved` / `rejected` / `in_review`. Synthetic mode via `X-Synthetic-Mode` header, separate tables.
- **Caching:** Rules (Redis 5min), Schema (memory + pub/sub), Dashboard (Redis 30s), CSRF (Redis 24h), Token revocation (Redis sets).

## Cross-Cutting Rules

- **i18n:** No hardcoded user-facing strings. Use `$t('key')` / `useI18n()`. Add to BOTH `pt-BR.json` and `en-US.json`.
- **Env vars:** All config via env vars. Update `.env.example` when adding new vars. All timeouts MUST be configurable.
- **Database:** ALWAYS parameterized queries (`$1, $2`) — NEVER concatenate user input.
- **Commits:** Conventional: `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`. Git hooks: `./scripts/install-hooks.sh`

## Security Essentials

- **SQL injection:** Parameterized queries only — use `$1, $2` placeholders instead of string concatenation
- **XSS:** Vue auto-escapes; use `{{ }}` instead of `v-html` with user input; use CSP instead of `eval()`
- **CSRF:** Redis-backed tokens (24h TTL), auto-added by `api.post/put/patch/delete` — use API client instead of raw `fetch()`
- **Auth:** bcrypt (cost 12+), JWT with Redis revocation — use env vars instead of hardcoded secrets
- **Tokens:** Pinia (memory) only — use `authStore` instead of `localStorage` for access/CSRF tokens
- **Headers:** `X-Content-Type-Options`, `X-Frame-Options`, `Strict-Transport-Security`, `CSP`
- **CORS:** Use `CORS_ALLOWED_ORIGINS` instead of `*` in production
- **Sensitive data:** Never log passwords, tokens, PII — mask in error messages instead
- **Scanning:** `govulncheck`, `npm audit`, `gitleaks`, `semgrep`

## Documentation Map

| Topic | Location | When to read |
|-------|----------|-------------|
| Go API conventions | `src/api/CLAUDE.md` | Working on backend code |
| Vue/TS conventions | `src/ui/CLAUDE.md` | Working on frontend code |
| Testing conventions | `docs/agents/CLAUDE.md` | Writing or reviewing tests |
| Full project spec & domain context | `openspec/project.md` | New features, architecture decisions |
| OpenSpec workflow | `openspec/AGENTS.md` | Proposals, breaking changes |
| Unit testing (full guide) | `docs/agents/unit-test.md` | Writing Go/Vue unit tests |
| Integration testing (full guide) | `docs/agents/integration-test.md` | Writing integration tests |
| i18n guidelines | `src/ui/src/locales/README.md` | Adding translations |
| API documentation | `README.md` § "API Documentation" | API endpoints reference |
| Environment variables | `.env.example` | Configuration reference |
| DI wiring | `src/api/internal/routes/routes.go` | Route registration |
