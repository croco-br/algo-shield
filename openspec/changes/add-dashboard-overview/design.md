# Design: Dashboard Overview

## Context

AlgoShield is a fraud detection/AML system with a target transaction processing latency of <50ms. The dashboard needs to provide operational insights to analysts without impacting core transaction processing performance. The frontend uses Vue 3 with Vuetify, and the backend follows a vertical slice architecture with Go/Fiber.

**Stakeholders**: Fraud analysts, compliance officers, system administrators

**Constraints**:
- Dashboard API response must be <300ms
- Must not impact transaction processing latency
- Must follow existing vertical slice architecture pattern
- Must follow existing migration pattern (sequential numbered SQL files)

## Goals / Non-Goals

**Goals**:
- Provide at-a-glance operational dashboard with key metrics
- Fast dashboard loading (<300ms) using caching and optimized queries
- Separate transactions page from dashboard overview
- Follow existing codebase patterns and conventions
- Real-time streaming updates
- Custom date range filtering

**Non-Goals**:
- Historical trend charts (can be added in future iteration)
- Dashboard customization/widget arrangement

## Decisions

### Decision 1: Single API Endpoint for All Dashboard Data
**What**: Return all dashboard statistics (summary, performance, rules) in a single API call.
**Why**: Reduces HTTP roundtrips, simplifies frontend loading state, enables efficient Redis caching as a single unit.
**Alternatives considered**:
- Multiple endpoints (one per widget): More REST-ful but adds latency from multiple roundtrips
- GraphQL: Overkill for this use case, not aligned with existing patterns

### Decision 2: Redis Caching with 30-second TTL
**What**: Cache computed dashboard stats in Redis with a short TTL.
**Why**: 
- Dashboard data is aggregate/statistical, slight staleness is acceptable
- Eliminates repeated expensive aggregate queries
- Aligns with existing Redis caching pattern used for rules
**Alternatives considered**:
- No caching, rely on database indexes: Risk of slow queries under load
- Longer TTL (5 min): Data would be too stale for operational dashboard
- Materialized views: More complex, harder to cache, PostgreSQL-specific

### Decision 3: Single Optimized Query vs Multiple Parallel Queries
**What**: Use a single PostgreSQL query with CTEs to fetch all aggregations in one roundtrip.
**Why**: Reduces database roundtrips and connection usage; PostgreSQL CTEs are efficient for this pattern.
**Alternatives considered**:
- Parallel queries: More Go code complexity, connection pool pressure
- Pre-computed stats table: Adds write-path complexity, sync issues

### Decision 4: Add Index on `risk_level` Column
**What**: Add B-tree index on `transactions.risk_level` via migration.
**Why**: Existing schema has indexes on `status` and `created_at`, but grouping by `risk_level` requires a scan without an index.
**Note**: Index on `processing_time` not needed initially; `percentile_cont` works efficiently with existing data size.

### Decision 5: Rules Triggered via JSONB Aggregation
**What**: Use PostgreSQL's `jsonb_array_elements_text()` to unnest `matched_rules` array and aggregate counts.
**Why**: 
- No schema change required; `matched_rules` is already stored as JSONB array
- PostgreSQL has efficient JSONB operators
- Query is cacheable, so complex aggregation only runs on cache miss
**Trade-off**: JSONB unnest is slower than a normalized join table, but caching mitigates this.

### Decision 6: Server-Sent Events (SSE) for Real-time Updates
**What**: Use SSE via `/api/v1/dashboard/stream` endpoint for pushing live updates to the frontend.
**Why**:
- SSE is simpler than WebSockets for unidirectional server-to-client updates
- Native browser support (EventSource API), no additional libraries needed
- Fiber has good SSE support
- Aligns with the read-only nature of dashboard updates (no client-to-server messages needed)
**Alternatives considered**:
- WebSockets: Bidirectional overkill; more complex connection management
- Polling: Higher latency, more server load, less efficient
- Long polling: More complex than SSE with no clear benefit

**Implementation**:
- Backend publishes dashboard invalidation events to Redis pub/sub when transactions are processed
- SSE endpoint subscribes to Redis channel and pushes updates to connected clients
- Frontend uses EventSource to receive updates and refresh stats

### Decision 7: Date Range Filtering via Query Parameters
**What**: Accept `start_date` and `end_date` query parameters on the stats endpoint.
**Why**:
- Simple REST-style filtering, easy to cache per date range
- Allows analysts to focus on specific time periods
- Date-filtered queries benefit from existing `created_at` index
**Alternatives considered**:
- Predefined periods only (last_24h, last_7d): Less flexible
- POST body with filter object: Breaks REST conventions for GET endpoint

**Cache key strategy**: Include date range in Redis cache key (e.g., `dashboard:stats:2026-01-01:2026-01-02`)

## API Response Schema

### Stats Endpoint: `GET /api/v1/dashboard/stats`

