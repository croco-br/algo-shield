# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AlgoShield is an open-source, high-performance fraud detection and AML transaction analysis system. It processes transactions with ultra-low latency (<50ms) using a custom expression-based rules engine powered by expr-lang.

**Key Capabilities:**
- Real-time transaction processing (<50ms latency)
- Custom rules engine with hot-reload support
- Event schema management with automatic field extraction
- Risk scoring and transaction classification
- JWT-based authentication with RBAC
- White-label branding customization
- Synthetic event generation

## Tech Stack

**Backend (Go 1.25.5):**
- Fiber v2 (web framework)
- pgx v5 (PostgreSQL driver with connection pooling)
- go-redis v9 (Redis client)
- Asynq (task queue for async processing)
- golang-jwt v5 + bcrypt (authentication)
- expr-lang v1.17.7 (rule expression engine)
- OpenTelemetry (observability infrastructure)

**Frontend (Vue.js 3.5):**
- Vue 3.5.26 with Composition API (`<script setup>`)
- TypeScript 5.9.3 (strict mode)
- Pinia 3.0 (state management)
- Vuetify 3.11.6 (Material Design components)
- Tailwind CSS 4.1.18 (utility-first styling)
- Vite 7.3.1 (build tool)
- vue-i18n 11.2.8 (internationalization - Composition API mode)
- Font Awesome 7.1.0 (icons)
- Prism.js 1.30 (syntax highlighting)

**Infrastructure:**
- PostgreSQL 16 (primary data store)
- Redis 7 (message queue + caching)
- Docker + Docker Compose (containerization)

## Common Commands

### Development Setup

```bash
# Install all dependencies (Go + npm + golangci-lint)
make install

# Install Git hooks (recommended - includes pre-commit checks)
./scripts/install-hooks.sh

# Start all services with Docker Compose
make up

# Start only infrastructure (postgres + redis)
docker-compose up -d postgres redis
```

### Local Development

```bash
# Run API locally (requires postgres + redis running)
cd src/api/cmd && go run main.go

# Run Worker locally (requires postgres + redis running)
cd src/workers/cmd && go run main.go

# Run UI locally (Vite dev server on port 5173)
cd src/ui && npm run dev
```

### Testing

```bash
# Run all tests
make test                    # API + UI
make test-api                # API only with race detector

# Run with better output
gotestsum --format testdox -- -race -parallel 8 ./src/...

# Run specific tests
go test -race ./src/api/internal/auth/service_test.go
go test -race -run TestServiceName ./src/api/internal/auth/...

# Integration tests (requires Docker)
gotestsum --format testdox -- -tags=integration -race -parallel 2 ./src/api/... ./src/workers/...

# Check flaky tests
go test -count=50 ./src/...

# UI tests
make test-ui                 # Run tests
cd src/ui && npm run test:coverage  # With coverage
cd src/ui && npm run test:ui        # Interactive

# Benchmarks
make bench
go test -bench=. -benchmem -benchtime=5s -run=^$$ ./src/workers/internal/rules/...
```

### Building and Linting

```bash
# Build all Docker images in parallel (with BuildKit)
make build

# Build only changed services (fast incremental)
make build-fast

# Run linters (golangci-lint)
make lint

# Type check UI
cd src/ui && npm run type-check
```

### Database Migrations

```bash
# Run migrations manually (local development)
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=algoshield
export POSTGRES_PASSWORD=algoshield_secret
export POSTGRES_DB=algoshield

# Run in order
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/001_schema.sql
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/002_indexes.sql
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/003_test_data.sql
```

**Note:** Current migration process is being improved - evaluate migration libraries for better workflow.

## Architecture

### High-Level Structure

AlgoShield follows a **vertical slice architecture** pattern, organizing code by feature rather than technical layers:

