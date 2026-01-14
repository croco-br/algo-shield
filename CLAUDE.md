# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AlgoShield is an open-source, high-performance fraud detection and AML transaction analysis system. It processes transactions with ultra-low latency (<50ms) using a custom expression-based rules engine powered by expr-lang. The system supports both real-time fraud prevention (pre-transaction) and post-transaction AML monitoring.

**Key Capabilities:**
- Real-time transaction processing with <50ms latency target
- Custom rules engine with hot-reload support
- Event schema management with automatic field extraction
- Risk scoring and transaction classification
- JWT-based authentication with role-based access control (RBAC)
- White-label branding customization
- Synthetic event generation for testing

## Tech Stack

**Backend (Go 1.25.4):**
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
# Run all tests (API + UI)
make test

# Run API tests only (with race detector)
make test-api

# Run API tests with gotestsum (better output)
gotestsum --format testdox -- -race -parallel 8 ./src/...

# Run specific test file
go test -race ./src/api/internal/auth/service_test.go

# Run specific test function
go test -race -run TestServiceName ./src/api/internal/auth/...

# Run integration tests (requires Docker containers)
gotestsum --format testdox -- -tags=integration -race -parallel 2 ./src/api/... ./src/workers/...

# Check for flaky tests
go test -count=50 ./src/...

# Run UI tests
make test-ui
cd src/ui && npm test

# Run UI tests with coverage
cd src/ui && npm run test:coverage

# Run UI tests interactively
cd src/ui && npm run test:ui

# Run benchmarks (rules engine)
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

All handlers and services use **dependency injection** following SOLID principles:

```go
// Service depends on interface (not concrete type)
type Service struct {
    repo Repository  // Interface, not *PostgresRepository
}

// Handler depends on interface
type Handler struct {
    service ServiceInterface  // Interface, not *Service
}

// Wiring happens in routes.go
func Setup(app *fiber.App, db *pgxpool.Pool, redis *redis.Client, ...) {
    // Create repositories (infrastructure layer)
    userRepo := user.NewPostgresUserRepository(db)

    // Create services with dependency injection (business layer)
    userService := user.NewService(userRepo, roleRepo, txManager, roleService, groupService)

    // Create handlers with dependency injection (presentation layer)
    userHandler := user.NewHandler(userService)
}
```

See `src/api/internal/routes/routes.go` for the complete dependency wiring.

### Message Queue Architecture (Asynq)

AlgoShield uses **Asynq** (Redis-based task queue) for async transaction processing:

```
API → Enqueue Transaction → Redis Queue → Worker → Process Transaction → Update DB
```

- **API**: Enqueues transactions to Redis via `LPUSH`
- **Worker**: Processes transactions from queue via `BRPOP` (blocking pop)
- **Three priority levels**: critical, default, low
- **Configurable**: Concurrency, retries, timeouts, batch size
- **Queue timeout**: 5s default (returns `ErrTimeout` if no events - expected behavior)
- **Batch processing**: Workers collect up to `batchSize` transactions and process in parallel

**Worker Processing Flow:**
1. Worker pops transaction(s) from Redis queue (`BRPOP` - blocking)
2. For batch mode: Collects up to `batchSize` transactions
3. Processes transactions in parallel (controlled concurrency via semaphore)
4. Each transaction evaluated against all enabled rules
5. Result saved to database with status, matched rules, processing time
6. Metrics recorded (success/failure, duration)
7. Worker continues to next batch (infinite loop until context cancellation)

### Rules Engine

The rules engine uses **expr-lang** for flexible expression evaluation:

**Architecture:**
- Rules defined as custom expressions (e.g., `amount > 10000`)
- Rules associated with event schemas for type safety
- Schemas define the structure of transaction events
- Hot-reload support: rules and schemas reload every 10s (configurable)
- Expression cache for performance (cleared on schema changes)
- Helper functions: `pointInPolygon()`, `velocityCount()`, `velocitySum()`

**Evaluation Flow:**
```
Transaction → Schema Validation → Rule Evaluation → Risk Score → Action (allow/block/review)
```

**Rule Actions:**
- `allow`: Sets status to `approved`
- `block`: Sets status to `rejected` (highest priority)
- `review`: Sets status to `in_review`