**Query Parameters**:
- `start_date` (optional): ISO 8601 date (e.g., `2026-01-01`)
- `end_date` (optional): ISO 8601 date (e.g., `2026-01-02`)
- If omitted, defaults to all-time statistics

```json
{
  "summary": {
    "total_transactions": 15420,
    "by_status": {
      "pending": 120,
      "approved": 14500,
      "rejected": 350,
      "review": 450
    },
    "by_risk_level": {
      "low": 12000,
      "medium": 2800,
      "high": 620
    }
  },
  "performance": {
    "avg_processing_time_ms": 42.5,
    "p95_processing_time_ms": 78.2,
    "target_latency_ms": 50,
    "within_target_percent": 87.3
  },
  "rules_triggered": [
    { "rule_name": "High Amount Transaction", "trigger_count": 2340, "percentage": 15.2 },
    { "rule_name": "Transaction Velocity Check", "trigger_count": 1890, "percentage": 12.3 }
  ],
  "filter": {
    "start_date": "2026-01-01T00:00:00Z",
    "end_date": "2026-01-02T23:59:59Z"
  },
  "generated_at": "2026-01-02T15:30:00Z",
  "cached": true
}
```

### Stream Endpoint: `GET /api/v1/dashboard/stream`

Server-Sent Events stream that pushes updates when dashboard data changes.

**Event format**:
```
event: stats_updated
data: {"type": "stats_updated", "timestamp": "2026-01-02T15:30:05Z"}

event: heartbeat
data: {"type": "heartbeat", "timestamp": "2026-01-02T15:30:10Z"}
```

**Client behavior**: On receiving `stats_updated` event, frontend re-fetches stats from the stats endpoint with current date range filters.

## SQL Query Design

```sql
-- Single query with CTEs for all dashboard aggregations (with date range filter)
-- $1 = start_date, $2 = end_date (NULL for all-time)
WITH filtered_transactions AS (
    SELECT *
    FROM transactions
    WHERE ($1::timestamptz IS NULL OR created_at >= $1)
      AND ($2::timestamptz IS NULL OR created_at <= $2)
),
status_counts AS (
    SELECT status, COUNT(*) as count
    FROM filtered_transactions
    GROUP BY status
),
risk_counts AS (
    SELECT risk_level, COUNT(*) as count
    FROM filtered_transactions
    GROUP BY risk_level
),
performance_stats AS (
    SELECT 
        AVG(processing_time) as avg_processing_time,
        PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY processing_time) as p95_processing_time,
        COUNT(*) FILTER (WHERE processing_time <= 50) * 100.0 / NULLIF(COUNT(*), 0) as within_target_percent
    FROM filtered_transactions
    WHERE processed_at IS NOT NULL
),
rules_triggered AS (
    SELECT 
        rule_name,
        COUNT(*) as trigger_count
    FROM filtered_transactions,
        LATERAL jsonb_array_elements_text(matched_rules) AS rule_name
    GROUP BY rule_name
    ORDER BY trigger_count DESC
    LIMIT 10
)
SELECT 
    (SELECT COUNT(*) FROM filtered_transactions) as total_transactions,
    (SELECT json_object_agg(status, count) FROM status_counts) as by_status,
    (SELECT json_object_agg(risk_level, count) FROM risk_counts) as by_risk_level,
    (SELECT row_to_json(p.*) FROM performance_stats p) as performance,
    (SELECT json_agg(r.*) FROM rules_triggered r) as rules_triggered;
```

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Cache stampede on expiry | Use single-flight pattern or lock when refreshing cache |
| Query performance degrades with millions of rows | Date-range filter is now supported; default to last 30 days if all-time is slow |
| JSONB unnest slow for large matched_rules | Date filtering reduces dataset; consider denormalized table in future |
| Frontend loading jank | Use skeleton loaders, optimize component rendering |
| SSE connection overhead | Use heartbeat to detect stale connections; limit max connections per user |
| Redis pub/sub message loss | SSE triggers refetch, not direct data push; missing events just delay update |
| Large date ranges slow | Cache per date range; consider limiting range to max 90 days |

## Migration Plan

1. Apply migration `007_dashboard_indexes.sql` (adds risk_level index)
2. Deploy backend with dashboard vertical slice (stats + stream endpoints)
3. Deploy frontend with new dashboard view (including SSE integration and date picker)
4. Monitor dashboard API response times via observability

**Rollback**: 
- Frontend can fall back to showing transactions page
- Backend endpoints can be removed without breaking existing functionality
- SSE can be disabled; frontend falls back to manual refresh
- Index can be dropped if causing write performance issues (unlikely)

## Open Questions

- [x] ~~Should dashboard show "last 24h" vs "last 7d" vs "all time"?~~ → User can select custom date range
- [x] ~~Should we add auto-refresh polling on the dashboard?~~ → Using SSE for real-time updates
- [ ] What is the maximum allowed date range? (Suggest: 90 days to limit query cost)
- [ ] Should we show a "live" indicator when SSE is connected?
