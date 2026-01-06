# 🛡️ AlgoShield

**AlgoShield** is an open-source, high-performance fraud detection and anti-money laundering (AML) transaction analysis tool designed to process transactions with ultra-low latency (<50ms).

## 🎯 Key Features

- **⚡ Ultra-Fast Processing**: Process transactions in <50ms with highly optimized Go workers
- **🔧 Custom Rules Engine**: Configure custom fraud detection rules with an intuitive UI
- **🔄 Hot-Reload Rules**: Update rules in real-time without restarting services
- **📋 Event Schema Management**: Define and manage event schemas with automatic field extraction from sample JSON
- **📊 Risk Scoring**: Flexible scoring system supporting OK/NOK or numeric scores
- **🧪 Synthetic Data Generation**: Generate test data to validate rules before production
- **🎯 Dual Processing Modes**: Support for pre-transaction (fraud prevention) and post-transaction (AML) analysis
- **🚀 High Scalability**: Horizontally scalable worker architecture
- **📈 Real-time Analysis**: Process events through Redis queues with minimal latency
- **🔐 Authentication & Authorization**: JWT-based authentication with role-based access control (RBAC)
- **👥 User Management**: Complete user, role, and group management system
- **🛡️ Permission System**: Fine-grained permissions for rule editing and administrative tasks
- **🎨 Branding Configuration**: White-label customization with configurable colors, logos, and app name

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

- **API Service**: RESTful API built with Fiber (Go) for high-performance HTTP handling with JWT authentication
- **Worker Service**: Transaction processing engine with custom rules evaluation, schema management, and hot-reload support
- **UI**: Vue.js 3-based modern web interface with Vuetify (Material Design) components, Pinia state management, and Tailwind CSS for rule management, schema management, and user administration
- **PostgreSQL**: Primary data store for transactions, rules, event schemas, users, roles, and groups
- **Redis**: Message queue for async processing, rules caching, and schema invalidation pub/sub

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.25.4+ (for local development)
- Node.js ^20.19.0 or >=22.12.0 (for UI development)

### Using Docker Compose (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/yourusername/algo-shield.git
cd algo-shield
```

2. Start all services:
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
```

**Note**: The migrations script (`migrations.sh`) is designed for Docker environments. For local development, run migrations manually as shown above.

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

```bash
POST /api/v1/rules
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "High Value Transaction",
  "description": "Flag transactions over $10,000",
  "type": "amount",
  "action": "review",
  "priority": 10,
  "enabled": true,
  "conditions": {
    "amount_threshold": 10000
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

AlgoShield uses a flexible custom expression-based rule system. All rules use the `custom` type and evaluate expressions using [expr-lang](https://github.com/expr-lang/expr) against event schemas.

### Basic Rule Structure

```json
{
  "name": "High Value Transaction",
  "description": "Flag transactions over $10,000",
  "type": "custom",
  "action": "review",
  "score": 30,
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
- **block**: Block the transaction immediately
- **review**: Flag for manual review
- **score**: Add risk score without blocking

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
- `VITE_API_URL`: API base URL (default: http://localhost:8080)
- `VITE_API_TIMEOUT`: API request timeout
- `VITE_API_RETRY_MAX_ATTEMPTS`: Maximum API retry attempts
- `VITE_API_RETRY_INITIAL_DELAY`: Initial API retry delay
- `VITE_API_RETRY_MAX_DELAY`: Maximum API retry delay
- `VITE_API_RETRY_MULTIPLIER`: API retry delay multiplier
- `VITE_UI_TOAST_DURATION`: Toast notification duration
- `VITE_UI_POLLING_INTERVAL`: Polling interval for data refresh

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
2. **Connection pooling** for PostgreSQL and Redis
3. **Rule caching** with Redis to minimize database queries
4. **Schema caching** with in-memory cache and Redis pub/sub invalidation
5. **Async processing** through Redis queues
6. **Horizontal scaling** of worker processes
7. **Optimized database indexes** for fast queries
8. **Hot-reload rules and schemas** without service restart
9. **Configurable timeouts** for transaction processing and rule evaluation
10. **Retry mechanisms** with exponential backoff

## 🧪 Testing

Run tests:
```bash
make test
```

## 📦 Building

Build all binaries:
```bash
make build
```

Build Docker images:
```bash
make docker-build
```

## 🚢 Deployment

### Docker Compose Production

Update `docker-compose.yml` with production settings and deploy:
```bash
docker-compose -f docker-compose.yml up -d
```

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

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Built with:
- [Go](https://golang.org/) - Programming language
- [Fiber](https://gofiber.io/) - Web framework
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [Vue.js](https://vuejs.org/) - UI framework
- [Vuetify](https://vuetifyjs.com/) - Material Design component framework
- [Pinia](https://pinia.vuejs.org/) - State management
- [Vue Router](https://router.vuejs.org/) - Routing
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework
- [Vite](https://vitejs.dev/) - Build tool
- [Redis](https://redis.io/) - Caching and message queue
- [JWT](https://jwt.io/) - Authentication tokens
- [Font Awesome](https://fontawesome.com/) - Icon library
- [Prism.js](https://prismjs.com/) - Syntax highlighting

## 📧 Support

For questions and support, please open an issue on GitHub.

---

Made with ❤️ for the fraud prevention community