**Helper Functions:**
- `velocityCount(fieldPath, timeWindowSeconds)`: Count transactions by field within time window
- `velocitySum(fieldPath, timeWindowSeconds)`: Sum numeric field (auto-detected) by field
- `velocitySum(fieldPath, sumFieldPath, timeWindowSeconds)`: Sum specified field
- `pointInPolygon(lat, lon, polygon)`: Check if point is inside polygon (ray casting)

**Security:** Field paths validated (only alphanumeric, underscore, dot) to prevent SQL injection.

### Caching Strategy

- **Rules Cache**: Redis-backed with TTL (default: 5 minutes)
- **Schema Cache**: In-memory + Redis pub/sub for invalidation
- **Dashboard Metrics**: Redis-backed with 30s TTL (separate keys for synthetic/real modes)
- **Branding Cache**: Redis-backed with TTL
- **CSRF Tokens**: Redis-backed with 24-hour expiration
- **Token Revocation**: Redis sets for blacklisting JWT tokens

### Authentication & Authorization

**JWT-based authentication:**
- Access tokens: 24 hours (configurable via JWT_EXPIRATION_HOURS)
- Refresh tokens: 7 days / 168 hours (configurable via JWT_REFRESH_EXPIRATION_HOURS)
- Token revocation: Redis-backed blacklist
- CSRF protection: Required for all state-changing requests (POST/PUT/PATCH/DELETE)
- Rate limiting: Applied to login/registration endpoints

**Role-based access control (RBAC):**
- **admin**: Full system access
- **rule_editor**: Can create/update/delete rules and schemas
- **viewer**: Read-only access

**Group-based permissions:**
- Users can belong to multiple groups
- Groups can have roles assigned
- Users inherit roles from groups

**Security Flow:**
```
Login → JWT + Refresh Token + CSRF Token → Store in Redis → Return to Client
Request → Validate JWT → Check Revocation → Verify CSRF → Process
Logout → Add Token to Blacklist → Expire CSRF Token
```

### Transaction Processing

**Transaction Statuses:**
- `pending`: Awaiting processing (synthetic events only)
- `approved`: Transaction approved (no blocking rules matched)
- `rejected`: Transaction blocked (block action matched)
- `in_review`: Transaction flagged for manual review

**Transaction Lifecycle:**
1. Transaction submitted via API (`POST /api/v1/transactions`)
2. Event validated (must be non-empty JSON object)
3. Event queued in Redis
4. API returns `202 Accepted` immediately
5. Worker picks up event from queue (FIFO via `BRPOP`)
6. Rules engine evaluates event against all enabled rules
7. Result saved to database with status, matched rules, processing time
8. Full event stored in `metadata` JSONB column
9. Transactions in `in_review` or `pending` can be manually approved/rejected

**Synthetic Mode:**
- Separate tables: `transactions` (real) and `transactions_synthetic` (synthetic)
- Synthetic transactions marked with `_synthetic: true`
- Synthetic transactions always have `pending` status (not processed by rules)
- Mode controlled via `X-Synthetic-Mode` header or system config
- Dashboard metrics separated by mode

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

**CRITICAL: All user-facing text MUST use translation keys. No hardcoded strings in components.**

**Usage in templates:**
```vue
<template>
  <div>{{ $t('common.save') }}</div>
</template>
```

**Usage in scripts:**
```vue
<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const message = t('common.save')
</script>
```

**Key naming convention (dot-notation hierarchy):**
- `common.*` - Shared translations (buttons, actions, labels)
- `auth.*` - Authentication-related
- `header.*` / `sidebar.*` - Layout component translations
- `views.<viewName>.*` - View-specific translations
- `components.<componentName>.*` - Component-specific translations
- `errors.*` - Error messages
- `languages.*` - Language names

**Adding new translations:**
1. Add the key to BOTH `src/ui/src/locales/pt-BR.json` AND `src/ui/src/locales/en-US.json`
2. Use descriptive key paths: `auth.errors.invalidCredentials` not `err1`
3. Use parameters for dynamic text: `{{ $t('welcome', { name: user.name }) }}`

**Supported languages:** Portuguese (pt-BR), English (en-US)
**Default language:** English (en-US)
**For detailed i18n documentation:** See `src/ui/src/locales/README.md`

### Database

