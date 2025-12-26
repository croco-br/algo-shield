# 🛡️ AlgoShield

**AlgoShield** is an open-source, high-performance fraud detection and anti-money laundering (AML) transaction analysis tool designed to process transactions with ultra-low latency (<50ms).

## 🎯 Key Features

- **⚡ Ultra-Fast Processing**: Process transactions in <50ms with highly optimized Go workers
- **🔧 Custom Rules Engine**: Configure custom fraud detection rules with an intuitive UI
- **🔄 Hot-Reload Rules**: Update rules in real-time without restarting services
- **📊 Risk Scoring**: Flexible scoring system supporting OK/NOK or numeric scores
- **🧪 Synthetic Data Generation**: Generate test data to validate rules before production
- **🎯 Dual Processing Modes**: Support for pre-transaction (fraud prevention) and post-transaction (AML) analysis
- **🚀 High Scalability**: Horizontally scalable worker architecture
- **📈 Real-time Analysis**: Process events through Redis queues with minimal latency
- **🔐 Authentication & Authorization**: JWT-based authentication with role-based access control (RBAC)
- **👥 User Management**: Complete user, role, and group management system
- **🛡️ Permission System**: Fine-grained permissions for rule editing and administrative tasks

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
- **Worker Service**: Transaction processing engine with custom rules evaluation and hot-reload support
- **UI**: Vue.js 3-based modern web interface with Pinia state management for rule management and user administration
- **PostgreSQL**: Primary data store for transactions, rules, users, roles, and groups
- **Redis**: Message queue for async processing and rules caching

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
  "from_account": "ACC001",
  "to_account": "ACC002",
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

## 🔧 Rule Types

### Amount Rule
Checks transaction amount against threshold:
```json
{
  "type": "amount",
  "conditions": {
    "amount_threshold": 10000
  }
}
```

### Velocity Rule
Checks transaction frequency:
```json
{
  "type": "velocity",
  "conditions": {
    "transaction_count": 10,
    "time_window_seconds": 3600
  }
}
```

### Blocklist Rule
Blocks specific accounts:
```json
{
  "type": "blocklist",
  "conditions": {
    "blocklisted_accounts": ["ACC123", "ACC456"]
  }
}
```

### Pattern Rule
Matches transaction patterns:
```json
{
  "type": "pattern",
  "conditions": {
    "pattern": "international_transfer"
  }
}
```

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
4. **Async processing** through Redis queues
5. **Horizontal scaling** of worker processes
6. **Optimized database indexes** for fast queries
7. **Hot-reload rules** without service restart
8. **Configurable timeouts** for transaction processing and rule evaluation
9. **Retry mechanisms** with exponential backoff

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
- [Pinia](https://pinia.vuejs.org/) - State management
- [Vue Router](https://router.vuejs.org/) - Routing
- [Tailwind CSS](https://tailwindcss.com/) - Styling
- [Vite](https://vitejs.dev/) - Build tool
- [Redis](https://redis.io/) - Caching and message queue
- [JWT](https://jwt.io/) - Authentication tokens

## 📧 Support

For questions and support, please open an issue on GitHub.

---

Made with ❤️ for the fraud prevention community
