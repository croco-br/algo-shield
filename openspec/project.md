# Project Context

## Purpose
AlgoShield is an open-source, high-performance fraud detection and anti-money laundering (AML) transaction analysis tool designed to process transactions with ultra-low latency (<50ms). The system provides:
- Real-time fraud prevention (pre-transaction analysis)
- Post-transaction AML compliance monitoring
- Custom rules engine with hot-reload capabilities
- Risk scoring and transaction classification
- User and permission management with RBAC

## Tech Stack

### Backend (Go 1.25.5)
- **Web Framework**: Fiber v2 (high-performance HTTP server)
- **Database Driver**: pgx v5 (PostgreSQL connection pooling)
- **Cache/Queue**: go-redis v9 (Redis client)
- **Authentication**: golang-jwt v5 + bcrypt (JWT tokens with password hashing)
- **Validation**: go-playground/validator v10
- **Observability**: OpenTelemetry (metrics and tracing)
- **Concurrency**: golang.org/x/sync (worker pools)
- **Expression Engine**: expr-lang/expr v1.17.7 (custom rule expressions)

### Frontend (Vue.js 3)
- **Framework**: Vue 3.5.26 with Composition API
- **Language**: TypeScript 5.9.3
- **State Management**: Pinia 3.0.4
- **Routing**: Vue Router 4.6.4
- **UI Components**: Vuetify 3.11.6 (Material Design)
- **Icons**: Font Awesome 7.1.0 (free solid icons)
- **Syntax Highlighting**: Prism.js 1.30
- **Styling**: Tailwind CSS 4.1.18 (with PostCSS)
- **Build Tool**: Vite 7.3.1 (with chunk optimization)
- **Dev Tools**: Vue DevTools, vue-tsc

### Infrastructure
- **Database**: PostgreSQL (primary data store)
- **Cache/Queue**: Redis (message queue + rules caching)
- **Containerization**: Docker + Docker Compose
- **Deployment**: Multi-stage Docker builds

### Internationalization (i18n)
- **Library**: vue-i18n v11.2.8
- **API Mode**: Composition API mode (legacy mode removed as of January 2026)
- **Supported Languages**: Portuguese (pt-BR), English (en-US)
- **Default Language**: English (en-US)
- **Locale Files**: `src/ui/src/locales/*.json`
- **Configuration**: `src/ui/src/plugins/i18n.ts`
- **Composable**: `src/ui/src/composables/useLocale.ts`
- **Note**: Migrated from legacy mode to Composition API mode for better performance and future compatibility

## Project Conventions

### Code Style

#### Go
- Follow standard Go conventions (gofmt, golangci-lint)
- Use `internal/` packages for private API code
- Package structure: `cmd/` for executables, `internal/` for application logic, `pkg/` for shared libraries
- Error handling: Always check errors, wrap with context
- Naming: Use camelCase for unexported, PascalCase for exported
- Use structured logging (with log levels: debug, info, warn, error)
- Connection pooling for all database and Redis connections

#### TypeScript/Vue
- Use TypeScript strict mode
- Prefer Composition API over Options API
- Use `<script setup>` syntax in Vue components
- Component naming: PascalCase for component files
- Organize by feature, not by type
- Use Pinia stores for shared state
- Tailwind utility classes for styling (avoid custom CSS when possible)
- Standardize base components for consistency across the application
- Use standardized validation fields and form components
- Remove unused code and components regularly
- Use fundamental abstractions to simplify complex logic

#### Internationalization (i18n)

**All user-facing text MUST use translation keys. No hardcoded strings in components.**

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
1. Add the key to BOTH `src/ui/src/locales/pt-BR.json` and `src/ui/src/locales/en-US.json`
2. Use descriptive key paths: `auth.errors.invalidCredentials` not `err1`
3. Use parameters for dynamic text: `{{ $t('welcome', { name: user.name }) }}`

**Language switching:**
- Users change language via avatar menu in header
- Preference saved to `localStorage` and persists across sessions
- Browser locale detection on first visit (defaults to English if unsupported)

**For detailed i18n documentation, see:** `src/ui/src/locales/README.md`

### Architecture Patterns

#### Microservices Architecture
```
UI (Vue.js) → API (Fiber) → Worker (Rules Engine)
                 ↓              ↓
          PostgreSQL ← → Redis (Queue + Cache)
```

**Key Patterns:**
- **API Service**: RESTful API with JWT middleware, handles HTTP requests, authentication, and user management
- **Worker Service**: Asynchronous transaction processing with custom expression-based rules evaluation, subscribes to Redis queues
- **Hot-Reload**: Rules and schemas cached in Redis with configurable reload interval (default: 10s) and pub/sub invalidation
- **Connection Pooling**: Reusable PostgreSQL and Redis connections
- **Async Processing**: Redis pub/sub for transaction queue management
- **RBAC**: Role-based access control with users, roles, groups, and permissions
- **Event Schemas**: Schema-based event structure definition with automatic field extraction from sample JSON
- **Custom Expressions**: expr-lang based rule expressions with helper functions (velocityCount, velocitySum, pointInPolygon)

#### Design Principles
- Ultra-low latency: Target <50ms transaction processing
- Horizontal scalability: Worker processes can scale independently
- Fail-fast: Use timeouts and circuit breakers
- Idempotency: Transaction processing should be repeatable
- Security-first: JWT authentication on all protected endpoints

### Testing Strategy

**For detailed testing guidelines, refer to:**
- **Unit Tests**: `docs/agents/unit-test.md` - Comprehensive guide for Go and Vue/TypeScript unit testing
- **Integration Tests**: `docs/agents/integration-test.md` - Guide for integration testing with real dependencies

#### Go Testing
- Use standard `go test` with race detector (`-race` flag)
- Run tests in parallel (`-parallel 4`)
- Unit tests for business logic
- Integration tests for database interactions (use `//go:build integration` tag)
- Benchmarks for performance-critical code (rules engine)
- Test file naming: `*_test.go` for unit tests, `*_integration_test.go` for integration tests
- Coverage target: **80% minimum** for new code
- Use `testcontainers-go` for integration tests with real databases

#### Commands
- `make test`: Run all tests with race detector
- `make bench`: Run rules engine benchmarks
- `make lint`: Run golangci-lint
- `go test -tags=integration ./...`: Run integration tests only

### Git Workflow

#### Branch Strategy
- **Main branch**: `main` (production-ready code)
- Feature branches: Created from `main`, merged via PR
- No force pushes to main/master

#### Commit Conventions
- Use descriptive commit messages
- Prefix commits with type: `feat:`, `fix:`, `chore:`, `docs:`, etc.
- Git hooks installed via `./scripts/install-hooks.sh`
- Pre-commit checks: linting, security scanning (gitleaks, semgrep)

#### CI/CD
- GitHub Actions for CI (likely based on `.github/` directory)
- Docker builds for deployment
- Automated testing on PRs

## Domain Context

### System Features Overview

AlgoShield provides a comprehensive fraud detection and AML platform with the following major features:

1. **Authentication & Authorization** - JWT-based auth with RBAC
2. **Transaction Processing** - Async queue-based processing with worker pool
3. **Rules Engine** - Custom expression-based rule evaluation
4. **Event Schemas** - Schema definition and automatic field extraction
5. **Dashboard & Metrics** - Real-time metrics and analytics
6. **User Management** - User CRUD with role and group assignment
7. **Branding** - White-label customization
8. **Synthetic Events** - Test data generation
9. **Transaction Approval/Rejection** - Manual review workflow

### Fraud Detection & AML

#### Transaction Analysis
- **Real-time Processing**: Transactions are queued via Redis and processed asynchronously by worker pool
- **Transaction Statuses**: 
  - `pending` - Awaiting processing (synthetic events only)
  - `approved` - Transaction approved (no rules matched or allow action)
  - `rejected` - Transaction blocked (block action matched)
  - `in_review` - Transaction flagged for review (review action matched)
- **Transaction Model**:
  - `id` - UUID primary key
  - `schema_id` - Optional UUID reference to event schema (for synthetic events)
  - `status` - Transaction status (pending, approved, rejected, in_review)
  - `processing_time` - Processing time in milliseconds
  - `matched_rules` - Array of rule names that matched
  - `metadata` - JSONB column containing full event data (all fields from schema)
  - `created_at` - Timestamp when transaction was created
  - `processed_at` - Timestamp when transaction was processed (null for synthetic/pending)
- **Transaction Lifecycle**:
  1. Transaction submitted via API (`POST /api/v1/transactions`)
  2. Event validated (must be non-empty JSON object)
  3. Event queued in Redis (`transaction:queue` list)
  4. API returns `202 Accepted` immediately
  5. Worker picks up event from queue (FIFO via `BRPOP`)
  6. Rules engine evaluates event against all enabled rules
  7. Result saved to database with status, matched rules, and processing time
  8. Full event stored in `metadata` JSONB column
  9. Transactions in `in_review` or `pending` status can be manually approved/rejected via API
