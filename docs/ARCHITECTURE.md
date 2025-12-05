# AlgoShield Architecture

## System Overview

AlgoShield is designed as a high-performance, distributed system for real-time fraud detection and anti-money laundering (AML) analysis.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         Client Layer                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │   Web    │  │ Mobile   │  │   API    │  │  Other   │   │
│  │    UI    │  │   App    │  │ Clients  │  │ Systems  │   │
│  └─────┬────┘  └─────┬────┘  └─────┬────┘  └─────┬────┘   │
└────────┼─────────────┼─────────────┼─────────────┼─────────┘
         │             │             │             │
         └─────────────┴─────────────┴─────────────┘
                       │
         ┌─────────────▼──────────────┐
         │        API Gateway          │
         │      (Fiber / Go)           │
         │  - Authentication           │
         │  - Rate Limiting            │
         │  - Request Validation       │
         └─────────────┬───────────────┘
                       │
         ┌─────────────▼──────────────┐
         │       API Service           │
         │  - Transaction Ingestion    │
         │  - Rule Management          │
         │  - Query Interface          │
         └─────┬───────────────────────┘
               │
               │ (Async via Redis Queue)
               │
         ┌─────▼──────────────────────┐
         │      Redis Queue            │
         │  - transaction:queue        │
         │  - rules:cache              │
         └─────┬───────────────────────┘
               │
               │ (Multiple Workers)
               │
    ┌──────────┴──────────┬────────────┐
    │                     │            │
┌───▼─────┐      ┌───────▼───┐  ┌────▼─────┐
│ Worker  │      │ Worker    │  │ Worker   │
│ #1      │      │ #2        │  │ #N       │
│         │      │           │  │          │
│ Rules   │      │ Rules     │  │ Rules    │
│ Engine  │      │ Engine    │  │ Engine   │
└───┬─────┘      └───────┬───┘  └────┬─────┘
    │                    │            │
    └────────────┬───────┴────────────┘
                 │
         ┌───────▼────────────────────┐
         │    PostgreSQL Database      │
         │  - Transactions             │
         │  - Rules                    │
         │  - Audit Logs               │
         └─────────────────────────────┘
```

## Component Details

### 1. API Service (Fiber + Go)

**Responsibilities:**
- Accept transaction events via REST API
- Manage rules (CRUD operations)
- Provide query interface for transactions
- Health checks and monitoring

**Key Features:**
- Ultra-fast HTTP handling with Fiber
- Connection pooling for database
- CORS support for web clients
- Structured logging

**Performance Optimizations:**
- Zero-allocation router
- Prefork mode support for multi-core
- Keep-alive connections
- Response caching where appropriate

### 2. Worker Service (Go)

**Responsibilities:**
- Pull transactions from Redis queue
- Evaluate transactions against rules
- Calculate risk scores
- Store results in database
- Send notifications (future)

**Key Features:**
- Configurable concurrency
- Hot-reload of rules (10-second intervals)
- Graceful shutdown
- Error handling and retries

**Performance Optimizations:**
- Goroutine-based workers
- Rules cached in Redis
- Batch processing capability
- In-memory rule evaluation

### 3. Rules Engine

**Responsibilities:**
- Load and cache rules from database
- Evaluate transactions against rules
- Calculate risk scores
- Determine transaction status

**Rule Types:**
- **Amount**: Threshold-based checks
- **Velocity**: Time-window transaction counting
- **Blacklist**: Account blocking
- **Pattern**: Custom pattern matching
- **Custom**: Extensible for future rules

**Evaluation Logic:**
```
For each transaction:
  1. Load all enabled rules (from cache)
  2. Sort rules by priority
  3. For each rule:
     a. Evaluate conditions
     b. If matched:
        - Add to matched rules list
        - Add score to total
        - Apply action (allow/block/review/score)
  4. Determine final risk level based on score
  5. Return result
```

### 4. Data Stores

#### PostgreSQL
- **Primary data store** for transactions and rules
- **Connection pooling** with pgx (50 max connections)
- **Optimized indexes** for common queries
- **JSONB columns** for flexible metadata

**Tables:**
- `transactions`: All processed transactions
- `rules`: Fraud detection rules

#### Redis
- **Message queue** for async processing
- **Rules cache** for fast access (5-minute TTL)
- **Session storage** (future)
- **Real-time metrics** (future)

### 5. UI (SvelteKit)

**Responsibilities:**
- Rule management interface
- Transaction monitoring
- Dashboard and analytics (future)

**Key Features:**
- Modern, responsive design
- Real-time updates
- Form validation
- Optimistic UI updates

## Data Flow

### Transaction Processing Flow

```
1. Client submits transaction
   POST /api/v1/transactions
   │
