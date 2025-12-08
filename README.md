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

## 🏗️ Architecture

AlgoShield is built with a modern microservices architecture:

```
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│     UI      │──────│     API     │──────│   Worker    │
│  (Svelte)   │      │   (Fiber)   │      │   (Rules)   │
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

- **API Service**: RESTful API built with Fiber (Go) for high-performance HTTP handling
- **Worker Service**: Transaction processing engine with custom rules evaluation
- **UI**: SvelteKit-based modern web interface for rule management
- **PostgreSQL**: Primary data store for transactions and rules
- **Redis**: Message queue for async processing and rules caching

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.23+ (for local development)
- Node.js 20+ (for UI development)

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
- **UI**: http://localhost:5173
- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health

### Local Development

1. Install Git hooks (recommended):
```bash
./src/scripts/install-hooks.sh
```

2. Install dependencies:
```bash
make deps
```

2. Start PostgreSQL and Redis:
```bash
docker-compose up -d postgres redis
```

3. Run database migrations:
```bash
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=algoshield
export POSTGRES_PASSWORD=algoshield_secret
export POSTGRES_DB=algoshield

psql -h localhost -U algoshield -d algoshield -f src/scripts/migrations/001_initial_schema.sql
```

4. Start the API:
```bash
make run-api
```

5. Start the Worker (in another terminal):
```bash
make run-worker
```

6. Start the UI (in another terminal):
```bash
make dev-ui
```

## 📖 API Documentation

### Process Transaction

Submit a transaction for analysis:

```bash
POST /api/v1/transactions
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
```

### List Transactions

```bash
GET /api/v1/transactions?limit=50&offset=0
```

### Create Rule

```bash
POST /api/v1/rules
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

```bash
PUT /api/v1/rules/{id}
```

### Delete Rule

```bash
DELETE /api/v1/rules/{id}
```

### List Rules

```bash
GET /api/v1/rules
```

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

### Blacklist Rule
Blocks specific accounts:
```json
{
  "type": "blacklist",
  "conditions": {
    "blacklisted_accounts": ["ACC123", "ACC456"]
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

Configuration is managed through environment variables:

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

### Worker
- `WORKER_CONCURRENCY`: Number of concurrent workers (default: 10)
- `WORKER_BATCH_SIZE`: Batch processing size (default: 100)

## 🏎️ Performance Optimization

AlgoShield is designed for maximum performance:

1. **Compiled with Go 1.23** using `GOEXPERIMENT=greenteagc,rangefunc` for enhanced performance
2. **Connection pooling** for PostgreSQL and Redis
3. **Rule caching** with Redis to minimize database queries
4. **Async processing** through Redis queues
5. **Horizontal scaling** of worker processes
6. **Optimized database indexes** for fast queries

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
- [SvelteKit](https://kit.svelte.dev/) - UI framework
- [Redis](https://redis.io/) - Caching and message queue

## 📧 Support

For questions and support, please open an issue on GitHub.

---

Made with ❤️ for the fraud prevention community
