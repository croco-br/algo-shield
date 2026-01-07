# Design: Synthetic Events and Dashboard Enhancement

## Context
This change adds synthetic event generation, enhanced transaction filtering, approval workflow, real-time updates, and a performance-optimized dashboard. The system processes transactions asynchronously via workers, and operators need better tools to manage and monitor transactions.

## Goals / Non-Goals

### Goals
- Generate synthetic events from configured schemas with type-aware randomization
- Enable comprehensive filtering of transactions
- Allow operators to approve transactions in review
- Provide real-time transaction updates in the UI
- Create a fast dashboard (<500ms) with aggregated metrics
- Maintain schema_id linkage for all transactions

### Non-Goals
- Complex event generation templates (start simple with type-based randomization)
- Advanced analytics beyond basic aggregations
- Historical trend analysis (focus on current state)
- Multi-user approval workflows (single operator approval)
- Real-time collaboration features

## Decisions

### Decision: Schema ID in Transactions
**What**: Add `schema_id` UUID column to transactions table, nullable initially for migration.

**Why**: 
- Enables filtering transactions by schema
- Allows schema-specific analytics
- Provides traceability from transaction to schema definition

**Alternatives considered**:
- Store schema_id in metadata JSONB - rejected because it's harder to query and index
- Infer schema from transaction structure - rejected because it's unreliable and slow

### Decision: Real-time Updates via Polling
**What**: Use polling with WebSocket fallback option for real-time transaction updates.

**Why**:
- Simpler implementation than WebSocket infrastructure
- Works with existing HTTP infrastructure
- Can upgrade to WebSocket later if needed
- Polling interval: 2-5 seconds for active views

**Alternatives considered**:
- WebSocket only - rejected due to added complexity
- Server-Sent Events (SSE) - considered but polling is simpler for MVP

### Decision: Dashboard Caching Strategy
**What**: Cache dashboard metrics in Redis with 30-second TTL, invalidate on transaction updates.

**Why**:
- Ensures <500ms response time even with high volume
- Reduces database load
- 30s TTL balances freshness with performance

**Alternatives considered**:
- No caching - rejected due to performance requirements
- Materialized views - considered but adds complexity, caching is simpler
- Longer TTL (60s+) - rejected because metrics need reasonable freshness

### Decision: Synthetic Event Randomization
**What**: Generate random values based on field type:
- `string`: Random alphanumeric strings, or use sample value as template
- `number`: Random number in reasonable range (0-1000000) or use sample as base
- `boolean`: Random true/false
- `array`: Generate 1-5 random elements based on array item type
- `object`: Recursively generate nested object fields

**Why**:
- Simple and predictable
- Works with any schema structure
- Type-aware ensures valid data

**Alternatives considered**:
- Faker library - considered but adds dependency, simple randomization is sufficient
- Template-based generation - rejected as too complex for MVP

### Decision: Filter Implementation
**What**: Support filters for: status, schema_id, date range (created_at), amount range.

**Why**:
- Covers most common use cases
- All filterable fields are indexed
- Simple query construction

**Alternatives considered**:
- Full-text search - rejected as not needed for MVP
- Complex query builder - rejected, keep it simple

### Decision: Dashboard Visualizations
**What**: Two basic visualizations:
1. Status distribution (pie/bar chart)
2. Temporal view (line/bar chart for 24h/7d/30d)

**Why**:
- Covers key operational metrics
- Simple to implement and understand
- Performance-friendly aggregations
- Aligns with current system abstractions (no risk level concept exists)

**Alternatives considered**:
- Risk level distribution - rejected, risk level concept doesn't exist in current abstractions
- More complex charts - rejected, start simple
- Geographic visualizations - rejected, not in scope

## Risks / Trade-offs

### Risk: Dashboard Performance with High Volume
**Mitigation**: 
- Use indexed aggregations
- Implement Redis caching
- Consider materialized views if caching insufficient
- Monitor query performance and optimize

### Risk: Real-time Updates Performance
**Mitigation**:
- Use efficient polling interval (not too frequent)
- Only poll when TransactionsView is active
- Use pagination to limit data transfer

### Risk: Schema ID Migration
**Mitigation**:
- Make schema_id nullable initially
- Backfill existing transactions with NULL (acceptable)
- Add NOT NULL constraint after migration period

### Trade-off: Real-time vs Performance
- Polling adds some server load but is simpler
- Can upgrade to WebSocket later if needed
- Acceptable trade-off for MVP

## Migration Plan

1. **Phase 1: Database**
   - Add schema_id column (nullable)
   - Add indexes
   - No data migration needed (existing transactions can have NULL)

2. **Phase 2: Backend**
   - Update models and services
   - Add synthetic event generation
   - Add filtering and approval
   - Add dashboard endpoints

3. **Phase 3: Frontend**
   - Update transactions view
   - Add dashboard
   - Add real-time updates

4. **Phase 4: Worker Update**
   - Ensure worker preserves schema_id when processing

5. **Rollback**: 
   - Remove schema_id column (if needed)
   - Revert API endpoints
   - Frontend changes are backward compatible

## Open Questions
- Should synthetic events be queued immediately or generated on-demand?
  - **Decision**: Generate and queue immediately for simplicity
- What should be the default polling interval for real-time updates?
  - **Decision**: 3 seconds when view is active, pause when inactive
- Should dashboard auto-refresh or manual only?
  - **Decision**: Auto-refresh every 30s with manual refresh option