```
src/
├── api/              # RESTful API service (Fiber v2)
│   ├── cmd/          # Application entry point
│   └── internal/     # Feature modules (auth, transactions, rules, schemas, etc.)
│       ├── auth/           # Authentication & JWT
│       ├── branding/       # White-label customization
│       ├── dashboard/      # Metrics & analytics
│       ├── groups/         # User group management
│       ├── health/         # Health checks
│       ├── permissions/    # User permission management
│       ├── roles/          # Role management
│       ├── routes/         # Route setup & middleware
│       ├── rules/          # Rule management endpoints
│       ├── schemas/        # Event schema management
│       ├── shared/         # Shared middleware & utilities
│       ├── system/         # System configuration
│       ├── transactions/   # Transaction processing API
│       └── user/           # User management
├── workers/          # Background job processing (Asynq)
│   ├── cmd/          # Worker entry point
│   └── internal/     # Rule engine, transaction processing
│       ├── config/         # Worker configuration
│       ├── rules/          # Rules engine & evaluation
│       ├── schemas/        # Schema loading & caching
│       └── transactions/   # Transaction processing logic
├── pkg/              # Shared packages
│   ├── config/       # Configuration management
│   ├── csrf/         # CSRF token handling
│   ├── database/     # PostgreSQL & Redis clients
│   ├── errors/       # Custom error types
│   ├── models/       # Shared data models
│   ├── queue/        # Asynq queue client & worker
│   ├── rules/        # Rule repository
│   ├── tokenrevoke/  # Token revocation service
│   └── utils/        # Shared utilities
└── ui/               # Vue.js 3.5 frontend
    └── src/
        ├── components/   # Reusable Vue components
        ├── composables/  # Vue composition functions
        ├── lib/          # API client, error handling
        ├── locales/      # i18n translations (pt-BR, en-US)
        ├── plugins/      # Vue plugins (vuetify, i18n)
        ├── router/       # Vue Router configuration
        ├── stores/       # Pinia stores (state management)
        ├── types/        # TypeScript type definitions
        └── views/        # Page components
```

### Dependency Flow (Clean Architecture)

The project follows clean architecture principles with clear dependency boundaries:

1. **Presentation Layer** (`internal/*/handler.go`): HTTP handlers receive requests and return responses
2. **Business Layer** (`internal/*/service.go`): Services contain business logic and orchestrate operations
3. **Data Layer** (`internal/*/repository.go`): Repositories handle database operations
4. **Shared Layer** (`src/pkg/`): Shared utilities, models, and infrastructure

**Key Pattern:** Services depend on interfaces (not concrete implementations) for better testability and flexibility.

### Dependency Injection

All handlers/services use **DI** following SOLID:

```go
type Service struct { repo Repository }  // Interface, not concrete type
type Handler struct { service ServiceInterface }

// Wiring in routes.go
func Setup(app *fiber.App, db *pgxpool.Pool, redis *redis.Client, ...) {
    userRepo := user.NewPostgresUserRepository(db)  // Infrastructure
    userService := user.NewService(userRepo, ...)    // Business
    userHandler := user.NewHandler(userService)      // Presentation
}
```
See `src/api/internal/routes/routes.go` for complete wiring.

### Message Queue Architecture (Asynq)

**Flow:** API → Enqueue (`LPUSH`) → Redis Queue → Worker (`BRPOP`) → Process → Update DB

- **3 priority levels**: critical, default, low | **Configurable**: concurrency, retries, timeouts, batch size
- **Queue timeout**: 5s (returns `ErrTimeout` if no events - expected)
- **Batch processing**: Collects up to `batchSize`, processes in parallel (semaphore-controlled)

**Worker Flow:** Pop (`BRPOP`) → Collect batch → Parallel processing → Evaluate rules → Save result (status, matched rules, time) → Record metrics → Continue

### Rules Engine

**expr-lang** for flexible evaluation. Hot-reload every 10s, expression cache (cleared on schema changes).

**Flow:** Transaction → Schema Validation → Rule Evaluation → Risk Score → Action
**Actions:** `allow` (approved) | `block` (rejected, highest priority) | `review` (in_review)

**Helpers:**
- `velocityCount(fieldPath, timeWindow)` - Count transactions by field
- `velocitySum(fieldPath, [sumField,] timeWindow)` - Sum numeric field
- `pointInPolygon(lat, lon, polygon)` - Point in polygon (ray casting)

**Security:** Field paths validated (`^[a-zA-Z0-9_.]+$`) to prevent SQL injection.

### Caching Strategy

- **Rules**: Redis, 5min TTL | **Schema**: In-memory + Redis pub/sub | **Dashboard**: Redis, 30s TTL (synthetic/real)
- **Branding**: Redis with TTL | **CSRF**: Redis, 24h | **Token Revocation**: Redis sets (JWT blacklist)