- Migrations are SQL files in `scripts/migrations/`
- Migration order: `001_schema.sql` → `002_indexes.sql` → `003_test_data.sql`
- Use `pgxpool` for connection pooling (not `database/sql`)
- **All queries MUST use parameterized statements** (prevent SQL injection)
- Transaction support via `pgx.Tx` for atomic operations
- Example (GOOD):
  ```go
  query := "SELECT * FROM users WHERE email = $1"
  row := db.QueryRow(ctx, query, userEmail)
  ```
- Example (BAD - NEVER DO THIS):
  ```go
  query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", userEmail) // VULNERABLE!
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

AlgoShield uses spec-driven testing with comprehensive documentation:

- **Unit Tests**: Follow `docs/agents/unit-test.md` strictly
- **Integration Tests**: Follow `docs/agents/integration-test.md` strictly
- **Coverage target**: **80% minimum** for new code
- **Race detection**: Always run tests with `-race` flag
- **Flaky test detection**: Run `go test -count=50` to verify deterministic tests
- **AAA Pattern**: Arrange-Act-Assert (without comments)

### Unit Tests (Go)

- Mock dependencies using interfaces (see `mock_*_test.go` files)
- Test services in isolation from repositories
- Test handlers in isolation from services
- Use table-driven tests for multiple scenarios
- Always run with `-race` flag to detect race conditions

**Example structure:**
```go
func TestServiceOperation(t *testing.T) {
    // Arrange: setup mocks
    mockRepo := &MockRepository{}
    service := NewService(mockRepo)

    // Act: call the function
    result, err := service.Operation(ctx, input)

    // Assert: verify results
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Integration Tests (Go)

- Use build tag `//go:build integration`
- Require Docker containers (postgres, redis)
- Test actual database operations
- Use `testcontainers-go` for isolated environments
- See `*_integration_test.go` files

### UI Tests (Vue/TypeScript)

- Use Vitest for unit tests
- Use Vue Test Utils for component testing
- Mock API calls with test fixtures
- Test user interactions and state management

## Important Patterns and Gotcas

### Security

**CRITICAL SECURITY REQUIREMENTS:**

1. **SQL Injection Prevention**:
   - ALWAYS use parameterized queries with `$1, $2, $3` placeholders
   - NEVER concatenate user input into SQL queries

2. **CSRF Protection**:
   - System uses CSRF tokens for all state-changing requests (POST/PUT/PATCH/DELETE)
   - **Backend**: Tokens generated on login/register, stored in Redis (24h TTL), validated via middleware
   - **Frontend**: Tokens stored in memory (Pinia), automatically added to requests via API client
   - **Common issue**: Users must logout and login after CSRF implementation (to get token)
   - **All views must use API client**: Never use `fetch()` or `axios` directly - always use `api.post()`, `api.put()`, etc.
   - **Troubleshooting**: Check if `X-CSRF-Token` header is present in DevTools Network tab

3. **Authentication**:
   - JWT tokens must be validated on every request
   - Passwords MUST be hashed with bcrypt (never log or return password hashes)
   - Implement token expiration (24h default for access tokens, 7 days for refresh tokens)
   - Support token revocation via Redis

4. **TLS/HTTPS**:
   - **TLS is required in production** (enforced at startup)
   - Configure via `TLS_ENABLE`, `TLS_CERT_PATH`, `TLS_KEY_PATH`

5. **CORS**:
   - Never use `*` in production - specify allowed origins
   - Configure via `CORS_ALLOWED_ORIGINS` environment variable

6. **Rate Limiting**:
   - Applied to auth endpoints to prevent brute force
   - Redis-backed for distributed rate limiting

7. **Request Body Size Limit**:
   - Default 4MB (configurable via `API_BODY_LIMIT`)
   - Prevents DoS attacks via large payloads

8. **Sensitive Data**:
   - NEVER log sensitive data (passwords, tokens, PII)
   - Mask sensitive fields in logs and error messages
   - Filter sensitive fields from API responses

### Performance

- **Connection pooling**: Always use `pgxpool.Pool`, never create new connections per request
- **Redis pipelining**: Use for bulk operations
- **Caching**: Rules are cached (check TTL before modifying cache logic)
- **Hot-reload**: Rules and schemas reload automatically (default: 10s interval)
- **Expression compilation**: Cached per schema to avoid recompilation overhead
- **Race detector**: Always run tests with `-race` flag
- **Target latency**: <50ms per transaction (from queue to database)
- **Batch processing**: Configurable batch size for efficient worker processing
- **Concurrent processing**: Controlled concurrency via semaphore to prevent resource exhaustion

### Error Handling

- API errors use custom error types (see `src/pkg/errors/`)
- Handlers should return appropriate HTTP status codes
- Always log errors with context for debugging
- Worker errors trigger Asynq retries (configurable max attempts)
- Return generic error messages to clients (don't expose stack traces)

### Transaction Processing

- Transactions processed asynchronously via Asynq
- Status flow: `queued` → `approved`/`in_review`/`rejected`
- Risk levels: Low (0-49), Medium (50-79), High (80-100)
- Actions: `allow`, `block`, `review` (block takes precedence)

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

**MANDATORY REQUIREMENTS FOR AI AGENTS:**

1. **Code Generation**:
   - Follow SOLID principles strictly
   - Use vertical slice architecture
   - Implement proper dependency injection
   - Ensure interfaces follow ISP (not too large)
   - Avoid unnecessary complexity
   - Use fundamental abstractions to simplify complex logic
   - Refactor magic numbers to named constants

2. **Validation**:
   - Always add validation to handlers - it's **mandatory**

3. **Environment Variables**:
   - Always update `.env.example` when adding new environment variables
   - Group related variables with clear documentation
   - All timeout values MUST be configurable (no hardcoded timeouts)

4. **Testing** (STRICTLY FOLLOW):
   - **Unit Tests**: Follow `docs/agents/unit-test.md`
   - **Integration Tests**: Follow `docs/agents/integration-test.md`
   - Use race condition flags (`-race`)
   - Minimum coverage: 80% for new code
   - Run `go test -count=50` to detect flaky tests
   - Use `testcontainers-go` for integration tests with databases

5. **Internationalization**:
   - **NEVER hardcode user-facing strings** - always use translation keys
   - Add translations to BOTH `pt-BR.json` AND `en-US.json` simultaneously
   - Use hierarchical key naming: `views.dashboard.title`, `errors.notFound`
   - Test UI with both languages to ensure text fits properly

6. **CSRF Protection**:
   - All state-changing operations MUST use `api.post/put/patch/delete`
   - Never use `fetch()` directly - it won't include CSRF token
   - Test with DevTools to verify `X-CSRF-Token` header is sent

7. **Code Quality**:
   - Remove unused files and code regularly
   - Fix concurrent issues and race conditions immediately
   - Standardize base components in front-end
   - Keep codebase clean and maintainable

8. **Documentation**:
   - Use OpenSpec for spec-driven development (see `openspec/AGENTS.md`)
   - Update README.md when adding significant features
   - Document architectural decisions in project.md

9. **Mandatory Code Review and Testing Verification**:
   - **REQUIRED**: When the user requests a correction or analysis, perform a comprehensive scan to verify:
     - Project adheres to all guidelines in this document
     - All related tests have been generated or updated
     - Code changes follow established patterns
     - Test coverage meets 80% minimum for new/modified code
     - Tests follow guidelines in `docs/agents/unit-test.md` and `docs/agents/integration-test.md`
   - This verification MUST be performed proactively
   - If violations or missing tests are found, address them before completing the task

10. **Regression Prevention**:
    - **CRITICAL**: When adding new features, follow "Regression Prevention Rules" in `openspec/project.md`
    - Complete all pre-implementation, mandatory test categories, and pre-merge checklists
    - Verify no breaking changes (or properly documented with migration plans)
    - Ensure backward compatibility maintained
    - Validate performance requirements met (<50ms latency target)

## Additional Resources

- **Rules engine expressions**: https://github.com/expr-lang/expr
- **Fiber framework**: https://gofiber.io/
- **Asynq task queue**: https://github.com/hibiken/asynq
- **Vue.js Composition API**: https://vuejs.org/guide/extras/composition-api-faq.html
- **Vuetify components**: https://vuetifyjs.com/
- **OpenSpec documentation**: `openspec/AGENTS.md`
- **Testing guidelines**: `docs/agents/unit-test.md`, `docs/agents/integration-test.md`
- **i18n guidelines**: `src/ui/src/locales/README.md`