- **Synthetic Mode**: System supports separate synthetic and real transaction tables
  - Synthetic transactions are marked with `_synthetic: true` flag in event
  - Synthetic transactions always have `pending` status and are not processed by rules
  - Synthetic transactions have `processed_at = null` and `processing_time = 0`
  - Mode controlled via `X-Synthetic-Mode` header or system configuration
  - Dashboard metrics separated by mode (cached separately)
  - Transaction queries respect mode (filter by table)
  - Separate database tables: `transactions` (real) and `transactions_synthetic` (synthetic)

#### Rule System
- **Custom Expression-Based**: All rules use `custom_expression` condition with [expr-lang](https://github.com/expr-lang/expr)
- **Schema-Based Evaluation**: Rules must reference an event schema for type safety
- **Rule Structure**:
  - `name` - Rule identifier
  - `description` - Human-readable description
  - `action` - One of: `allow`, `block`, `review`
  - `priority` - Numeric priority (higher = evaluated first, though all rules are evaluated)
  - `enabled` - Boolean flag to enable/disable rule
  - `schema_id` - UUID reference to event schema (required)
  - `conditions.custom_expression` - Expression string (required)
  - `score` - Risk score (0-100, stored but not currently used in status determination)
- **Rule Actions**:
  - `allow` - Sets status to `approved` (if no block action matches)
  - `block` - Sets status to `rejected` (highest priority, overrides other actions)
  - `review` - Sets status to `in_review` (if not already rejected)
- **Rule Evaluation**:
  - All enabled rules are evaluated for each transaction
  - Rules are evaluated in priority order (though all are checked)
  - Multiple rules can match; status determined by action priority (block > review > allow)
  - Matched rule names are stored in `matched_rules` array
- **Expression Syntax**:
  - Supports comparisons (`>`, `<`, `==`, `!=`, `>=`, `<=`)
  - Logical operators (`and`, `or`, `not`)
  - Array operations (`in`, `contains`)
  - Nested field access via dot notation (`user.country`, `metadata.ip_address`)
  - String operations and type coercion
- **Helper Functions**:
  - `velocityCount(fieldPath, timeWindowSeconds)` - Count transactions grouped by field value within time window
  - `velocitySum(fieldPath, timeWindowSeconds)` - Sum numeric field (auto-detected) grouped by field value within time window
  - `velocitySum(fieldPath, sumFieldPath, timeWindowSeconds)` - Sum specified numeric field grouped by field value
  - `pointInPolygon(lat, lon, polygon)` - Check if geographic point is inside polygon (ray casting algorithm)
- **Velocity Functions Details**:
  - Field path can be any schema field (e.g., `"origin"`, `"user.id"`, `"customer_id"`)
  - Time window in seconds (e.g., 3600 for 1 hour)
  - Queries transaction history from database `metadata` JSONB column
  - Supports nested field paths with dot notation
  - Returns 0 on error (graceful degradation)
  - Field path validation: Only alphanumeric, underscore, and dot characters allowed (prevents SQL injection)
  - Queries use PostgreSQL JSONB operators (`->` and `->>`) for safe field extraction
  - Time window calculated from current transaction's `created_at` timestamp
  - Only queries processed transactions (excludes pending synthetic events)
  - Count/Sum operations grouped by field value within time window
- **Hot-Reload**: Rules cached in Redis with configurable reload interval (default: 10s)
  - Schema changes trigger immediate cache invalidation via Redis pub/sub
  - Expression cache cleared on reload to ensure fresh compilation

#### Event Schemas
- **Purpose**: Define transaction event structure with automatic field extraction
- **Schema Structure**:
  - `name` - Unique schema identifier
  - `description` - Human-readable description
  - `sample_json` - Example JSON event (used for field extraction)
  - `extracted_fields` - Auto-generated array of field definitions
- **Field Extraction**:
  - Automatically extracts all fields from `sample_json`
  - Field definition includes: `path` (dot-notation), `type` (string, number, boolean, datetime, array, object, null), `nullable` flag
  - Supports nested objects and arrays
  - Can be re-extracted via `POST /api/v1/schemas/:id/parse` endpoint
- **Schema Operations**:
  - Create: `POST /api/v1/schemas` (requires `rule_editor` or `admin` role)
  - List: `GET /api/v1/schemas`
  - Get: `GET /api/v1/schemas/:id`
  - Update: `PUT /api/v1/schemas/:id` (requires `rule_editor` or `admin` role)
  - Delete: `DELETE /api/v1/schemas/:id` (requires `rule_editor` or `admin` role, fails if referenced by rules)
  - Parse: `POST /api/v1/schemas/:id/parse` - Re-extract fields from sample JSON
  - Generate Events: `POST /api/v1/schemas/:id/generate-events` - Generate synthetic events
- **Schema Validation**:
  - Rules must reference a valid schema
  - Schema deletion prevented if referenced by rules (returns list of referencing rules)
  - Schema name must be unique
- **Cache Invalidation**: Schema changes trigger Redis pub/sub notification to invalidate worker caches

#### Synthetic Event Generation
- **Purpose**: Generate test transaction events from schemas
- **Generation Process**:
  1. Operator requests N events from a schema via `POST /api/v1/schemas/:id/generate-events`
  2. System generates N events with random values based on field types
  3. Each event marked with `_schema_id` and `_synthetic: true`
  4. Events queued in Redis for processing
  5. Worker processes events but marks them as `pending` (not evaluated by rules)
- **Value Generation**:
  - String fields: Random alphanumeric strings (10-30 chars)
  - Number fields: Random numbers (0-10000)
  - Boolean fields: Random true/false
  - DateTime fields: Random RFC3339 dates within last 90 days
  - Date-like string fields (contains "date", "time", "timestamp", etc.): Random dates
  - Array fields: Empty arrays
  - Object fields: Empty objects
  - Null fields: null
- **Event Structure**: Generated events match schema's `extracted_fields` structure
- **Response**: Returns count of successfully generated and queued events

### Authentication & Authorization

#### Authentication
- **JWT-Based**: All protected endpoints require JWT token in `Authorization: Bearer <token>` header
- **Public Endpoints**:
  - `POST /api/v1/auth/register` - User registration
  - `POST /api/v1/auth/login` - User login (returns JWT token)
  - `GET /api/v1/branding` - Get branding configuration
  - `GET /health` - Health check
  - `GET /ready` - Readiness check
- **Protected Endpoints**: All other endpoints require valid JWT token
- **Token Management**:
  - Tokens stored in Redis for revocation support
  - Logout invalidates token via Redis
  - Token expiration configurable (default: 24 hours)
- **Password Security**:
  - Passwords hashed with bcrypt
  - Never returned in API responses
  - Password hash stored in database, never logged

#### Role-Based Access Control (RBAC)
- **Roles**:
  - `admin` - Full system access (user management, all rules, branding, system config)
  - `rule_editor` - Can create/update/delete rules and schemas
  - `viewer` - Read-only access (can view transactions, rules, schemas, dashboard)
- **Role Assignment**:
  - Users can have multiple roles
  - Roles assigned via `POST /api/v1/permissions/users/:userId/roles` (admin only)
  - Roles removed via `DELETE /api/v1/permissions/users/:userId/roles/:roleId` (admin only)
- **Groups**:
  - Users can belong to multiple groups
  - Groups can have roles assigned
  - Group roles inherited by group members
  - Groups managed via `/api/v1/groups` endpoints (admin only)
- **Permission Enforcement**:
  - Middleware checks roles on protected endpoints
  - `RequireRole("admin")` - Requires exact role
  - `RequireAnyRole("admin", "rule_editor")` - Requires one of the roles
  - Frontend routes protected via router guards
  - Admin-only routes redirect non-admin users to dashboard

#### User Management (Admin Only)
- **User Operations**:
  - List users: `GET /api/v1/permissions/users`
  - Get user: `GET /api/v1/permissions/users/:id`
  - Update active status: `PUT /api/v1/permissions/users/:id/active`
  - Assign role: `POST /api/v1/permissions/users/:userId/roles`
  - Remove role: `DELETE /api/v1/permissions/users/:userId/roles/:roleId`
- **User Active Status**:
  - Users can be activated/deactivated
  - Inactive users cannot login
  - Protection: Admin cannot deactivate themselves
  - Protection: Cannot deactivate last active admin
- **User Model**:
  - `id` - UUID
  - `email` - Unique email address
  - `name` - Display name
  - `password_hash` - Bcrypt hash (never returned)
  - `auth_type` - `local` or `sso` (SSO not yet implemented)
  - `active` - Boolean status
  - `roles` - Array of assigned roles
  - `groups` - Array of assigned groups
  - `created_at`, `updated_at`, `last_login_at` - Timestamps

#### Roles Management (Admin Only)
- **Role Endpoints**:
  - List roles: `GET /api/v1/roles` - Returns all roles in system
  - Get role: `GET /api/v1/roles/:id` - Returns role details with timestamps
- **Role Model**:
  - `id` - UUID
  - `name` - Unique role name (e.g., "admin", "rule_editor", "viewer")
  - `description` - Human-readable description
  - `created_at`, `updated_at` - Timestamps
- **Role Assignment to Users**:
  - Assign: `POST /api/v1/permissions/users/:userId/roles` (body: `{"role_id": "uuid"}`)
  - Remove: `DELETE /api/v1/permissions/users/:userId/roles/:roleId`
  - Users can have multiple roles
  - Role permissions are additive (user has union of all role permissions)

#### Groups Management (Admin Only)
- **Group Endpoints**:
  - List groups: `GET /api/v1/groups` - Returns all groups with their roles
  - Get group: `GET /api/v1/groups/:id` - Returns group details with assigned roles
- **Group Model**:
  - `id` - UUID
  - `name` - Unique group name
  - `description` - Human-readable description
  - `roles` - Array of roles assigned to the group (inherited by all members)
  - `created_at`, `updated_at` - Timestamps
- **Group-Role Inheritance**:
  - Groups can have multiple roles assigned
  - All members of a group inherit all roles assigned to that group
  - User permissions = direct roles + group roles (union)
  - Groups provide convenient way to manage permissions for multiple users
- **User-Group Relationship**:
  - Users can belong to multiple groups
  - Groups can have multiple users
  - Many-to-many relationship managed in database

### Transaction Processing

#### API Endpoints
- **Process Transaction**: `POST /api/v1/transactions`
  - Accepts JSON event object
  - Validates event is non-empty JSON object
  - Queues event in Redis (`transaction:queue`)
  - Returns `202 Accepted` with status "queued"
  - Extracts `external_id` from event if present
- **List Transactions**: `GET /api/v1/transactions`
  - Query parameters: `limit` (default: 50, max: 100), `offset` (default: 0)
  - Filters: `status`, `schema_id`, `start_date` (RFC3339), `end_date` (RFC3339), `min_amount`, `max_amount`
  - Returns paginated list with total count
  - Respects synthetic mode (separate tables)
- **Get Transaction**: `GET /api/v1/transactions/:id`
  - Returns full transaction details including metadata
  - Respects synthetic mode
- **Approve Transaction**: `PATCH /api/v1/transactions/:id/approve`
  - Changes status from `pending` or `in_review` to `approved`
  - Returns updated transaction
  - Fails if transaction not in reviewable status
- **Reject Transaction**: `PATCH /api/v1/transactions/:id/reject`
  - Changes status from `pending` or `in_review` to `rejected`
  - Returns updated transaction
  - Fails if transaction not in reviewable status

#### Worker Processing
- **Architecture**: Worker pool processes transactions from Redis queue
- **Queue Service**:
  - Redis list: `transaction:queue` (FIFO queue)
  - API pushes events via `LPUSH` (left push, adds to head)
  - Workers pop events via `BRPOP` (blocking right pop, removes from tail)
  - Queue pop timeout: Configurable timeout (default: 5s) - returns `ErrTimeout` if no events available (expected behavior)
  - Events serialized as JSON before queuing
  - Invalid JSON or unmarshal errors return `ErrInvalidData`
- **Worker Configuration**:
  - Concurrency: Configurable number of workers (default: 10)
  - Batch size: Configurable batch processing (default: 50 transactions per batch)
  - Timeouts: Configurable transaction timeout (default: 300ms) and rule evaluation timeout (default: 300ms)
  - Queue pop timeout: Configurable timeout for queue operations (default: 5s)
- **Processing Flow**:
  1. Worker pops transaction(s) from Redis queue using `BRPOP` (blocking, waits for events)
  2. For batch mode: Collects up to `batchSize` transactions (stops on timeout or batch full)
  3. Processes transactions in parallel (controlled concurrency via semaphore)
  4. Each transaction evaluated against all enabled rules
  5. Result saved to database with status, matched rules, processing time
  6. Metrics recorded (success/failure, duration)
  7. Worker continues to next batch/transaction (infinite loop until context cancellation)
- **Batch Processing**:
  - When `batchSize > 1`, workers process transactions in batches
  - Batch collected from queue (up to `batchSize` items)
  - Transactions in batch processed in parallel with controlled concurrency
  - Uses `golang.org/x/sync/semaphore` for concurrency control
  - Uses `golang.org/x/sync/errgroup` for error handling
  - Individual transaction failures don't stop batch processing
- **Retry Mechanism**:
  - Exponential backoff retry for failed transactions
  - Configurable retry attempts and delays
  - Retries on transient errors
- **Metrics Collection**:
  - Tracks total processed, total failed, average duration
  - Metrics per transaction and per batch
  - Logged on worker shutdown

#### Rules Reload
- **Periodic Reload**: Rules reloaded periodically (configurable interval, default: 10s)
- **Schema Invalidation**: Schema changes trigger immediate cache invalidation via Redis pub/sub
- **Cache Management**:
  - Rules cached in Redis with TTL
  - Schemas cached in-memory with Redis pub/sub invalidation
  - Expression cache cleared on reload to ensure fresh compilation

### Dashboard & Metrics

#### Dashboard Endpoints
- **Get Metrics**: `GET /api/v1/dashboard/metrics`
  - Returns comprehensive dashboard metrics
  - Cached in Redis (30s TTL)
  - Separate cache keys for synthetic and real modes
- **Metrics Data**:
  - `status_distribution` - Count of transactions per status
  - `temporal_24h` - Transaction count per hour for last 24 hours
  - `temporal_7d` - Transaction count per day for last 7 days
  - `temporal_30d` - Transaction count per day for last 30 days
  - `total_count` - Total transaction count
  - `cached_at` - Timestamp of cache generation
- **Cache Invalidation**: Cache invalidated on transaction creation (via Redis pub/sub or direct call)

#### Frontend Dashboard
- **Views**: Dashboard, Transactions, Rules, Schemas, Permissions, Branding, Profile
- **Navigation**:
  - Sidebar navigation with role-based menu items
  - Admin-only items hidden from non-admin users
  - Responsive design with mobile overlay
  - Collapsible sidebar
- **Features**:
  - Real-time metrics visualization
  - Transaction filtering and search
  - Rule management with expression builder
  - Schema management with field extraction
  - User and permission management (admin only)
  - Branding customization (admin only)
  - Profile management
  - Internationalization (i18n) with language switching
- **Dashboard View**:
  - Total transaction count
  - Status distribution chart
  - Temporal charts (24h, 7d, 30d)
  - Auto-refresh with configurable polling interval
  - Separate metrics for synthetic mode
- **Transaction View**:
  - Live transaction feed (polling-based, configurable interval)
  - Filter by status, schema, date range, amount range
  - Schema-based dynamic column generation (columns adapt to selected schema)
  - Schema selection required to view transactions (filters transactions by schema)
  - Approve/reject actions for reviewable transactions (pending/in_review status)
  - Pagination support (limit/offset)
  - Real-time updates with pause/resume controls
  - Transaction details modal
- **Rules View**:
  - List all rules with status indicators (enabled/disabled)
  - Create/edit rules with expression builder
  - Visual builders for common patterns:
    - Condition builder (field, operator, value with logic operators)
    - Velocity builder (count/sum with time window and threshold)
    - Polygon builder (geographic point-in-polygon checks)
  - Expression validation and sanitization
  - Schema selection for rule creation (required)
  - Rule priority and action configuration
  - Enable/disable rules
  - Delete rules (with confirmation)
- **Schemas View**:
  - List all schemas with field counts
  - Create/edit schemas with sample JSON editor
  - JSON syntax highlighting (Prism.js)
  - Parse schema to extract fields (re-extract from sample JSON)
  - Generate synthetic events from schema (configurable count)
  - View schema details and extracted fields
  - Delete schemas (fails if referenced by rules, shows referencing rules)
- **Permissions View** (Admin Only):
  - List all users with roles and groups
  - View user details
  - Activate/deactivate users (with protection: cannot deactivate self or last admin)
  - Assign/remove roles to/from users
  - User search and filtering
- **Branding View** (Admin Only):
  - Configure app name, colors, icons
  - Preview branding changes
  - Save branding configuration
- **Profile View**:
  - View current user information
  - View assigned roles and groups
  - Language preference (stored in localStorage)
- **Authentication**:
  - Login page with email/password
  - Registration page (public)
  - Protected routes with automatic redirect to login
  - JWT token stored in Pinia store
  - Automatic token refresh handling
  - Logout functionality

### Branding

#### Branding Configuration
- **Purpose**: White-label customization of UI appearance
- **Configuration Fields**:
  - `app_name` - Application name displayed in UI
  - `icon_url` - URL to application icon
  - `favicon_url` - URL to favicon
  - `primary_color` - Primary theme color (hex)
  - `secondary_color` - Secondary theme color (hex)
  - `header_color` - Header background color (hex)
- **Endpoints**:
  - Get: `GET /api/v1/branding` (public, no auth required)
  - Update: `PUT /api/v1/branding` (admin only)
- **Caching**: Branding cached in Redis with TTL
- **Frontend Integration**: Branding store loads and applies configuration on app initialization

### System Configuration

#### Synthetic Mode
- **Purpose**: Separate synthetic and real transaction data
- **Configuration**:
  - Get mode: `GET /api/v1/system/mode` (any authenticated user)
  - Set mode: `PUT /api/v1/system/mode` (admin only)
- **Mode Behavior**:
  - When enabled, system uses separate transaction table for synthetic events
  - Dashboard metrics separated by mode
  - Transaction queries respect mode (via `X-Synthetic-Mode` header or system config)
  - Mode stored in Redis and propagated via context
- **Header Support**: `X-Synthetic-Mode: true` header can override system mode per request

### API Endpoints Summary

#### Public Endpoints (No Authentication)
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/branding` - Get branding configuration
- `GET /health` - Health check (database and Redis connectivity)
- `GET /ready` - Readiness check

#### Protected Endpoints (Require JWT)
- **Authentication**:
  - `GET /api/v1/auth/me` - Get current user
  - `POST /api/v1/auth/logout` - Logout (invalidate token)
- **Transactions**:
  - `POST /api/v1/transactions` - Process transaction (queue for worker)
  - `GET /api/v1/transactions` - List transactions (with filters and pagination)
  - `GET /api/v1/transactions/:id` - Get transaction details
  - `PATCH /api/v1/transactions/:id/approve` - Approve transaction
  - `PATCH /api/v1/transactions/:id/reject` - Reject transaction
- **Dashboard**:
  - `GET /api/v1/dashboard/metrics` - Get dashboard metrics
- **Rules** (Read: all authenticated, Write: `rule_editor` or `admin`):
  - `GET /api/v1/rules` - List all rules
  - `GET /api/v1/rules/:id` - Get rule details
  - `POST /api/v1/rules` - Create rule
  - `PUT /api/v1/rules/:id` - Update rule
  - `DELETE /api/v1/rules/:id` - Delete rule
- **Schemas** (Read: all authenticated, Write: `rule_editor` or `admin`):
  - `GET /api/v1/schemas` - List all schemas
  - `GET /api/v1/schemas/:id` - Get schema details
  - `POST /api/v1/schemas` - Create schema
  - `PUT /api/v1/schemas/:id` - Update schema
  - `DELETE /api/v1/schemas/:id` - Delete schema
  - `POST /api/v1/schemas/:id/parse` - Re-extract fields from sample JSON
  - `POST /api/v1/schemas/:id/generate-events` - Generate synthetic events
- **Permissions** (Admin only):
  - `GET /api/v1/permissions/users` - List all users
  - `GET /api/v1/permissions/users/:id` - Get user details
  - `PUT /api/v1/permissions/users/:id/active` - Update user active status
  - `POST /api/v1/permissions/users/:userId/roles` - Assign role to user
  - `DELETE /api/v1/permissions/users/:userId/roles/:roleId` - Remove role from user
- **Roles** (Admin only):
  - `GET /api/v1/roles` - List all roles
  - `GET /api/v1/roles/:id` - Get role details
- **Groups** (Admin only):
  - `GET /api/v1/groups` - List all groups
  - `GET /api/v1/groups/:id` - Get group details
- **Branding** (Admin only):
  - `PUT /api/v1/branding` - Update branding configuration
- **System**:
  - `GET /api/v1/system/mode` - Get synthetic mode status (all authenticated)
  - `PUT /api/v1/system/mode` - Set synthetic mode (admin only)

### Middleware & Request Processing

#### Middleware Stack
1. **Logger** - Request logging with structured format
2. **Security Headers** - CORS and security headers (Brave browser compatibility)
3. **CORS** - Cross-origin resource sharing configuration
4. **Synthetic Mode** - Extracts `X-Synthetic-Mode` header and sets context
5. **Auth Middleware** - Validates JWT token on protected routes
6. **Role Middleware** - Enforces role-based access control

#### Request Validation
- **Input Validation**: All request bodies validated using `go-playground/validator`
- **Pagination Validation**: Limit (max 100) and offset validated
- **Date Validation**: RFC3339 format required for date filters
- **UUID Validation**: All ID parameters validated as UUIDs
- **Error Responses**: Consistent JSON error format with descriptive messages

#### Context Propagation
- **Synthetic Mode**: Propagated via context from middleware to repositories
- **User Context**: Current user ID available in context after auth middleware
- **Timeout Context**: All handlers use timeout contexts (default: 30s)

### Performance Requirements
- **Target Latency**: <50ms per transaction (from queue to database)
- **Throughput**: High-volume transaction processing with horizontal scaling
- **Scalability**: Worker processes can scale independently
- **Reliability**: Retry mechanisms with exponential backoff
- **Caching**: Aggressive caching of rules, schemas, and dashboard metrics
- **Connection Pooling**: PostgreSQL and Redis connection pooling for optimal performance
- **Batch Processing**: Configurable batch size for efficient worker processing
- **Concurrent Processing**: Controlled concurrency via semaphore to prevent resource exhaustion

## Code Security Guidelines

### Backend Security (Go)

#### OWASP Top 10 Prevention

**1. Injection Attacks (SQL Injection, NoSQL Injection, Command Injection)**

- **SQL Injection Prevention**:
  - ALWAYS use parameterized queries with pgx prepared statements
  - NEVER concatenate user input directly into SQL queries
  - Use `$1, $2, $3` placeholders for all dynamic values
  - Example (GOOD):
    ```go
    query := "SELECT * FROM users WHERE email = $1"
    row := db.QueryRow(ctx, query, userEmail)
    ```
  - Example (BAD - NEVER DO THIS):
    ```go
    query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", userEmail) // VULNERABLE!
    ```

- **Command Injection Prevention**:
  - Avoid executing shell commands when possible
  - If shell commands are necessary, use `exec.Command` with separate arguments (not shell strings)
  - Validate and sanitize all input used in system commands
  - Use allowlists for command arguments
  - Example (GOOD):
    ```go
    cmd := exec.Command("ls", "-l", safeDirectory)
    ```
  - Example (BAD):
    ```go
    cmd := exec.Command("sh", "-c", "ls -l " + userInput) // VULNERABLE!
    ```

- **Expression Injection Prevention** (Rules Engine):
  - Validate expression syntax before compilation
  - Field path validation: Only allow alphanumeric, underscore, and dot characters
  - Use expr-lang's type safety features
  - Sanitize field paths in velocity functions
  - Example (GOOD):
    ```go
    if !isValidFieldPath(fieldPath) {
        return 0 // Reject invalid field paths
    }
    ```

**2. Broken Authentication**

- **JWT Token Security**:
  - Use strong secret keys (minimum 32 bytes, random)
  - Store secrets in environment variables, never in code
  - Implement token expiration (default: 24 hours)
  - Support token revocation via Redis
  - Validate token signature and expiration on every request
  - Example:
    ```go
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method")
        }
        return []byte(config.JWTSecret), nil
    })
    ```

- **Password Security**:
  - ALWAYS use bcrypt for password hashing (cost factor: 12+)
  - NEVER log or return password hashes in API responses
  - Enforce password complexity requirements (minimum length, character types)
  - Implement rate limiting on login attempts
  - Example:
    ```go
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    ```

- **Session Management**:
  - Store active tokens in Redis for quick invalidation
  - Implement logout functionality that invalidates tokens
  - Use secure session IDs (UUIDs)
  - Clear sessions on password change

**3. Sensitive Data Exposure**

- **Data Protection**:
  - NEVER log sensitive data (passwords, tokens, credit cards, PII)
  - Mask sensitive fields in logs and error messages
  - Use HTTPS/TLS for all production traffic (configure TLS in Fiber)
  - Encrypt sensitive data at rest in database (if required by compliance)
  - Return minimal error details to clients (avoid stack traces)
  - Example:
    ```go
    // GOOD
    log.Error().Str("user_id", userID).Msg("authentication failed")

    // BAD - NEVER DO THIS
    log.Error().Str("password", password).Msg("authentication failed") // VULNERABLE!
    ```

- **API Response Security**:
  - Filter sensitive fields from API responses (password_hash, internal IDs)
  - Use separate DTOs for API responses (don't return database models directly)
  - Implement field-level access control

**4. XML External Entities (XXE)**

- **Not applicable**: AlgoShield does not process XML
- If XML processing is added in the future, disable external entity processing

**5. Broken Access Control**

- **RBAC Enforcement**:
  - Implement role-based access control on ALL protected endpoints
  - Use middleware to enforce permissions before handler execution
  - Validate user permissions in service layer (defense in depth)
  - Prevent privilege escalation (users can't modify their own roles)
  - Example:
    ```go
    // Middleware check
    func RequireRole(role string) fiber.Handler {
        return func(c *fiber.Ctx) error {
            user := c.Locals("user").(*User)
            if !user.HasRole(role) {
                return fiber.ErrForbidden
            }
            return c.Next()
        }
    }
    ```

- **Authorization Checks**:
  - Verify user owns resource before allowing access/modification
  - Implement object-level authorization (not just endpoint-level)
  - Prevent horizontal privilege escalation (user A can't access user B's data)
  - Example:
    ```go
    // Check ownership before update
    if transaction.UserID != currentUser.ID && !currentUser.IsAdmin() {
        return fiber.ErrForbidden
    }
    ```

- **Admin Protection**:
  - Admins cannot deactivate themselves
  - Cannot deactivate last active admin
  - Require re-authentication for critical operations

**6. Security Misconfiguration**

- **Secure Defaults**:
  - Disable debug mode in production
  - Set secure HTTP headers (via middleware)
  - Configure CORS properly (restrict origins in production)
  - Use environment-specific configurations
  - Remove default credentials and test accounts

- **Security Headers** (implement via Fiber middleware):
  ```go
  app.Use(func(c *fiber.Ctx) error {
      c.Set("X-Content-Type-Options", "nosniff")
      c.Set("X-Frame-Options", "DENY")
      c.Set("X-XSS-Protection", "1; mode=block")
      c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
      c.Set("Content-Security-Policy", "default-src 'self'")
      return c.Next()
  })
  ```

- **Error Handling**:
  - Return generic error messages to clients
  - Log detailed errors internally with correlation IDs
  - Don't expose stack traces or internal paths

**7. Cross-Site Scripting (XSS)**

- **Input Validation**:
  - Validate all input using `go-playground/validator`
  - Sanitize HTML input if accepting rich text (use bluemonday library)
  - Escape output in templates (if using HTML templates)
  - Store data as plain text, escape on output

- **API Security**:
  - Set `Content-Type: application/json` header
  - Validate JSON structure before processing
  - Reject malformed requests early

**8. Insecure Deserialization**

- **JSON Security**:
  - Use standard library `encoding/json` for serialization
  - Validate JSON structure against expected schema
  - Set size limits for request bodies
  - Don't deserialize untrusted binary data
  - Example:
    ```go
    // Set max body size in Fiber
    app.Use(bodyparser.New(bodyparser.Config{
        BodyLimit: 1 * 1024 * 1024, // 1MB limit
    }))
    ```

**9. Using Components with Known Vulnerabilities**

- **Dependency Management**:
  - Keep all dependencies up to date
  - Run `go mod tidy` regularly
  - Use Dependabot for automated updates
  - Scan for vulnerabilities with `govulncheck`
  - Review security advisories for critical dependencies
  - Pin dependency versions (don't use `latest`)

- **Security Scanning**:
  - Run `gitleaks` to detect secrets in code
  - Run `semgrep` for static security analysis
  - Integrate security scans in CI/CD pipeline
  - Example commands:
    ```bash
    go install golang.org/x/vuln/cmd/govulncheck@latest
    govulncheck ./...
    ```

**10. Insufficient Logging & Monitoring**

- **Security Logging**:
  - Log all authentication attempts (success and failure)
  - Log authorization failures
  - Log critical operations (user creation, role assignment, config changes)
  - Include correlation IDs for request tracing
  - Use structured logging (JSON format)
  - Example:
    ```go
    log.Warn().
        Str("user_id", userID).
        Str("ip_address", c.IP()).
        Str("endpoint", c.Path()).
        Msg("unauthorized access attempt")
    ```

- **Monitoring**:
  - Monitor failed authentication attempts (detect brute force)
  - Monitor unusual transaction patterns
  - Set up alerts for security events
  - Track rate of 401/403 responses
  - Use OpenTelemetry for observability

#### Additional Backend Security Practices

**Input Validation**:
- Validate ALL input in handlers using `go-playground/validator`
- Check data types, lengths, formats, and ranges
- Use allowlists instead of blocklists
- Validate business logic constraints
- Example:
  ```go
  type CreateUserRequest struct {
      Email    string `json:"email" validate:"required,email"`
      Password string `json:"password" validate:"required,min=8"`
      Name     string `json:"name" validate:"required,min=2,max=100"`
  }
  ```

**Rate Limiting**:
- Implement rate limiting on authentication endpoints
- Use Redis for distributed rate limiting
- Different limits for authenticated vs unauthenticated users
- Return 429 Too Many Requests when limit exceeded

**Timeout & Resource Limits**:
- Set timeouts for all database queries
- Set timeouts for Redis operations
- Limit request body size
- Set maximum connection pool sizes
- Prevent resource exhaustion attacks

**Database Security**:
- Use connection pooling (pgx)
- Implement query timeouts
- Use read replicas for analytics queries
- Encrypt database connections (TLS)
- Follow principle of least privilege for database users

**Redis Security**:
- Use authentication (requirepass)
- Use TLS for Redis connections (production)
- Set memory limits to prevent OOM
- Use separate Redis instances for queue and cache

**CORS Configuration**:
- Restrict allowed origins (don't use `*` in production)
- Specify allowed methods explicitly
- Don't allow credentials with wildcard origins
- Example:
  ```go
  app.Use(cors.New(cors.Config{
      AllowOrigins: "https://yourdomain.com",
      AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
      AllowHeaders: "Authorization,Content-Type",
      AllowCredentials: true,
  }))
  ```

### Frontend Security (Vue.js/TypeScript)

#### OWASP Top 10 for Frontend

**1. Cross-Site Scripting (XSS) Prevention**

- **Template Security**:
  - Vue.js escapes expressions automatically (using `{{ }}`)
  - NEVER use `v-html` with user input
  - If `v-html` is necessary, sanitize with DOMPurify
  - Example (GOOD):
    ```vue
    <template>
      <div>{{ userInput }}</div> <!-- Automatically escaped -->
    </template>
    ```
  - Example (BAD):
    ```vue
    <template>
      <div v-html="userInput"></div> <!-- VULNERABLE if userInput contains scripts -->
    </template>
    ```

- **JavaScript Security**:
  - Avoid `eval()`, `Function()`, `setTimeout(string)`, `setInterval(string)`
  - Don't create DOM elements from user input
  - Sanitize user input before storing or displaying
  - Use Content Security Policy (CSP) headers

**2. Injection Attacks Prevention**

- **API Request Security**:
  - Use Axios or Fetch with proper encoding
  - Don't construct URLs by string concatenation
  - Use URL parameters properly
  - Validate data before sending to API
  - Example (GOOD):
    ```typescript
    const response = await axios.get('/api/users', {
      params: { id: userId } // Properly encoded
    })
    ```

**3. Broken Authentication (Frontend)**

- **Token Storage**:
  - Store JWT tokens in memory (Pinia store) for SPA security
  - DON'T store tokens in localStorage (vulnerable to XSS)
  - If localStorage is necessary, use HttpOnly cookies instead
  - Clear tokens on logout
  - Implement token refresh mechanism
  - Example:
    ```typescript
    // GOOD - Store in Pinia
    const authStore = useAuthStore()
    authStore.setToken(token)

    // BAD - Don't store sensitive tokens in localStorage
    localStorage.setItem('token', token) // VULNERABLE to XSS
    ```

- **Session Management**:
  - Implement automatic logout on token expiration
  - Redirect to login on 401 responses
  - Clear all state on logout

**4. Sensitive Data Exposure (Frontend)**

- **Data Handling**:
  - Never log sensitive data to console
  - Don't store sensitive data in localStorage/sessionStorage
  - Clear sensitive forms on unmount
  - Mask sensitive input fields (passwords, credit cards)
  - Example:
    ```vue
    <template>
      <input type="password" v-model="password" autocomplete="new-password">
    </template>
    ```

- **API Response Handling**:
  - Validate API responses before using data
  - Handle errors gracefully without exposing details
  - Don't display raw error messages to users

**5. Broken Access Control (Frontend)**

- **Route Protection**:
  - Implement router guards for protected routes
  - Check user roles before rendering admin components
  - Hide UI elements based on permissions
  - Example:
    ```typescript
    router.beforeEach((to, from, next) => {
      const authStore = useAuthStore()

      if (to.meta.requiresAuth && !authStore.isAuthenticated) {
        next('/login')
      } else if (to.meta.requiresAdmin && !authStore.isAdmin) {
        next('/dashboard')
      } else {
        next()
      }
    })
    ```

- **Component-Level Protection**:
  - Check permissions in components
  - Don't just hide UI elements (also check in backend)
  - Example:
    ```vue
    <template>
      <button v-if="userHasPermission('admin')" @click="deleteUser">
        Delete User
      </button>
    </template>
    ```

**6. Security Misconfiguration (Frontend)**

- **Build Security**:
  - Remove console.log statements in production
  - Enable source map only for debugging (not in production)
  - Set proper CSP headers
  - Configure HTTPS redirects
  - Use SRI (Subresource Integrity) for CDN resources

- **Vite Configuration**:
  ```typescript
  // vite.config.ts
  export default defineConfig({
    build: {
      minify: true,
      sourcemap: false, // Disable in production
    },
    server: {
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
        }
      }
    }
  })
  ```

**7. Using Components with Known Vulnerabilities**

- **Dependency Security**:
  - Keep all npm packages up to date
  - Run `npm audit` regularly
  - Use Dependabot for automated updates
  - Review security advisories for Vue.js ecosystem
  - Avoid packages with known vulnerabilities
  - Commands:
    ```bash
    npm audit
    npm audit fix
    npm outdated
    ```

**8. Cross-Site Request Forgery (CSRF)**

- **CSRF Protection**:
  - Use CSRF tokens for state-changing requests
  - Verify origin/referer headers in backend
  - Use SameSite cookie attribute
  - For JWT auth, implement double-submit pattern
  - Example:
    ```typescript
    // Include CSRF token in requests
    axios.defaults.headers.common['X-CSRF-Token'] = csrfToken
    ```

**9. Insecure Client-Side Storage**

- **Storage Security**:
  - Don't store sensitive data in localStorage/sessionStorage
  - Use HttpOnly cookies for tokens (if possible)
  - Clear storage on logout
  - Encrypt sensitive data before storing (if necessary)
  - Example:
    ```typescript
    // GOOD - Only store non-sensitive preferences
    localStorage.setItem('theme', 'dark')
    localStorage.setItem('language', 'en-US')

    // BAD - Don't store sensitive data
    localStorage.setItem('password', password) // NEVER DO THIS
    ```

**10. Insufficient Input Validation**

- **Form Validation**:
  - Validate all user input on frontend
  - Use schema validation libraries (Yup, Zod)
  - Implement real-time validation feedback
  - Sanitize input before displaying
  - Example:
    ```typescript
    import * as yup from 'yup'

    const schema = yup.object({
      email: yup.string().email().required(),
      password: yup.string().min(8).required(),
    })

    try {
      await schema.validate(formData)
    } catch (error) {
      // Show validation errors
    }
    ```

#### Additional Frontend Security Practices

**Content Security Policy (CSP)**:
- Implement strict CSP headers
- Whitelist trusted sources
- Disable inline scripts (use nonces or hashes if necessary)
- Report CSP violations
- Example CSP header:
  ```
  Content-Security-Policy:
    default-src 'self';
    script-src 'self' 'nonce-random123';
    style-src 'self' 'unsafe-inline';
    img-src 'self' data: https:;
    font-src 'self' data:;
    connect-src 'self' https://api.yourdomain.com;
  ```

**HTTPS Enforcement**:
- Always use HTTPS in production
- Implement HSTS headers
- Redirect HTTP to HTTPS
- Use secure cookies (Secure flag)

**Third-Party Dependencies**:
- Audit third-party libraries before use
- Use SRI for CDN resources
- Minimize third-party dependencies
- Review library permissions

**Error Handling**:
- Don't expose sensitive information in error messages
- Log errors to backend (not console)
- Show user-friendly error messages
- Implement global error handler

**API Communication**:
- Use Axios interceptors for global error handling
- Implement request/response validation
- Set timeout for all requests
- Handle network errors gracefully
- Example:
  ```typescript
  axios.interceptors.response.use(
    response => response,
    error => {
      if (error.response?.status === 401) {
        // Redirect to login
        router.push('/login')
      }
      return Promise.reject(error)
    }
  )
  ```

**Component Security**:
- Sanitize props if they come from user input
- Validate event payloads
- Don't pass raw user input to dangerous components
- Use TypeScript for type safety

**Vue.js Specific Security**:
- Don't use user input in `v-bind:is` (dynamic components)
- Be careful with `v-html` directive
- Validate props with TypeScript interfaces
- Use computed properties for derived data (not raw user input)

### Security Testing & Verification

**Backend Testing**:
- Test authentication and authorization logic
- Test input validation with malicious payloads
- Test rate limiting effectiveness
- Test SQL injection prevention (parameterized queries)
- Test error handling (no sensitive data leakage)
- Run security scanners (gitleaks, semgrep, govulncheck)

**Frontend Testing**:
- Test XSS prevention (try injecting scripts)
- Test CSRF protection
- Test route guards and permissions
- Test secure storage practices
- Audit dependencies with `npm audit`
- Test CSP effectiveness

**Integration Testing**:
- Test authentication flows end-to-end
- Test authorization at each layer
- Test secure API communication
- Test error handling across stack

**Security Checklist**:
- [ ] All SQL queries use parameterized statements
- [ ] All user input is validated and sanitized
- [ ] Passwords are hashed with bcrypt
- [ ] JWT tokens are validated and have expiration
- [ ] Sensitive data is never logged
- [ ] HTTPS is enforced in production
- [ ] Security headers are configured
- [ ] Rate limiting is implemented
- [ ] CORS is properly configured
- [ ] Dependencies are up to date and scanned
- [ ] XSS prevention is implemented (no v-html with user input)
- [ ] CSRF protection is in place
- [ ] Route guards protect sensitive pages
- [ ] Error messages don't expose sensitive information
- [ ] Security tests are passing

## Important Constraints

### Technical Constraints
- **Go Version**: Must use Go 1.25.5 for optimal performance
- **Node Version**: ^20.19.0 or >=22.12.0
- **Latency SLA**: <50ms transaction processing time
- **Timeout Configuration**: Configurable timeouts for transaction processing (300ms default) and rule evaluation (300ms default)
- **Timeout Environment Variables**: All timeout values MUST be configurable via environment variables (no hardcoded timeout values)

### Security Constraints
- All protected endpoints require JWT authentication
- Passwords must be hashed with bcrypt
- No sensitive data in logs or error messages
- Environment variables for secrets (never commit `.env`)
- RBAC enforcement on all administrative operations
- File exclusions: `.gitleaksignore`, `.semgrepignore` for security scanning

### Performance Constraints
- Connection pooling required for all database/Redis connections
- Rules must be cached in Redis to minimize database queries
- Worker concurrency configurable (default: 10 workers)
- Batch processing support (default: 50 transactions per batch)

## External Dependencies

### Required Services
- **PostgreSQL**: Primary database for transactions, rules, event schemas, users, roles, groups
  - Connection pooling via pgx
  - Migration scripts in `scripts/migrations/`
- **Redis**: Message queue and caching layer
  - Pub/sub for transaction queue and schema invalidation
  - Rules and schema caching with TTL
- **Docker**: Required for local development and deployment

### Optional Services
- **Observability**: OpenTelemetry support for metrics and tracing
- **TLS**: Optional TLS configuration for API (disabled by default)

### Documentation System
- **OpenSpec**: Spec-driven development framework for change proposals and capability documentation
  - Located in `openspec/` directory
  - Supports change proposals, design documents, and capability specs
  - See `openspec/AGENTS.md` for AI agent instructions

## Project Evolution & Decisions

This section documents key architectural decisions and evolution history to provide context for future development and AI agents.

### Dependency Updates (January 2026)

**Decision**: Updated all major dependencies to their latest stable versions

**Key Changes**:
- Updated Vue.js from 3.5.25 to 3.5.26
- Updated vue-i18n from 11.0.0 to 11.2.8 and migrated from legacy mode to Composition API mode
- Updated Vuetify from 3.7.0 to 3.11.6 (major version bump with new features)
- Updated Font Awesome from 6.5.1 to 7.1.0 (major version bump)
- Updated Vite from 7.2.4 to 7.3.1
- Updated TypeScript from 5.9.0 to 5.9.3
- Updated vue-router from 4.6.3 to 4.6.4

**Rationale**:
- Keep dependencies up to date for security and performance improvements
- Remove deprecated legacy mode from vue-i18n before v12 release
- Take advantage of new features and bug fixes in latest versions
- Font Awesome 7 brings new icons and better performance
- Vuetify 3.11 brings significant UI improvements and bug fixes

**Impact**:
- All 408 tests passing after migration
- No breaking changes to existing functionality
- Better TypeScript support with Composition API mode
- Improved build performance with latest Vite
- See `docs/dependency-updates-2026-01.md` for complete details

### Frontend Framework Migration

**Decision**: Migrated from React to Vue.js 3

**Rationale**:
- React codebase was confusing and difficult to maintain
- Vue.js provides a good balance between ease of use, simplicity, and is widely adopted in open-source projects
- Better developer experience for this project's needs

**Status**: Migration completed, but some layout bugs remain unresolved from the migration process

### Architecture: Vertical Slice

**Decision**: Adopted vertical slice architecture pattern

**Rationale**:
- Initial code was AI-generated and suboptimal
- Vertical slice significantly improved code readability for humans
- Better organization by feature rather than by technical layer

**Impact**: Improved maintainability and code organization

### Dependency Injection

**Decision**: Implemented dependency injection throughout the codebase

**Rationale**:
- Initial AI-generated code did not follow best practices
- Refactored to adhere to SOLID principles
- Improved testability and modularity

**Status**: Applied to handlers and services

### Code Quality Improvements

**Race Conditions**:
- Fixed race conditions in rule service detected via race condition flags in tests
- Solution was suggested by an AI agent and successfully applied

**ISP Violations**:
- Fixed Interface Segregation Principle violations in worker service
- Interfaces were too large and needed refactoring
- Solution suggested by AI agent

**Validation**:
- Added validation in handlers as a mandatory practice
- System not yet in production, but validation is required for any system

### Infrastructure Simplification

**Proxy Removal**:
- Removed front-end proxy that was added by an AI agent without necessity
- The proxy complicated the codebase without providing value
- Simplified deployment and architecture

### CI/CD Evolution

**Status**: Stable after learning curve

**History**:
- Initial setup had multiple iterations and fixes
- Process improved through learning GitHub Actions over time
- Current pipeline is stable and functional

### Observability

**OpenTelemetry**:
- Metrics infrastructure defined using OpenTelemetry
- Not yet fully evolved or implemented
- Plan to use OpenTelemetry for observability, but implementation pending

**Performance Monitoring**:
- Need help with performance optimization and monitoring
- Target latency: <50ms per transaction
- OpenTelemetry should help identify bottlenecks once fully implemented

### Security & Configuration

**TLS**:
- TLS configuration prepared for production
- Not yet configured in production
- Plan to make HTTPS mandatory in the future

**Environment Variables**:
- All required variables MUST be documented in `.env.example`
- If any variable is missing from `.env.example`, it must be corrected
- This is a mandatory requirement for the project
- All timeout values MUST be configurable via environment variables (no hardcoded timeout values)

### Database Migrations

**Current State**: Unsatisfactory, needs improvement

**Issues**:
- Current migration process via scripts is not ideal
- Process needs to be simplified and improved with a library
- Migration workflow needs refactoring

**Action Required**: Evaluate and implement a proper migration library/tool

### Testing Strategy

**Current Approach**: Spec-driven testing with comprehensive documentation

**Documentation**:
- **Unit Tests**: `docs/agents/unit-test.md` - Complete guide for Go and Vue/TypeScript
- **Integration Tests**: `docs/agents/integration-test.md` - Guide for testing with real dependencies

**Requirements**:
- Minimum 80% coverage for new code
- All tests must be deterministic (pass with `-count=50`)
- Race condition detection via `-race` flag
- Integration tests use `testcontainers-go` for real databases
- Follow AAA pattern (Arrange-Act-Assert) without comments

**Note**: AI agents MUST follow the testing documentation strictly when generating tests

### Regression Prevention Rules

**CRITICAL**: These rules MUST be followed when adding new features to prevent regressions in existing functionality.

#### Pre-Implementation Checklist

Before implementing any new feature, verify:

1. **OpenSpec Compliance**:
   - [ ] Change proposal created in `openspec/changes/` if required (new capability, breaking change, architecture shift)
   - [ ] Proposal reviewed and approved before implementation
   - [ ] Spec deltas created with proper ADDED/MODIFIED/REMOVED sections
   - [ ] All requirements have at least one scenario
   - [ ] `openspec validate <change-id> --strict` passes

2. **Impact Analysis**:
   - [ ] All affected modules identified
   - [ ] Breaking changes documented and marked with **BREAKING**
   - [ ] Migration plan created for breaking changes
   - [ ] Backward compatibility verified or migration path provided

3. **Test Coverage Requirements**:
   - [ ] Unit tests for all new business logic (minimum 80% coverage)
   - [ ] Integration tests for database/Redis interactions
   - [ ] API endpoint tests for new/modified endpoints
   - [ ] Frontend component tests for new UI components
   - [ ] E2E tests for critical user flows (if applicable)
   - [ ] Tests run with `-race` flag (Go) and pass deterministically (`-count=50`)

#### Mandatory Test Categories

Every new feature MUST include tests in these categories:

1. **Unit Tests** (Go):
   - [ ] All new functions/methods have unit tests
   - [ ] Edge cases and error conditions covered
   - [ ] Mock dependencies properly (use interfaces)
   - [ ] Tests are deterministic (no flaky tests)
   - [ ] Race conditions tested with `-race` flag

2. **Integration Tests** (Go):
   - [ ] Database operations tested with real PostgreSQL (testcontainers-go)
   - [ ] Redis operations tested with real Redis (testcontainers-go)
   - [ ] Transaction processing flow tested end-to-end
   - [ ] Rules engine evaluation tested with real schemas
   - [ ] Velocity functions tested with transaction history

3. **API Tests** (Go):
   - [ ] All new endpoints have handler tests
   - [ ] Authentication/authorization tested
   - [ ] Input validation tested
   - [ ] Error responses tested
   - [ ] Status codes verified

4. **Frontend Tests** (Vue/TypeScript):
   - [ ] Component unit tests for new components
   - [ ] View tests for new pages
   - [ ] Store tests for new Pinia stores
   - [ ] Router guard tests for protected routes
   - [ ] i18n tests (verify translation keys exist)

5. **Regression Tests**:
   - [ ] Existing tests still pass (no broken tests)
   - [ ] Existing API endpoints still work
   - [ ] Existing UI components still render correctly
   - [ ] Existing rules still evaluate correctly
   - [ ] Existing schemas still parse correctly

#### Breaking Changes Protocol

If a feature introduces breaking changes:

1. **Documentation**:
   - [ ] Breaking change marked with **BREAKING** in proposal
   - [ ] Migration guide created
   - [ ] Deprecation notice added (if applicable)
   - [ ] Version bump planned (if using semantic versioning)

2. **Backward Compatibility**:
   - [ ] Old API versions supported (if versioning exists)
   - [ ] Data migration scripts created
   - [ ] Feature flags considered for gradual rollout
   - [ ] Rollback plan documented

3. **Communication**:
   - [ ] Breaking changes communicated in proposal
   - [ ] Impact assessment shared with stakeholders
   - [ ] Timeline for migration provided

#### Performance Regression Prevention

1. **Latency Requirements**:
   - [ ] New features don't increase transaction processing time beyond <50ms target
   - [ ] Database queries optimized (indexes added if needed)
   - [ ] Redis caching used appropriately
   - [ ] No N+1 query problems introduced

2. **Load Testing**:
   - [ ] High-volume scenarios tested (if applicable)
   - [ ] Worker concurrency tested
   - [ ] Batch processing performance verified
   - [ ] Memory leaks checked (long-running tests)

3. **Resource Usage**:
   - [ ] Connection pool limits respected
   - [ ] Redis memory usage monitored
   - [ ] Database connection limits not exceeded
   - [ ] No unbounded growth in data structures

#### Database Schema Changes

When modifying database schema:

1. **Migration Safety**:
   - [ ] Migration scripts tested on sample data
   - [ ] Rollback scripts created and tested
   - [ ] Data loss prevention verified
   - [ ] Index creation/dropping optimized (concurrent operations)

2. **Compatibility**:
   - [ ] Old code works with new schema (if backward compatible)
   - [ ] New code handles missing columns gracefully (if adding columns)
   - [ ] Foreign key constraints validated
   - [ ] Unique constraints preserved

3. **Performance Impact**:
   - [ ] New indexes don't slow down writes significantly
   - [ ] Query plans analyzed for new queries
   - [ ] Table statistics updated after migrations

#### API Contract Stability

When modifying API endpoints:

1. **Request/Response Validation**:
   - [ ] Request validation doesn't break existing clients
   - [ ] Response structure changes documented
   - [ ] Optional fields added (not required fields removed)
   - [ ] Default values preserve existing behavior

2. **Error Handling**:
   - [ ] Error codes don't change for existing errors
   - [ ] New error codes documented
   - [ ] Error messages remain descriptive
   - [ ] HTTP status codes follow REST conventions

3. **Versioning** (if applicable):
   - [ ] API versioning strategy followed
   - [ ] Old versions deprecated with timeline
   - [ ] Version negotiation handled correctly

#### Frontend Compatibility

When modifying frontend:

1. **Browser Compatibility**:
   - [ ] Tested in supported browsers
   - [ ] No breaking changes to public component APIs
   - [ ] CSS changes don't break existing layouts
   - [ ] JavaScript errors don't break existing functionality

2. **State Management**:
   - [ ] Pinia store changes don't break existing components
   - [ ] Router changes don't break existing navigation
   - [ ] i18n changes don't break existing translations

3. **Component APIs**:
   - [ ] Props changes are backward compatible
   - [ ] Events emitted remain consistent
   - [ ] Slots remain functional

#### Rules Engine Stability

When modifying rules engine:

1. **Expression Compatibility**:
   - [ ] Existing expressions still evaluate correctly
   - [ ] New helper functions don't break existing rules
   - [ ] Schema changes don't invalidate existing rules
   - [ ] Expression cache invalidation works correctly

2. **Performance**:
   - [ ] Rule evaluation time doesn't increase significantly
   - [ ] Velocity functions still perform within limits
   - [ ] Schema loading doesn't slow down processing

3. **Data Integrity**:
   - [ ] Transaction history queries still work
   - [ ] Metadata structure preserved
   - [ ] Field path extraction remains accurate

#### Pre-Merge Validation Checklist

Before merging any PR with new features:

1. **Code Quality**:
   - [ ] All linters pass (`make lint`)
   - [ ] All tests pass (`make test`)
   - [ ] No race conditions detected (`go test -race`)
   - [ ] No flaky tests (`go test -count=50`)
   - [ ] Code coverage meets 80% minimum
   - [ ] No security vulnerabilities (gitleaks, semgrep)

2. **Documentation**:
   - [ ] README.md updated (if user-facing changes)
   - [ ] API documentation updated (if API changes)
   - [ ] `.env.example` updated (if new env vars)
   - [ ] OpenSpec proposal archived (if change completed)
   - [ ] Code comments added for complex logic

3. **Integration**:
   - [ ] Integration tests pass with real databases
   - [ ] Docker builds succeed
   - [ ] CI/CD pipeline passes
   - [ ] No breaking changes to existing workflows

4. **Manual Testing**:
   - [ ] Feature works as specified
   - [ ] Existing features still work
   - [ ] UI renders correctly in both languages (i18n)
   - [ ] Error handling works correctly
   - [ ] Performance acceptable

#### Post-Implementation Monitoring

After deploying new features:

1. **Observability**:
   - [ ] Metrics collected for new features
   - [ ] Error rates monitored
   - [ ] Performance metrics tracked
   - [ ] Alerts configured (if applicable)

2. **Rollback Plan**:
   - [ ] Rollback procedure documented
   - [ ] Database rollback scripts ready
   - [ ] Feature flags available (if used)
   - [ ] Communication plan for rollback

3. **Validation**:
   - [ ] Production behavior matches test expectations
   - [ ] No unexpected errors in logs
   - [ ] Performance within acceptable limits
   - [ ] User feedback collected

#### AI Agent Requirements

When AI agents implement new features, they MUST:

1. **Before Implementation**:
   - Create OpenSpec proposal if required
   - Perform impact analysis
   - Identify all affected modules
   - Plan test coverage

2. **During Implementation**:
   - Write tests alongside code (TDD approach preferred)
   - Run tests frequently
   - Verify no existing tests break
   - Check for race conditions

3. **After Implementation**:
   - Run full test suite
   - Verify code coverage
   - Check linter errors
   - Update documentation
   - Verify all checklist items completed

4. **Before Completion**:
   - Confirm all regression prevention rules followed
   - Verify no breaking changes introduced (or properly documented)
   - Ensure backward compatibility maintained
   - Validate performance requirements met

#### Enforcement

These rules are **MANDATORY** and must be enforced:

- **CI/CD Pipeline**: Automated checks for test coverage, linters, and security scans
- **Code Review**: Reviewers must verify checklist completion
- **Pre-merge Hooks**: Git hooks should run basic validation
- **Documentation**: All new features must include appropriate documentation

**Violation Consequences**:
- PRs missing required tests will be rejected
- Breaking changes without migration plans will be rejected
- Performance regressions will block deployment
- Missing documentation will delay merge

### Known Issues

**Layout Bugs**:
- Some layout bugs remain from Vue.js migration
- Root cause not yet identified
- Bugs introduced during code generation by AI agent
- Needs investigation and resolution

### Development Guidelines for AI Agents

**Important Notes**:
1. **Code Generation**: Initial code was AI-generated and required significant refactoring. When generating code:
   - Follow SOLID principles strictly
   - Use vertical slice architecture
   - Implement proper dependency injection
   - Ensure interfaces follow ISP (not too large)
   - Avoid unnecessary complexity (like the proxy that was removed)
   - Use fundamental abstractions to simplify complex logic
   - Refactor magic numbers to named constants

2. **Validation**: Always add validation to handlers - it's mandatory

3. **Environment Variables**: Always update `.env.example` when adding new environment variables. Improve environment variable management by grouping related variables and providing clear documentation.

4. **Testing**: 
   - **Unit Tests**: Follow `docs/agents/unit-test.md` strictly
   - **Integration Tests**: Follow `docs/agents/integration-test.md` strictly
   - Use race condition flags (`-race`) when running tests
   - Minimum coverage: 80% for new code
   - Run `go test -count=50` to detect flaky tests
   - Use `testcontainers-go` for integration tests with databases

5. **Performance**: Keep <50ms latency target in mind. Need help with performance optimization and monitoring

6. **Migrations**: Current migration process needs improvement - consider suggesting migration libraries when working with database changes

7. **Avoid Unnecessary Complexity**: Don't add features (like proxies) without clear necessity

8. **Frontend Component Standards**:
   - Standardize base components for consistency (BaseButton, BaseInput, BaseTable, BaseBadge, etc.)
   - Use standardized validation fields across forms
   - Remove unused code and components regularly
   - Padronize (standardize) components in front-end for maintainability

9. **Build Optimization**:
   - Use build chunks for better code splitting and faster builds
   - Optimize build time by building only changed services when possible
   - Use Docker BuildKit for faster builds with better caching

10. **Code Quality**:
    - Remove unused files and code regularly
    - Fix concurrent issues and race conditions immediately
    - Fix layout bugs introduced during migrations
    - Keep codebase clean and maintainable

11. **Documentation**:
    - Use OpenSpec for spec-driven development
    - Update README.md when adding significant features
    - Document architectural decisions in project.md

12. **Internationalization (i18n)**:
    - **NEVER hardcode user-facing strings** - always use translation keys
    - Add translations to BOTH `pt-BR.json` and `en-US.json` simultaneously
    - Use hierarchical key naming: `views.dashboard.title`, `errors.notFound`
    - Test UI with both languages to ensure text fits properly
    - Refer to `src/ui/src/locales/README.md` for detailed guidelines

13. **Mandatory Code Review and Testing Verification**:
    - **REQUIRED**: Whenever the user requests a correction or analysis, AI agents MUST perform a comprehensive scan to verify:
      - The project adheres to all guidelines defined in this document
      - All related tests have been generated or updated appropriately
      - Code changes follow the established patterns and conventions
      - Test coverage meets the minimum 80% requirement for new/modified code
      - Tests follow the guidelines in `docs/agents/unit-test.md` and `docs/agents/integration-test.md`
    - This verification MUST be performed proactively, not just when explicitly requested
    - If any guideline violations or missing tests are found, they MUST be addressed before considering the task complete

14. **Regression Prevention**:
    - **CRITICAL**: When adding new features, AI agents MUST follow the "Regression Prevention Rules" section in this document
    - Complete all pre-implementation, mandatory test categories, and pre-merge checklists
    - Verify no breaking changes introduced (or properly documented with migration plans)
    - Ensure backward compatibility maintained
    - Validate performance requirements met (<50ms latency target)
    - See "Regression Prevention Rules" section for complete checklist and requirements

15. **CSRF Protection**:
    - **CRITICAL**: The system uses CSRF tokens for all state-changing requests (POST, PUT, PATCH, DELETE)
    - **Backend Implementation**:
      - CSRF tokens generated on login/register (`src/api/internal/auth/handler.go`)
      - Tokens stored in Redis with 24-hour TTL (`src/pkg/csrf/csrf.go`)
      - Middleware validates `X-CSRF-Token` header on all non-GET requests (`src/api/internal/shared/middleware/csrf.go`)
      - Excluded paths: `/api/v1/auth/login`, `/api/v1/auth/register`, `/api/v1/branding` (public GET)
    - **Frontend Implementation**:
      - Auth store stores both JWT and CSRF tokens in memory (`src/ui/src/stores/auth.ts`)
      - API client automatically adds `X-CSRF-Token` header to POST/PUT/PATCH/DELETE requests (`src/ui/src/lib/api.ts`)
      - Login/Register views extract `csrf_token` from response (`src/ui/src/views/LoginView.vue`)
      - CSRF token cleared on logout
    - **Common Issues and Solutions**:
      - **403 "CSRF token missing"**: User must logout and login again after CSRF implementation
      - **All views use API client**: Never use `fetch()` or `axios` directly - always use `api.post()`, `api.put()`, etc.
      - **Stores use API client**: Pinia stores (like `branding.ts`, `systemMode.ts`) must use `api` from `@/lib/api`
      - **Token in memory only**: CSRF token is NOT in localStorage - stored in memory for security
    - **Troubleshooting CSRF Issues**:
      1. Check if user logged in after CSRF implementation (token issued at login)
      2. Verify `X-CSRF-Token` header in browser DevTools Network tab
      3. Check backend logs for "CSRF token missing" or "Invalid CSRF token"
      4. Ensure all state-changing operations use `api.post/put/patch/delete`
      5. Verify `csrfTokenGetterHolder.getter` is registered in `api.ts`
    - **When Adding New Features**:
      - Use `api.post()`, `api.put()`, `api.patch()`, `api.delete()` for all HTTP requests
      - Never use `fetch()` directly - it won't include CSRF token
      - Test with browser DevTools to verify `X-CSRF-Token` header is sent
      - Verify backend accepts the request (200/201 status, not 403)
    - **Testing CSRF Protection**:
      - Manual test: Open DevTools → Network → Try POST/PUT operation → Check Request Headers
      - Should see: `X-CSRF-Token: <base64-token>`
      - Backend should return 200/201, not 403
    - **CSRF Token Lifecycle**:
      1. User logs in → Backend generates CSRF token → Returns in response
      2. Frontend stores token in memory (auth store)
      3. Frontend registers token getter with API client
      4. API client adds token to headers on POST/PUT/PATCH/DELETE
      5. Backend validates token on each request
      6. User logs out → Token cleared from memory and Redis