### Authentication & Authorization

**JWT:** Access 24h, Refresh 168h (configurable), Redis revocation, CSRF for state-changing ops, rate limiting on auth endpoints

**RBAC:** admin (full access) | rule_editor (rules/schemas CRUD) | viewer (read-only)
**Groups:** Users → multiple groups → roles → inherited permissions

**Flow:** Login → JWT+Refresh+CSRF → Redis → Client | Request → Validate JWT → Check revocation → Verify CSRF → Process | Logout → Blacklist token

### Transaction Processing

**Statuses:** `pending` (synthetic only) | `approved` (no blocks) | `rejected` (blocked) | `in_review` (flagged)

**Lifecycle:** Submit (`POST /api/v1/transactions`) → Validate → Queue Redis → `202 Accepted` → Worker picks (`BRPOP`) → Evaluate rules → Save (status, matched rules, time, metadata JSONB) → Manual review if `in_review`/`pending`

**Synthetic Mode:** Separate tables (`transactions` vs `transactions_synthetic`), marked `_synthetic: true`, always `pending` (not rule-processed), controlled via `X-Synthetic-Mode` header, separate dashboard metrics

## Key Conventions

### Go Code

- Follow standard Go conventions (gofmt, golangci-lint)
- Use `internal/` packages for private API code
- Package structure: `cmd/` for executables, `internal/` for application logic, `pkg/` for shared libraries
- Error handling: Always check errors, wrap with context using `fmt.Errorf("context: %w", err)`
- Naming: camelCase for unexported, PascalCase for exported
- Structured logging with log levels (debug, info, warn, error)
- Connection pooling for all database and Redis connections
- **Validation**: All handlers MUST include request validation (mandatory practice)
- Use `context.Context` for all database and Redis operations
- Package names match directory names (feature-based)
- Interfaces defined in the same package that uses them
- Mock implementations use `Mock*` prefix and `_test.go` suffix

### TypeScript/Vue Code

- Use TypeScript strict mode
- Prefer Composition API over Options API (use `<script setup lang="ts">` syntax)
- Component naming: PascalCase for component files
- Organize by feature, not by type
- Use Pinia stores for shared state
- Tailwind utility classes for styling (avoid custom CSS when possible)
- Standardize base components for consistency (BaseButton, BaseInput, BaseTable, etc.)
- Use standardized validation fields and form components
- Remove unused code and components regularly
- Use fundamental abstractions to simplify complex logic
- Component structure: base components in `src/ui/src/components/`
- Type-safe API calls (see `src/ui/src/lib/api.ts`)
- Standardized error handling (see `src/ui/src/lib/error-handler.ts`)

### Internationalization (i18n)

**CRITICAL: All user-facing text MUST use translation keys. No hardcoded strings.**

**Usage:** `{{ $t('common.save') }}` in templates, `const { t } = useI18n()` in scripts

**Key naming (dot-notation):**
- `common.*` - Shared (buttons, actions, labels)
- `auth.*`, `header.*`, `sidebar.*` - Auth/layout
- `views.<viewName>.*` - View-specific
- `components.<componentName>.*` - Component-specific
- `errors.*`, `languages.*` - Errors and language names

**Adding translations:**
1. Add to BOTH `pt-BR.json` AND `en-US.json`
2. Descriptive paths: `auth.errors.invalidCredentials` not `err1`
3. Use parameters: `{{ $t('welcome', { name: user.name }) }}`

**Languages:** pt-BR, en-US (default: en-US) | **Details:** `src/ui/src/locales/README.md`

### Database

- Migrations: SQL files in `scripts/migrations/` (`001_schema.sql` → `002_indexes.sql` → `003_test_data.sql`)
- Use `pgxpool` for connection pooling, `pgx.Tx` for atomic operations
- **ALWAYS parameterized queries** (`$1, $2, $3`) - NEVER concatenate user input
```go
// ✅ query := "SELECT * FROM users WHERE email = $1"
// ❌ query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", userEmail)
```

### Environment Variables