2. API validates and queues
   └─> Redis LPUSH transaction:queue
       │
3. Worker picks up transaction
   └─> Redis BRPOP transaction:queue
       │
4. Load rules (cached)
   └─> Redis GET rules:cache
       │ (cache miss)
       └─> PostgreSQL SELECT rules
           │
5. Evaluate rules
   └─> Rules Engine processes
       │
6. Store result
   └─> PostgreSQL INSERT transaction
       │
7. Client queries result
   GET /api/v1/transactions/{id}
   └─> PostgreSQL SELECT transaction
```

### Rule Update Flow

```
1. User updates rule via UI
   │
2. API receives update
   PUT /api/v1/rules/{id}
   │
3. Update database
   └─> PostgreSQL UPDATE rules
       │
4. Invalidate cache
   └─> Redis DEL rules:cache
       │
5. Workers reload rules (within 10s)
   └─> PostgreSQL SELECT rules
       └─> Redis SET rules:cache
```

## Performance Characteristics

### Latency Targets

- **API Response**: <5ms (queue operation only)
- **Transaction Processing**: <50ms (end-to-end)
- **Rule Evaluation**: <30ms (average)
- **Database Queries**: <10ms (with proper indexes)

### Throughput

- **API**: 10,000+ requests/second per instance
- **Workers**: 1,000+ transactions/second per worker
- **Horizontal Scaling**: Near-linear with worker count

### Resource Usage

- **API Service**: ~50MB RAM, <5% CPU (idle)
- **Worker Service**: ~100MB RAM, 10-20% CPU per worker
- **PostgreSQL**: Varies with data volume
- **Redis**: ~100MB RAM (typical)

## Scalability

### Horizontal Scaling

1. **API Service**: Run multiple instances behind load balancer
2. **Worker Service**: Scale replicas up/down based on queue depth
3. **Database**: Read replicas for queries, primary for writes
4. **Redis**: Redis Cluster for high availability

### Vertical Scaling

1. **Increase worker concurrency**: `WORKER_CONCURRENCY`
2. **Increase database connections**: Pool size configuration
3. **Optimize queries**: Add indexes, query optimization

## Security Considerations

### Current

- Input validation on all endpoints
- SQL injection protection (parameterized queries)
- CORS configuration
- Health check endpoints

### Future Enhancements

- Authentication & Authorization (JWT/OAuth2)
- API rate limiting per client
- Request signing
- Encryption at rest
- Audit logging
- Role-based access control (RBAC)

## Monitoring & Observability

### Current

- Health check endpoints
- Structured logging
- Processing time metrics

### Future Enhancements

- Prometheus metrics
- Distributed tracing (OpenTelemetry)
- Grafana dashboards
- Alert manager integration
- Performance profiling

## High Availability

### Database

- PostgreSQL replication (primary-replica)
- Automatic failover with patroni/stolon
- Point-in-time recovery (PITR)

### Redis

- Redis Sentinel for HA
- Redis Cluster for sharding
- Persistent storage (AOF + RDB)

### Application

- Multiple API instances
- Multiple worker instances
- Graceful shutdown handling
- Health-based routing

## Disaster Recovery

1. **Database Backups**: Daily snapshots + continuous WAL archiving
2. **Configuration Backups**: Git-based versioning
3. **Recovery Time Objective (RTO)**: <1 hour
4. **Recovery Point Objective (RPO)**: <5 minutes

## Future Enhancements

1. **Machine Learning Integration**
   - ML-based fraud detection models
   - Anomaly detection
   - Adaptive risk scoring

2. **Advanced Analytics**
   - Real-time dashboards
   - Transaction trends
   - Rule effectiveness metrics

3. **Synthetic Data Generation**
   - Configurable data generators
   - Rule testing framework
   - Load testing support

4. **Notification System**
   - Email/SMS alerts
   - Webhook integrations
   - Slack/Teams notifications

5. **Multi-tenancy**
   - Tenant isolation
   - Per-tenant rules
   - Custom branding

6. **API Versioning**
   - Version management
   - Backward compatibility
   - Deprecation strategy

