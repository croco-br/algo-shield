# AlgoShield API Examples

This document provides practical examples of using the AlgoShield API.

## Base URL

```
http://localhost:8080
```

## Health Checks

### Check System Health

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "timestamp": "2024-12-05T10:00:00Z",
  "postgres": "healthy",
  "redis": "healthy"
}
```

### Check Readiness

```bash
curl http://localhost:8080/ready
```

## Transaction Management

### Submit a Transaction for Analysis

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "external_id": "txn_001",
    "amount": 1500.00,
    "currency": "USD",
    "from_account": "ACC12345",
    "to_account": "ACC67890",
    "type": "transfer",
    "metadata": {
      "ip_address": "192.168.1.100",
      "device_id": "device_abc123",
      "user_agent": "Mozilla/5.0"
    },
    "timestamp": "2024-12-05T10:00:00Z"
  }'
```

Response:
```json
{
  "status": "queued",
  "external_id": "txn_001",
  "processing_time": 3,
  "message": "Transaction queued for processing"
}
```

### Get Transaction by ID

```bash
# Replace {transaction_id} with actual UUID
curl http://localhost:8080/api/v1/transactions/{transaction_id}
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "external_id": "txn_001",
  "amount": 1500.00,
  "currency": "USD",
  "from_account": "ACC12345",
  "to_account": "ACC67890",
  "type": "transfer",
  "status": "approved",
  "risk_score": 25.0,
  "risk_level": "low",
  "processing_time_ms": 15,
  "matched_rules": ["High Amount Transaction"],
  "metadata": {
    "ip_address": "192.168.1.100",
    "device_id": "device_abc123"
  },
  "created_at": "2024-12-05T10:00:00Z",
  "processed_at": "2024-12-05T10:00:00.015Z"
}
```

### List All Transactions

```bash
# Default: limit=50, offset=0
curl http://localhost:8080/api/v1/transactions

# With pagination
curl "http://localhost:8080/api/v1/transactions?limit=20&offset=0"
```

Response:
```json
{
  "transactions": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "external_id": "txn_001",
      "amount": 1500.00,
      "status": "approved",
      "risk_score": 25.0,
      "risk_level": "low"
    }
  ],
  "limit": 20,
  "offset": 0
}
```

## Rule Management

### Create a New Rule

#### Amount Threshold Rule

```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "High Value Alert",
    "description": "Flag transactions over $10,000",
    "type": "amount",
    "action": "review",
    "priority": 10,
    "enabled": true,
    "conditions": {
      "amount_threshold": 10000
    },
    "score": 50
  }'
```

#### Velocity Check Rule

```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Rapid Fire Transactions",
    "description": "Flag accounts with >10 transactions per hour",
    "type": "velocity",
    "action": "score",
    "priority": 20,
    "enabled": true,
    "conditions": {
      "transaction_count": 10,
      "time_window_seconds": 3600
    },
    "score": 40
  }'
```

#### Blacklist Rule

```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Blocked Accounts",
    "description": "Block transactions from blacklisted accounts",
    "type": "blacklist",
    "action": "block",
    "priority": 1,
    "enabled": true,
    "conditions": {
      "blacklisted_accounts": ["ACC99999", "ACC88888"]
    },
    "score": 100
  }'
```

#### Pattern Match Rule

```bash
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "International Transfer Check",
    "description": "Review international transfers",
    "type": "pattern",
    "action": "review",
    "priority": 15,
    "enabled": true,
    "conditions": {
      "pattern": "international_transfer"
    },
    "score": 30
  }'
```

### Get Rule by ID

```bash
# Replace {rule_id} with actual UUID
curl http://localhost:8080/api/v1/rules/{rule_id}
```

### List All Rules

```bash
curl http://localhost:8080/api/v1/rules
```

Response:
```json
{
  "rules": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "High Value Alert",
      "description": "Flag transactions over $10,000",
      "type": "amount",
      "action": "review",
      "priority": 10,
      "enabled": true,
      "conditions": {
        "amount_threshold": 10000
      },
      "score": 50,
      "created_at": "2024-12-05T10:00:00Z",
      "updated_at": "2024-12-05T10:00:00Z"
    }
  ]
}
```

### Update a Rule

```bash
# Replace {rule_id} with actual UUID
curl -X PUT http://localhost:8080/api/v1/rules/{rule_id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated High Value Alert",
    "description": "Flag transactions over $15,000",
    "type": "amount",
    "action": "block",
    "priority": 5,
    "enabled": true,
    "conditions": {
      "amount_threshold": 15000
    },
    "score": 75
  }'
```

### Delete a Rule

```bash
# Replace {rule_id} with actual UUID
curl -X DELETE http://localhost:8080/api/v1/rules/{rule_id}
```

## Complete Workflow Example

### 1. Create Rules

```bash
# Create amount rule
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Large Transaction",
    "description": "Block transactions over $50,000",
    "type": "amount",
    "action": "block",
    "priority": 5,
    "enabled": true,
    "conditions": {"amount_threshold": 50000},
    "score": 100
  }'

# Create velocity rule
curl -X POST http://localhost:8080/api/v1/rules \
  -H "Content-Type: application/json" \
  -d '{
    "name": "High Frequency",
    "description": "Flag high transaction frequency",
    "type": "velocity",
    "action": "score",
    "priority": 10,
    "enabled": true,
    "conditions": {
      "transaction_count": 5,
      "time_window_seconds": 600
    },
    "score": 40
  }'
```

### 2. Submit Transactions

```bash
# Normal transaction
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "external_id": "txn_safe_001",
    "amount": 500.00,
    "currency": "USD",
    "from_account": "ACC001",
    "to_account": "ACC002",
    "type": "transfer",
    "metadata": {},
    "timestamp": "2024-12-05T10:00:00Z"
  }'

# High-value transaction (should be blocked)
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "external_id": "txn_large_001",
    "amount": 75000.00,
    "currency": "USD",
    "from_account": "ACC001",
    "to_account": "ACC003",
    "type": "transfer",
    "metadata": {},
    "timestamp": "2024-12-05T10:01:00Z"
  }'
```

### 3. Check Transaction Status

```bash
# Wait a moment for processing, then check
curl http://localhost:8080/api/v1/transactions
```

## Testing with JQ

If you have `jq` installed, you can format responses:

```bash
curl -s http://localhost:8080/api/v1/rules | jq '.'

curl -s http://localhost:8080/api/v1/transactions | jq '.transactions[] | {external_id, status, risk_score}'
```

## Rate Testing

### Simple Load Test

```bash
# Send 100 transactions
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/v1/transactions \
    -H "Content-Type: application/json" \
    -d "{
      \"external_id\": \"txn_load_$i\",
      \"amount\": $((RANDOM % 10000 + 100)),
      \"currency\": \"USD\",
      \"from_account\": \"ACC001\",
      \"to_account\": \"ACC002\",
      \"type\": \"transfer\",
      \"metadata\": {},
      \"timestamp\": \"2024-12-05T10:00:00Z\"
    }" &
done
wait
```

## Error Responses

### Invalid Request

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{"invalid": "data"}'
```

Response (400):
```json
{
  "error": "Invalid request body"
}
```

### Not Found

```bash
curl http://localhost:8080/api/v1/transactions/invalid-uuid
```

Response (404):
```json
{
  "error": "Transaction not found"
}
```

## Performance Metrics

Each transaction response includes processing time:

```json
{
  "processing_time": 3
}
```

This indicates the time in milliseconds it took to queue the transaction. The actual rule evaluation happens asynchronously in the worker and is recorded in the transaction record.