- All configuration uses environment variables (see `.env.example`)
- **Required variables**: `POSTGRES_*`, `REDIS_*`, `JWT_SECRET`, `ENVIRONMENT`
- **JWT settings**: `JWT_EXPIRATION_HOURS` (default: 24), `JWT_REFRESH_EXPIRATION_HOURS` (default: 168)
- Security settings: `TLS_ENABLE`, `CORS_ALLOWED_ORIGINS`, `API_BODY_LIMIT`
- Worker settings: `WORKER_CONCURRENCY`, timeouts, retry configuration
- **MANDATORY**: Update `.env.example` when adding new environment variables
- **ALL timeout values MUST be configurable** via environment variables (no hardcoded timeouts)

## Testing Patterns

### Test Strategy

- **Unit Tests**: Follow `docs/agents/unit-test.md` | **Integration Tests**: Follow `docs/agents/integration-test.md`
- **Coverage**: 80% minimum | **Race detection**: `-race` flag | **Flaky tests**: `-count=50` (Go), `--repeat=20` (Vue)
- **AAA Pattern**: No comments, use blank lines | **Lint**: All tests must pass linters

### Unit Tests (Go)

- Mock dependencies (interfaces), test in isolation, table-driven for multiple scenarios, `-race` flag, 80%+ coverage
```go
func TestServiceOperation(t *testing.T) {
    mockRepo := &MockRepository{}
    service := NewService(mockRepo)

    result, err := service.Operation(ctx, input)

    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Integration Tests (Go)

- Build tag `//go:build integration`, Docker containers, testcontainers-go, <500ms per test
- **Isolation:** Transactions (fastest) → Truncate tables (moderate) → Recreate DB (slowest)

### UI Tests (Vue/TypeScript)

- Vitest + Vue Test Utils, Mock API with MSW, test interactions/state, `flushPromises()` for async, <300ms per test

### Coverage & Quality Verification

```bash
# Go: Lint then coverage
go vet ./... && golangci-lint run ./...
go test -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -func=coverage.txt | grep total  # ≥80%

# Vue: Lint then coverage
npm run lint && npm run test:coverage  # ≥80%

# Flaky test detection
go test -count=50 ./...           # Go
npm run test -- --repeat=20       # Vue

# Race conditions (MANDATORY)
go test -race ./...
go test -tags=integration -race ./...
```

## Important Patterns and Gotcas

### Security

**CRITICAL: OWASP Top 10 Compliance Required**

#### 1. SQL Injection Prevention
- **ALWAYS** use parameterized queries (`$1, $2, $3`)
- **NEVER** concatenate user input into SQL
```go
// ✅ GOOD: query := "SELECT * FROM users WHERE email = $1"
// ❌ BAD: query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", userEmail)
```

#### 2. XSS Prevention
**Backend:** Validate all input, set `Content-Type: application/json`, sanitize HTML with bluemonday
**Frontend:** Vue auto-escapes `{{ }}`, NEVER use `v-html` with user input, avoid `eval()`, use CSP headers

#### 3. Command Injection Prevention
- Use `exec.Command("ls", "-l", dir)` with separate args
- NEVER `exec.Command("sh", "-c", "ls " + userInput)`

#### 4. Expression Injection (Rules Engine)
- Validate field paths: `^[a-zA-Z0-9_.]+$`
- Use expr-lang type safety

#### 5. CSRF Protection (Implemented)
- Tokens in Redis (24h TTL), auto-added via API client
- **MUST use** `api.post()` not `fetch()` directly

#### 6. Authentication
**JWT:** Strong secrets (32+ bytes), expiration (24h access/168h refresh), Redis revocation
**Passwords:** bcrypt (cost 12+), never log hashes, rate limit login

#### 7. Sensitive Data
- NEVER log passwords, tokens, PII
- Use HTTPS/TLS in production
- Filter sensitive fields from API responses

#### 8. Access Control
- RBAC on ALL endpoints via middleware
- Verify user owns resource (prevent horizontal escalation)
- Admins can't deactivate themselves or last admin

#### 9. Security Misconfiguration
**Headers (MANDATORY):**
```go
c.Set("X-Content-Type-Options", "nosniff")
c.Set("X-Frame-Options", "DENY")
c.Set("X-XSS-Protection", "1; mode=block")
c.Set("Strict-Transport-Security", "max-age=31536000")
c.Set("Content-Security-Policy", "default-src 'self'")
```

