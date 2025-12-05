# AlgoShield Quick Start Guide

Get AlgoShield up and running in 5 minutes!

## Prerequisites

- Docker & Docker Compose installed
- 4GB RAM available
- Ports 5173, 8080, 5432, 6379 available

## Step 1: Start the Application

```bash
# Clone the repository
git clone https://github.com/yourusername/algo-shield.git
cd algo-shield

# Start all services
docker-compose up -d

# Watch the logs (optional)
docker-compose logs -f
```

Wait for all services to be healthy (about 30 seconds).

## Step 2: Verify Installation

```bash
# Check API health
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "postgres": "healthy",
  "redis": "healthy"
}
```

## Step 3: Access the UI

Open your browser and navigate to:
```
http://localhost:5173
```

You should see the AlgoShield dashboard with the pre-configured sample rules.

## Step 4: Submit Your First Transaction

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "external_id": "test_001",
    "amount": 500.00,
    "currency": "USD",
    "from_account": "ACC001",
    "to_account": "ACC002",
    "type": "transfer",
    "metadata": {},
    "timestamp": "2024-12-05T10:00:00Z"
  }'
```

## Step 5: Check Transaction Status

```bash
# List all transactions
curl http://localhost:8080/api/v1/transactions
```

You'll see your transaction with risk score and status!

## Step 6: Create a Custom Rule

### Option 1: Via UI

1. Click "Create Rule" button
2. Fill in the form:
   - Name: "My First Rule"
   - Type: "amount"
   - Action: "review"
   - Conditions: `{"amount_threshold": 1000}`
   - Score: 50
3. Click "Create"

### Option 2: Via API

```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My First Rule",
    "description": "Review transactions over $1,000",
    "type": "amount",
    "action": "review",
    "priority": 10,
    "enabled": true,
    "conditions": {
      "amount_threshold": 1000
    },
    "score": 50
  }'
```

## Step 7: Test Your Rule

Submit a transaction that triggers your rule:

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "external_id": "test_002",
    "amount": 1500.00,
    "currency": "USD",
    "from_account": "ACC001",
    "to_account": "ACC003",
    "type": "transfer",
    "metadata": {},
    "timestamp": "2024-12-05T10:01:00Z"
  }'
```

Check the transaction - it should have status "review" and show your rule in matched_rules!

## Common Tasks

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f worker
```

### Restart Services

```bash
# Restart all
docker-compose restart

# Restart specific service
docker-compose restart api
docker-compose restart worker
```

### Stop Services

```bash
docker-compose down
```

### Stop and Clean Everything

```bash
# Remove containers, networks, and volumes
docker-compose down -v
```

## Development Mode

If you want to run locally for development:

```bash
# Start infrastructure
docker-compose up -d postgres redis

# Install Go dependencies
go mod download

# Install UI dependencies
cd ui && npm install && cd ..

# Run API (terminal 1)
make run-api

# Run Worker (terminal 2)
make run-worker

# Run UI (terminal 3)
make dev-ui
```

## Troubleshooting

### Port Already in Use

If ports are already in use, edit `docker-compose.yml` to change the exposed ports:

```yaml
ports:
  - "8081:8080"  # Changed from 8080:8080
```

### Database Connection Failed

```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# View PostgreSQL logs
docker-compose logs postgres

# Restart PostgreSQL
docker-compose restart postgres
```

### Worker Not Processing

```bash
# Check worker logs
docker-compose logs worker

# Check Redis connection
docker-compose exec redis redis-cli ping

# Check queue length
docker-compose exec redis redis-cli llen transaction:queue
```

### UI Not Loading

```bash
# Check UI logs
docker-compose logs ui

# Rebuild UI
docker-compose up -d --build ui
```

## Next Steps

1. Read the [Architecture Documentation](./ARCHITECTURE.md)
2. Explore [API Examples](./API_EXAMPLES.md)
3. Learn about [Rule Types](../README.md#rule-types)
4. Check [Contributing Guidelines](../CONTRIBUTING.md)

## Performance Benchmarking

Test the system performance:

```bash
# Install Apache Bench (if not installed)
# macOS: brew install ab
# Ubuntu: apt-get install apache2-utils

# Benchmark API
ab -n 1000 -c 10 -T 'application/json' \
   -p test_transaction.json \
   http://localhost:8080/api/v1/transactions
```

Create `test_transaction.json`:
```json
{
  "external_id": "bench_test",
  "amount": 100.00,
  "currency": "USD",
  "from_account": "ACC001",
  "to_account": "ACC002",
  "type": "transfer",
  "metadata": {},
  "timestamp": "2024-12-05T10:00:00Z"
}
```

## Getting Help

- 📖 [Full Documentation](../README.md)
- 🐛 [Report Issues](https://github.com/yourusername/algo-shield/issues)
- 💬 [Discussions](https://github.com/yourusername/algo-shield/discussions)

Happy fraud detecting! 🛡️

