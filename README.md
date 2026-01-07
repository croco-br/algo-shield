# 🛡️ AlgoShield

**AlgoShield** is an open-source, high-performance fraud detection and anti-money laundering (AML) transaction analysis tool designed to process transactions with ultra-low latency (<50ms).

## 🎯 Key Features

- **⚡ Ultra-Fast Processing**: Process transactions in <50ms with highly optimized Go workers
- **🔧 Custom Expression Rules Engine**: Flexible rule system using [expr-lang](https://github.com/expr-lang/expr) for powerful fraud detection expressions
- **🔄 Hot-Reload Rules & Schemas**: Update rules and event schemas in real-time without restarting services
- **📋 Event Schema Management**: Define and manage event schemas with automatic field extraction from sample JSON
- **📊 Risk Scoring**: Flexible scoring system with rule-based risk accumulation
- **🎯 Dual Processing Modes**: Support for pre-transaction (fraud prevention) and post-transaction (AML) analysis
- **🚀 High Scalability**: Horizontally scalable worker architecture
- **📈 Real-time Analysis**: Process events through Redis queues with minimal latency
- **🔐 Authentication & Authorization**: JWT-based authentication with role-based access control (RBAC)
- **👥 User Management**: Complete user, role, and group management system
- **🛡️ Permission System**: Fine-grained permissions for rule editing and administrative tasks
- **🎨 Branding Configuration**: White-label customization with configurable colors, logos, and app name
- **📚 OpenSpec Documentation**: Spec-driven development with change proposals and capability documentation

## 🏗️ Architecture

AlgoShield is built with a modern microservices architecture:

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│     UI      │──────│     API     │──────│   Worker    │
│  (Vue.js)   │      │   (Fiber)   │      │   (Rules)   │
└─────────────┘      └─────────────┘      └─────────────┘
                            │                     │
                            └──────┬──────────────┘
                                   │
                            ┌──────┴──────┐
                            │  PostgreSQL │
                            │    Redis    │
                            └─────────────┘
```

### Components

- **API Service**: RESTful API built with Fiber v2 (Go 1.25.4) for high-performance HTTP handling with JWT authentication, validation, dependency injection, and OpenTelemetry observability support
- **Worker Service**: Transaction processing engine with custom expression-based rules evaluation using expr-lang, schema management, and hot-reload support
- **UI**: Vue.js 3.5+ (TypeScript 5.9) modern web interface with Vuetify 3.7 (Material Design) components, Pinia 3.0 state management, Tailwind CSS 4.1, standardized base components, Font Awesome 6.5 icons, and Prism.js syntax highlighting for rule management, schema management, and user administration
- **PostgreSQL**: Primary data store for transactions, rules, event schemas, users, roles, and groups (PostgreSQL 16)
- **Redis**: Message queue for async processing, rules caching, and schema invalidation pub/sub (Redis 7)

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.25.4 (for local development)
- Node.js ^20.19.0 or >=22.12.0 (for UI development)

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/yourusername/algo-shield.git
cd algo-shield
```

2. Start all services:
```bash
make up
```

Or using Docker Compose directly:
```bash
docker-compose up -d
```

3. Access the services:
- **UI**: http://localhost:3000
- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

**Note**: Default admin credentials are created via migration. Check `scripts/migrations/004_insert_admin.sql` for details.

### Local Development

1. Install Git hooks (recommended):
```bash
./scripts/install-hooks.sh
```

2. Install dependencies:
```bash
make install
```

3. Start PostgreSQL and Redis:
```bash
docker-compose up -d postgres redis
```

4. Run database migrations:
```bash
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=algoshield
export POSTGRES_PASSWORD=algoshield_secret
export POSTGRES_DB=algoshield

# Run all migrations in order
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/001_initial_schema.sql
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/002_auth_schema.sql
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/003_local_auth.sql
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/004_insert_admin.sql
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/005_branding_config.sql
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/006_add_header_color.sql
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/007_event_schemas.sql
psql -h localhost -U algoshield -d algoshield -f scripts/migrations/008_test_data.sql
```

**Note**: The migrations script (`migrations.sh`) is designed for Docker environments. For local development, run migrations manually as shown above. The project includes 8 migration files:
- `001_initial_schema.sql` - Initial database schema
- `002_auth_schema.sql` - Authentication tables
- `003_local_auth.sql` - Local authentication setup
- `004_insert_admin.sql` - Default admin user
- `005_branding_config.sql` - Branding configuration
- `006_add_header_color.sql` - Header color customization
- `007_event_schemas.sql` - Event schema management
- `008_test_data.sql` - Test data (optional)

5. Start the API:
```bash
cd src/api/cmd && go run main.go
```

6. Start the Worker (in another terminal):
```bash
cd src/workers/cmd && go run main.go
```

7. Start the UI (in another terminal):
```bash
cd src/ui && npm run dev
```

The UI will be available at http://localhost:5173 (Vite dev server).

## 📖 API Documentation

### Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

#### Register User

```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secure_password"
}
```

#### Login

```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "secure_password"
}
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "active": true,
    "roles": []
  }
}
```

#### Get Current User

```bash
GET /api/v1/auth/me
Authorization: Bearer <token>
```

#### Logout

```bash
POST /api/v1/auth/logout
Authorization: Bearer <token>
```

### User Management (Admin Only)

#### List Users

```bash
GET /api/v1/permissions/users
Authorization: Bearer <token>
```

#### Get User

```bash
GET /api/v1/permissions/users/{id}
Authorization: Bearer <token>
```

#### Update User Active Status

```bash
PUT /api/v1/permissions/users/{id}/active
Authorization: Bearer <token>
Content-Type: application/json

{
  "active": true
}
```

#### Assign Role to User

```bash
POST /api/v1/permissions/users/{userId}/roles
Authorization: Bearer <token>
Content-Type: application/json

{
  "role_id": "role_uuid"
}
```

#### Remove Role from User

```bash
DELETE /api/v1/permissions/users/{userId}/roles/{roleId}
Authorization: Bearer <token>
```

### Roles Management (Admin Only)

#### List Roles

```bash
GET /api/v1/roles
Authorization: Bearer <token>
```

#### Get Role

```bash
GET /api/v1/roles/{id}
Authorization: Bearer <token>
```

### Groups Management (Admin Only)

#### List Groups

```bash
GET /api/v1/groups
Authorization: Bearer <token>
```

#### Get Group

```bash
GET /api/v1/groups/{id}
Authorization: Bearer <token>
```

### Branding Configuration

#### Get Branding Configuration

Retrieve the current branding configuration. This is a public endpoint and does not require authentication.

```bash
GET /api/v1/branding
```

Response:
```json
{
  "id": 1,
  "app_name": "AlgoShield",
  "icon_url": "/assets/logo.svg",
  "favicon_url": "/favicon.ico",
  "primary_color": "#3B82F6",
  "secondary_color": "#10B981",
  "header_color": "#1e1e1e",
  "created_at": "2024-12-05T10:00:00Z",
  "updated_at": "2024-12-05T10:00:00Z"
}
```

**Note**: If no branding configuration exists, the API returns default values.

#### Update Branding Configuration

**Requires `admin` role**

Update the branding configuration for white label customization.

```bash
PUT /api/v1/branding
Authorization: Bearer <token>
Content-Type: application/json

{
  "app_name": "My Company",
  "icon_url": "/assets/custom-logo.png",
  "favicon_url": "/favicon-custom.ico",
  "primary_color": "#FF5733",
  "secondary_color": "#33FF57",
  "header_color": "#2C3E50"
}
```

Response:
```json
{
  "id": 1,
  "app_name": "My Company",
  "icon_url": "/assets/custom-logo.png",
  "favicon_url": "/favicon-custom.ico",
  "primary_color": "#FF5733",
  "secondary_color": "#33FF57",
  "header_color": "#2C3E50",
  "created_at": "2024-12-05T10:00:00Z",
  "updated_at": "2024-12-05T10:30:00Z"
}
```

**Validation Rules:**
- `app_name`: Required, 1-100 characters
- `icon_url`: Optional, must be a valid URI or file path
- `favicon_url`: Optional, must be a valid URI or file path
- `primary_color`: Required, must be in hex format (#RGB or #RRGGBB)
- `secondary_color`: Required, must be in hex format (#RGB or #RRGGBB)
- `header_color`: Required, must be in hex format (#RGB or #RRGGBB)

**Error Responses:**
- `400 Bad Request`: Invalid request body or validation errors
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: User does not have admin role
- `500 Internal Server Error`: Server error during update

### Transactions

#### Process Transaction

Submit a transaction for analysis:

```bash
POST /api/v1/transactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "external_id": "txn_123456",
  "amount": 5000.00,
  "currency": "USD",
  "origin": "ACC001",
  "destination": "ACC002",
  "type": "transfer",
  "metadata": {
    "ip_address": "192.168.1.1",
    "device_id": "device_123"
  },
  "timestamp": "2024-12-05T10:00:00Z"
}
```

Response:
```json
{
  "status": "queued",
  "external_id": "txn_123456",
  "processing_time": 5,
  "message": "Transaction queued for processing"
}
```

### Get Transaction

Retrieve transaction details:

```bash
GET /api/v1/transactions/{id}
Authorization: Bearer <token>
```

### List Transactions

```bash
GET /api/v1/transactions?limit=50&offset=0
Authorization: Bearer <token>
```

### Create Rule

**Requires `admin` or `rule_editor` role**

All rules use custom expressions. Rules must be associated with an event schema.

```bash
POST /api/v1/rules
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "High Value Transaction",
  "description": "Flag transactions over $10,000",
  "action": "review",
  "priority": 500,
  "enabled": true,
  "schema_id": "uuid-of-event-schema",
  "conditions": {
    "custom_expression": "amount > 10000"
  },
  "score": 50
}
```

### Update Rule

**Requires `admin` or `rule_editor` role**

```bash
PUT /api/v1/rules/{id}
Authorization: Bearer <token>
Content-Type: application/json
```

### Delete Rule

**Requires `admin` or `rule_editor` role**

```bash
DELETE /api/v1/rules/{id}
Authorization: Bearer <token>
```

### List Rules

```bash
GET /api/v1/rules
Authorization: Bearer <token>
```

### Get Rule

```bash
GET /api/v1/rules/{id}
Authorization: Bearer <token>
```

**Note**: Rule creation, update, and deletion require `admin` or `rule_editor` role.

### Event Schemas

Event schemas define the structure of transaction events and enable automatic field extraction from sample JSON. Rules can be associated with specific schemas to ensure type safety and proper field validation.

#### List Schemas

```bash
GET /api/v1/schemas
Authorization: Bearer <token>
```

#### Get Schema

```bash
GET /api/v1/schemas/{id}
Authorization: Bearer <token>
```

#### Create Schema

**Requires `admin` or `rule_editor` role**

Create a new event schema from sample JSON. The system automatically extracts all fields from the sample JSON.

```bash
POST /api/v1/schemas
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Payment Transaction",
  "description": "Schema for payment transactions",
  "sample_json": {
    "amount": 100.50,
    "currency": "USD",
    "origin": "ACC001",
    "destination": "ACC002",
    "timestamp": "2024-12-05T10:00:00Z",
    "metadata": {
      "ip_address": "192.168.1.1",
      "device_id": "device_123"
    }
  }
}
```

Response:
```json
{
  "id": "uuid",
  "name": "Payment Transaction",
  "description": "Schema for payment transactions",
  "sample_json": { ... },
  "extracted_fields": [
    "amount",
    "currency",
    "origin",
    "destination",
    "timestamp",
    "metadata.ip_address",
    "metadata.device_id"
  ],
  "created_at": "2024-12-05T10:00:00Z",
  "updated_at": "2024-12-05T10:00:00Z"
}
```

#### Update Schema

**Requires `admin` or `rule_editor` role**

```bash
PUT /api/v1/schemas/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Updated Schema Name",
  "description": "Updated description",
  "sample_json": { ... }
}
```

#### Delete Schema

**Requires `admin` or `rule_editor` role**

**Note**: Schemas cannot be deleted if they are referenced by any rules.

```bash
DELETE /api/v1/schemas/{id}
Authorization: Bearer <token>
```

#### Parse Schema

**Requires `admin` or `rule_editor` role**

Re-extract fields from a schema's sample JSON without updating other fields.

```bash
POST /api/v1/schemas/{id}/parse
Authorization: Bearer <token>
```

## 🔧 Custom Expression Rules

AlgoShield uses a flexible custom expression-based rule system. All rules evaluate expressions using [expr-lang](https://github.com/expr-lang/expr) against event schemas.

### Basic Rule Structure

```json
{
  "name": "High Value Transaction",
  "description": "Flag transactions over $10,000",
  "action": "review",
  "priority": 500,
  "enabled": true,
  "schema_id": "uuid-of-event-schema",
  "conditions": {
    "custom_expression": "amount > 10000"
  }
}
```

### Expression Examples

**Amount Threshold:**
```javascript
amount > 10000
```

**Multiple Conditions:**
```javascript
amount > 5000 and currency != "USD"
```

**String Matching:**
```javascript
user.country in ["RU", "CN", "KP"]
```

**Boolean Checks:**
```javascript
metadata.is_suspicious == true
```

**Complex Logic:**
```javascript
(amount > 10000 and currency == "USD") or (amount > 5000 and user.country == "RU")
```

### Helper Functions

#### Polygon Checks

Check if a geographic point is inside a polygon:

```javascript
pointInPolygon(location.lat, location.lon, [[37.7749, -122.4194], [37.7849, -122.4094], [37.7649, -122.4294]])
```

The polygon is defined as a 2D array of `[latitude, longitude]` coordinate pairs. The function uses the ray casting algorithm to determine if a point is inside the polygon.

#### Velocity Checks

Check transaction velocity (count or sum) within a time window:

```javascript
// Count transactions in the last hour (3600 seconds)
velocityCount(origin, 3600) > 10

// Sum transaction amounts in the last hour
velocitySum(origin, 3600) > 10000
```

**Note:** Velocity checks query transaction history from the database. The `origin` field should match the account identifier in your event schema.

### Recreating Legacy Rule Types

The following examples show how to recreate common rule patterns using custom expressions:

#### Amount Rule (Threshold Check)
```javascript
// Old: type: "amount", conditions: { "amount_threshold": 10000 }
// New:
amount > 10000
```

#### Velocity Rule (Transaction Count)
```javascript
// Old: type: "velocity", conditions: { "transaction_count": 10, "time_window_seconds": 3600 }
// New:
velocityCount(origin, 3600) > 10
```

#### Velocity Rule (Amount Sum)
```javascript
// Old: type: "velocity", conditions: { "amount_threshold": 10000, "time_window_seconds": 3600 }
// New:
velocitySum(origin, 3600) > 10000
```

#### Blocklist Rule
```javascript
// Old: type: "blocklist", conditions: { "blocklisted_accounts": ["ACC123", "ACC456"] }
// New:
origin in ["ACC123", "ACC456"]
```

#### Geography Rule (Polygon Check)
```javascript
// Old: type: "geography", conditions: { "polygon": [[lat1, lon1], [lat2, lon2], ...] }
// New:
pointInPolygon(location.lat, location.lon, [[lat1, lon1], [lat2, lon2], [lat3, lon3]])
```

### Expression Syntax

Expressions support:
- **Comparisons**: `==`, `!=`, `>`, `<`, `>=`, `<=`
- **Logical Operators**: `and`, `or`, `not`
- **Array Operations**: `in`, `contains`
- **Nested Fields**: Use dot notation (e.g., `user.country`, `metadata.ip_address`)
- **Helper Functions**: `pointInPolygon()`, `velocityCount()`, `velocitySum()`

For complete expression syntax, see the [expr-lang documentation](https://github.com/expr-lang/expr).

## 📊 Rule Actions

- **allow**: Explicitly allow the transaction
- **block**: Block the transaction immediately (highest priority)
- **review**: Flag for manual review (medium priority)

## 🎯 Risk Levels

Transactions are automatically assigned risk levels based on cumulative scores:

- **Low**: Score 0-49
- **Medium**: Score 50-79
- **High**: Score 80-100

## ⚙️ Configuration

Configuration is managed through environment variables. See `.env.example` for a complete list.

### Database
- `POSTGRES_HOST`: PostgreSQL host (default: localhost)
- `POSTGRES_PORT`: PostgreSQL port (default: 5432)
- `POSTGRES_USER`: Database user (default: algoshield)
- `POSTGRES_PASSWORD`: Database password
- `POSTGRES_DB`: Database name (default: algoshield)

### Redis
- `REDIS_HOST`: Redis host (default: localhost)
- `REDIS_PORT`: Redis port (default: 6379)

### API
- `API_HOST`: API bind address (default: 0.0.0.0)
- `API_PORT`: API port (default: 8080)
- `TLS_ENABLE`: Enable TLS (default: false)
- `TLS_CERT_PATH`: Path to TLS certificate
- `TLS_KEY_PATH`: Path to TLS private key
- `JWT_SECRET`: Secret key for JWT token signing (required)
- `JWT_EXPIRATION_HOURS`: JWT token expiration in hours (default: 24)
- `ENVIRONMENT`: Environment name (development, staging, production)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

### Worker
- `WORKER_CONCURRENCY`: Number of concurrent workers (default: 10)
- `WORKER_BATCH_SIZE`: Batch processing size (default: 50)
- `WORKER_TIMEOUT_TRANSACTION_PROCESSING`: Timeout for transaction processing (default: 300ms)
- `WORKER_TIMEOUT_RULE_EVALUATION`: Timeout for rule evaluation (default: 300ms)
- `WORKER_RETRY_MAX_ATTEMPTS`: Maximum retry attempts (default: 3)
- `WORKER_RETRY_INITIAL_DELAY`: Initial retry delay (default: 100ms)
- `WORKER_RETRY_MAX_DELAY`: Maximum retry delay (default: 5s)
- `WORKER_RETRY_MULTIPLIER`: Retry delay multiplier (default: 2.0)
- `WORKER_QUEUE_POP_TIMEOUT`: Queue pop timeout (default: 1s)
- `WORKER_RULES_RELOAD_INTERVAL`: Rules reload interval (default: 10s)

### UI
- `VITE_API_URL`: API base URL (required, must be set at build time)
- `VITE_API_TIMEOUT`: API request timeout in milliseconds (default: 30000)
- `VITE_API_RETRY_MAX_ATTEMPTS`: Maximum API retry attempts (default: 3)
- `VITE_API_RETRY_INITIAL_DELAY`: Initial API retry delay in milliseconds (default: 1000)
- `VITE_API_RETRY_MAX_DELAY`: Maximum API retry delay in milliseconds (default: 10000)
- `VITE_API_RETRY_MULTIPLIER`: API retry delay multiplier (default: 2.0)
- `VITE_UI_TOAST_DURATION`: Toast notification duration in milliseconds (default: 5000)
- `VITE_UI_POLLING_INTERVAL`: Polling interval for data refresh in milliseconds (default: 10000)
- `VITE_SERVICE_NAME`: Service identification name

## 🔐 Authentication & Authorization

AlgoShield implements a comprehensive authentication and authorization system:

### Roles

- **admin**: Full system access, can manage users, roles, groups, and all rules
- **rule_editor**: Can create, update, and delete rules
- **viewer**: Read-only access (can be extended)

### Groups

Users can be organized into groups for easier permission management. Groups can have roles assigned, which are inherited by group members.

### Permissions

Permissions are managed through roles. Each role defines what actions users can perform:
- Rule management (create, update, delete)
- User management (admin only)
- Transaction viewing
- System administration

### Security Features

- JWT-based authentication with configurable expiration
- Password hashing using bcrypt
- Role-based access control (RBAC)
- Group-based permission inheritance
- User active/inactive status management
- Protected routes in the UI with automatic redirects

## 🏎️ Performance Optimization

AlgoShield is designed for maximum performance:

1. **Compiled with Go 1.25.4** for enhanced performance
2. **Connection pooling** for PostgreSQL (pgx v5) and Redis (go-redis v9)
3. **Rule caching** with Redis to minimize database queries
4. **Schema caching** with in-memory cache and Redis pub/sub invalidation
5. **Async processing** through Redis queues
6. **Horizontal scaling** of worker processes
7. **Optimized database indexes** for fast queries
8. **Hot-reload rules and schemas** without service restart (configurable reload interval, default: 10s)
9. **Configurable timeouts** for transaction processing (default: 300ms) and rule evaluation (default: 300ms)
10. **Retry mechanisms** with exponential backoff
11. **Docker BuildKit** for faster builds with better caching
12. **Parallel builds** for Docker images
13. **OpenTelemetry** infrastructure for observability (metrics and tracing support)

## 🧪 Testing

### Test Strategy

- **API Tests**: Unit and integration tests with race condition detection (`-race` flag)
- **UI Tests**: Vitest-based unit tests with Vue Test Utils
- **Integration Tests**: Database integration tests using testcontainers
- **Benchmarks**: Performance benchmarks for the rules engine
- **Coverage**: Test coverage reporting for both backend and frontend

### Running Tests

Run all tests (API + UI):
```bash
make test
```

Run API tests only (with race detector):
```bash
make test-api
```

Run UI tests only:
```bash
make test-ui
```

Run UI tests with coverage:
```bash
cd src/ui && npm run test:coverage
```

Run UI tests with UI (interactive):
```bash
cd src/ui && npm run test:ui
```

Run rules engine benchmarks:
```bash
make bench
```

Run linters:
```bash
make lint
```

## 📦 Building

Build all Docker images in parallel:
```bash
make build
```

Build only changed services (fast incremental build):
```bash
make build-fast
```

Install all dependencies (Go + npm):
```bash
make install
```

Start all services with Docker Compose:
```bash
make up
```

Stop all services:
```bash
make down
```

View service logs:
```bash
make logs
```

## 🚢 Deployment

### Docker Compose Production

Update `docker-compose.yml` with production settings and deploy:
```bash
docker-compose -f docker-compose.yml up -d
```

Or use the Makefile:
```bash
make up
```

**Note**: Ensure all required environment variables are set in your `.env` file. See `.env.example` for reference.

### Kubernetes

Helm charts and Kubernetes manifests coming soon!

## 📈 Monitoring

Health check endpoints:

- `/health`: Overall system health
- `/ready`: Readiness check

Example response:
```json
{
  "status": "ok",
  "timestamp": "2024-12-05T10:00:00Z",
  "postgres": "healthy",
  "redis": "healthy"
}
```

### Observability

AlgoShield includes OpenTelemetry infrastructure for metrics and tracing:
- OpenTelemetry SDK integration for Go services
- Metrics and tracing support (implementation in progress)
- HTTP instrumentation middleware
- Ready for integration with observability platforms

## 🔄 CI/CD

The project includes GitHub Actions workflows for continuous integration:

- **Backend CI**: Runs on changes to `src/api/`, `src/workers/`, `src/pkg/`, and Go dependencies
  - Linting with golangci-lint
  - Unit and integration tests with race detector
  - Docker image builds
  - Test coverage reporting

- **Frontend CI**: Runs on changes to `src/ui/` and UI dependencies
  - TypeScript type checking
  - Linting (ESLint)
  - Unit tests with Vitest
  - Docker image builds
  - Test coverage reporting

Workflows trigger on pushes and pull requests to `main` and `develop` branches.

## 🏛️ Architecture & Development

### Code Organization

AlgoShield follows a **vertical slice architecture** pattern, organizing code by feature rather than by technical layer. This improves maintainability and code readability.

### Key Patterns

- **Dependency Injection**: All handlers and services use dependency injection for better testability and modularity
- **Validation**: All handlers include request validation as a mandatory practice
- **Error Handling**: Comprehensive error handling with proper context wrapping
- **Structured Logging**: All services use structured logging with configurable log levels
- **Connection Pooling**: Reusable PostgreSQL and Redis connections for optimal performance

### Development Guidelines

- Follow SOLID principles strictly
- Use vertical slice architecture for new features
- Implement proper dependency injection
- Always add validation to handlers
- Use race condition flags (`-race`) when running tests
- Update `.env.example` when adding new environment variables
- Standardize base components in the frontend for consistency
- Remove unused code and components regularly

### Git Workflow

- **Main branch**: `main` (production-ready code)
- Feature branches: Created from `main`, merged via PR
- Commit conventions: Use prefixes like `feat:`, `fix:`, `chore:`, `docs:`, etc.
- Git hooks: Install via `./scripts/install-hooks.sh` (includes pre-commit checks)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Install Git hooks: `./scripts/install-hooks.sh`
4. Make your changes following the development guidelines above
5. Run tests: `make test`
6. Run linters: `make lint`
7. Commit your changes (`git commit -m 'feat: Add some AmazingFeature'`)
8. Push to the branch (`git push origin feature/AmazingFeature`)
9. Open a Pull Request

**Note**: The project uses OpenSpec for spec-driven development. For significant changes, see `openspec/AGENTS.md` for guidance on creating change proposals.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Built with:
- [Go 1.25.4](https://golang.org/) - Programming language
- [Fiber v2](https://gofiber.io/) - High-performance web framework
- [pgx v5](https://github.com/jackc/pgx) - PostgreSQL driver with connection pooling
- [go-redis v9](https://github.com/redis/go-redis) - Redis client
- [expr-lang v1.17.7](https://github.com/expr-lang/expr) - Expression evaluation engine
- [golang-jwt v5](https://github.com/golang-jwt/jwt) - JWT authentication
- [OpenTelemetry](https://opentelemetry.io/) - Observability framework
- [Vue.js 3.5+](https://vuejs.org/) - Progressive JavaScript framework
- [TypeScript 5.9](https://www.typescriptlang.org/) - Typed JavaScript
- [Vuetify 3.7](https://vuetifyjs.com/) - Material Design component framework
- [Pinia 3.0](https://pinia.vuejs.org/) - State management
- [Vue Router 4.6](https://router.vuejs.org/) - Routing
- [Tailwind CSS 4.1](https://tailwindcss.com/) - Utility-first CSS framework
- [Vite 7.2](https://vitejs.dev/) - Next-generation build tool
- [Vitest 4.0](https://vitest.dev/) - Fast unit test framework
- [Font Awesome 6.5](https://fontawesome.com/) - Icon library
- [Prism.js 1.30](https://prismjs.com/) - Syntax highlighting
- [PostgreSQL 16](https://www.postgresql.org/) - Database
- [Redis 7](https://redis.io/) - Caching and message queue

## 📧 Support

For questions and support, please open an issue on GitHub.

---

Made with ❤️ for the fraud prevention community