#### 10. CSP
- Strict CSP: `default-src 'self'; script-src 'self'`
- No inline scripts (use nonces if needed)

#### 11. Dependency Vulnerabilities
**Scan regularly:**
```bash
govulncheck ./...          # Go
npm audit && npm audit fix # Frontend
gitleaks detect --source . # Secrets
semgrep --config auto      # SAST
```

#### 12. Security Logging
- Log auth attempts, authorization failures, critical ops
- Include correlation IDs, structured JSON format

#### 13-14. Rate Limiting & Resource Limits
- Rate limit auth endpoints (Redis-backed)
- Limit body size (4MB default), set query/Redis timeouts

#### 15. TLS/HTTPS
- **MANDATORY in production** (`TLS_ENABLE=true`)
- HSTS headers, secure cookies

#### 16. CORS
- NEVER `*` in production
- Configure via `CORS_ALLOWED_ORIGINS`

#### 17. Token Storage
- Store in Pinia (memory), NEVER localStorage
- Clear on logout

### Performance

- **Connection pooling** (`pgxpool.Pool`), **Redis pipelining** (bulk ops), **Caching** (rules 5min TTL)
- **Hot-reload** (10s interval), **Expression compilation** (cached per schema), **Race detector** (`-race`)
- **Target**: <50ms per transaction, batch processing with controlled concurrency (semaphore)

### Error Handling & Transaction Flow

- Custom error types (`src/pkg/errors/`), appropriate HTTP codes, log with context
- Worker errors → Asynq retries, generic client messages (no stack traces)
- **Status flow:** `queued` → `approved`/`in_review`/`rejected`
- **Risk levels:** Low (0-49), Medium (50-79), High (80-100)
- **Actions:** `allow`, `block` (highest priority), `review`

## Testing Anti-Patterns

**AI agents must avoid generating useless tests.**

### What NOT to Test

1. **Trivial Constructors** - Don't test constructors that only assign fields (no logic)
2. **Interface Implementation** - Compiler already validates
3. **Duplicate Validation** - Test validation ONCE at service layer, not at handler/private function layers
4. **Repetitive Tests** - Use table-driven tests instead of separate functions for similar scenarios
5. **Component Existence Only** - Verify actual prop values, not just `exists()`

### Decision Tree

```
Constructor only assigns fields? → ❌ Don't test
Test only verifies interface impl? → ❌ Don't test (compiler checks)
Validation tested in another layer? → ❌ Don't duplicate
3+ similar tests with same assertions? → ⚠️  Consolidate to table-driven
Test only checks exists() without props? → ❌ Verify actual values
Test verifies actual behavior/logic? → ✅ Write it
```

### Cleanup Checklist

- [ ] No trivial constructor/interface tests
- [ ] No duplicate validation across layers
- [ ] Repetitive tests consolidated (table-driven)
- [ ] Component tests verify props not just existence
- [ ] Tests pass linting and `-count=50` (Go) or `--repeat=20` (Vue)
- [ ] Coverage ≥80%

**Details:** See `docs/agents/unit-test.md` section 7

## Regression Prevention Rules

**CRITICAL: Follow when adding new features to prevent regressions.**

### Pre-Implementation Checklist

1. **OpenSpec Compliance** - Create proposal if required, validate with `openspec validate --strict`
2. **Impact Analysis** - Identify affected modules, document breaking changes with migration plan
3. **Test Coverage** - Plan unit (80%+), integration, API, frontend, E2E tests; run with `-race` and `-count=50`

### Mandatory Test Categories

1. **Unit Tests** - All functions, edge cases, mocked dependencies, deterministic, 80%+ coverage
2. **Integration Tests** - Real DB/Redis (testcontainers-go), transaction flows, rules engine, velocity functions
3. **API Tests** - Endpoints, auth/authz, validation, errors, status codes
4. **Frontend Tests** - Components, views, stores, router guards, i18n keys
5. **Regression Tests** - Existing tests pass, APIs work, UI renders, rules evaluate correctly

### Performance Requirements

- Transaction processing <50ms
- Optimized queries with indexes
- No N+1 queries
- Connection pool limits respected

### Pre-Merge Validation

