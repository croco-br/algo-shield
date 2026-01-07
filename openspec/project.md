# Project Context

## Purpose
AlgoShield is an open-source, high-performance fraud detection and anti-money laundering (AML) transaction analysis tool designed to process transactions with ultra-low latency (<50ms). The system provides:
- Real-time fraud prevention (pre-transaction analysis)
- Post-transaction AML compliance monitoring
- Custom rules engine with hot-reload capabilities
- Risk scoring and transaction classification
- User and permission management with RBAC

## Tech Stack

### Backend (Go 1.25.4)
- **Web Framework**: Fiber v2 (high-performance HTTP server)
- **Database Driver**: pgx v5 (PostgreSQL connection pooling)
- **Cache/Queue**: go-redis v9 (Redis client)
- **Authentication**: golang-jwt v5 + bcrypt (JWT tokens with password hashing)
- **Validation**: go-playground/validator v10
- **Observability**: OpenTelemetry (metrics and tracing)
- **Concurrency**: golang.org/x/sync (worker pools)
- **Expression Engine**: expr-lang/expr v1.17.7 (custom rule expressions)

### Frontend (Vue.js 3)
- **Framework**: Vue 3.5+ with Composition API
- **Language**: TypeScript 5.9
- **State Management**: Pinia 3.0
- **Routing**: Vue Router 4.6
- **UI Components**: Vuetify 3.7 (Material Design)
- **Icons**: Font Awesome 6.5 (free solid icons)
- **Syntax Highlighting**: Prism.js 1.30
- **Styling**: Tailwind CSS 4.1 (with PostCSS)
- **Build Tool**: Vite 7.2 (with chunk optimization)
- **Dev Tools**: Vue DevTools, vue-tsc

### Infrastructure
- **Database**: PostgreSQL (primary data store)
- **Cache/Queue**: Redis (message queue + rules caching)
- **Containerization**: Docker + Docker Compose
- **Deployment**: Multi-stage Docker builds

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

#### Go Testing
- Use standard `go test` with race detector (`-race` flag)
- Run tests in parallel (`-parallel 4`)
- Unit tests for business logic
- Integration tests for database interactions
- Benchmarks for performance-critical code (rules engine)
- Test file naming: `*_test.go`
- Coverage target: Not explicitly defined, but comprehensive coverage expected

#### Commands
- `make test`: Run all tests with race detector
- `make bench`: Run rules engine benchmarks
- `make lint`: Run golangci-lint

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

### Fraud Detection & AML
- **Transaction Analysis**: Real-time evaluation of financial transactions
- **Risk Scoring**: Cumulative scores based on matched rules with risk levels (Low: 0-49, Medium: 50-79, High: 80-100)
- **Rule System**: Custom expression-based rules using [expr-lang](https://github.com/expr-lang/expr)
  - **Custom Expressions**: All rules use `custom` type with `custom_expression` condition
  - **Schema-Based**: Rules are associated with event schemas for type safety and field validation
  - **Helper Functions**: Built-in functions for velocity checks (`velocityCount`, `velocitySum`) and geography (`pointInPolygon`)
  - **Expression Syntax**: Supports comparisons, logical operators, array operations, and nested field access
- **Event Schemas**: Define transaction event structure with automatic field extraction from sample JSON
  - Schemas enable automatic field discovery and validation
  - Rules must reference a schema for proper evaluation
  - Schema changes trigger cache invalidation via Redis pub/sub
- **Rule Actions**: `allow`, `block`, `review` (score action removed, scoring is implicit)
- **Processing Modes**: Pre-transaction (fraud prevention) and post-transaction (AML compliance)

### Performance Requirements
- **Target Latency**: <50ms per transaction
- **Throughput**: High-volume transaction processing
- **Scalability**: Horizontal scaling of workers
- **Reliability**: Retry mechanisms with exponential backoff

### User Roles
- **admin**: Full system access (user management, all rules)
- **rule_editor**: Can create/update/delete rules
- **viewer**: Read-only access

## Important Constraints

### Technical Constraints
- **Go Version**: Must use Go 1.25.4 for optimal performance
- **Node Version**: ^20.19.0 or >=22.12.0
- **Latency SLA**: <50ms transaction processing time
- **Timeout Configuration**: Configurable timeouts for transaction processing (300ms default) and rule evaluation (300ms default)

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

### Database Migrations

**Current State**: Unsatisfactory, needs improvement

**Issues**:
- Current migration process via scripts is not ideal
- Process needs to be simplified and improved with a library
- Migration workflow needs refactoring

**Action Required**: Evaluate and implement a proper migration library/tool

### Testing Strategy

**Current Approach**: Waiting for scope stabilization

**Rationale**:
- System is still evolving and scope is not fully defined
- Plan to add tests via AI once the scope is more stable
- Tests will be added after core functionality is stabilized

**Note**: Race condition detection is already in place via test flags (`-race`)

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

4. **Testing**: Use race condition flags (`-race`) when running tests to detect concurrency issues

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