**Code Quality:** Linters pass, tests pass, no race conditions, 80%+ coverage, security scans clean
**Documentation:** Update README, API docs, `.env.example`, archive OpenSpec proposals
**Integration:** Tests pass with real DBs, Docker builds, CI/CD passes
**Manual Testing:** Feature works, existing features work, i18n correct, performance acceptable

### AI Agent Requirements

**Before:** OpenSpec proposal, impact analysis, identify modules, plan tests
**During:** Write tests with code (TDD), run tests frequently, check for race conditions
**After:** Run full suite, verify coverage ≥80%, check linters, update docs
**Before Completion:** Verify rules followed, no breaking changes (or documented), backward compatible, <50ms latency

### Enforcement

**MANDATORY** - CI/CD automated checks, code review verification, pre-merge hooks
**Violations:** Missing tests/docs rejected, breaking changes without migration rejected, performance regressions block deployment

**Details:** See `openspec/project.md` "Regression Prevention Rules"

## Development Workflow

1. Install Git hooks: `./scripts/install-hooks.sh`
2. Make changes following vertical slice architecture
3. Add tests for new functionality (unit + integration if needed)
4. Run tests locally: `make test`
5. Run linters: `make lint`
6. Commit with conventional commit prefixes: `feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, etc.
7. Pre-commit hook runs tests and checks formatting automatically

## Known Issues & Areas for Improvement

**Layout Bugs:**
- Some layout bugs remain from Vue.js migration
- Root cause not yet identified
- Needs investigation and resolution

**Database Migrations:**
- Current migration process via scripts is not ideal
- Needs improvement with a proper migration library
- Action required: Evaluate migration tools

**OpenTelemetry:**
- Metrics infrastructure defined but not fully implemented
- Plan to use for observability and performance monitoring

## AI Agent Guidelines

**MANDATORY REQUIREMENTS:**

1. **Code Generation** - SOLID principles, vertical slice architecture, dependency injection, ISP compliance, avoid complexity, refactor magic numbers
2. **Validation** - Always add validation to handlers (mandatory)
3. **Environment Variables** - Update `.env.example`, group related vars, no hardcoded timeouts
4. **Testing** - Follow `docs/agents/unit-test.md` & `docs/agents/integration-test.md`, use `-race`, 80%+ coverage, `-count=50` for flaky tests, testcontainers-go for integration
5. **i18n** - NEVER hardcode strings, add to BOTH `pt-BR.json` AND `en-US.json`, hierarchical keys (`views.dashboard.title`)
6. **CSRF** - Use `api.post/put/patch/delete` ONLY (never `fetch()` directly), verify `X-CSRF-Token` header
7. **Code Quality** - Remove unused code, fix race conditions immediately, standardize base components
8. **Documentation** - OpenSpec for spec-driven dev, update README for significant features
9. **Code Review** - Proactively scan: guidelines adherence, tests generated/updated, patterns followed, 80%+ coverage
10. **Regression Prevention** - Follow "Regression Prevention Rules", complete checklists, verify no breaking changes, <50ms latency
11. **Security** - OWASP Top 10, security headers, never log PII, parameterized queries, validate all input, rate limiting, run scanners (`govulncheck`, `npm audit`, `gitleaks`, `semgrep`), CSP, TLS in production
12. **Anti-Patterns** - No trivial constructor tests, no interface tests, no duplicate validation, consolidate to table-driven, verify behavior not existence
13. **Coverage** - 80% minimum, verify with `go test -coverprofile=coverage.txt` and `npm run test:coverage`
14. **Lint** - All tests pass linters, fix ALL errors, no unused vars/imports/debugging code

## Additional Resources

- **Rules engine expressions**: https://github.com/expr-lang/expr
- **Fiber framework**: https://gofiber.io/
- **Asynq task queue**: https://github.com/hibiken/asynq
- **Vue.js Composition API**: https://vuejs.org/guide/extras/composition-api-faq.html
- **Vuetify components**: https://vuetifyjs.com/
- **OpenSpec documentation**: `openspec/AGENTS.md`
- **Testing guidelines**: `docs/agents/unit-test.md`, `docs/agents/integration-test.md`
- **i18n guidelines**: `src/ui/src/locales/README.md`
