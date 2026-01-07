## 1. Database Changes
- [ ] 1.1 Create migration to add `schema_id` column to transactions table
- [ ] 1.2 Add foreign key constraint to event_schemas table
- [ ] 1.3 Add indexes for filtering (status, schema_id, created_at ranges)
- [ ] 1.4 Add indexes for dashboard aggregations (status, created_at)

## 2. Backend - Synthetic Events
- [ ] 2.1 Create synthetic event generator service
- [ ] 2.2 Implement field type-based randomization (string, number, boolean, etc.)
- [ ] 2.3 Create API endpoint POST /api/v1/schemas/{id}/generate-events
- [ ] 2.4 Add validation for generation parameters (count, seed for reproducibility)
- [ ] 2.5 Integrate with transaction processing to set schema_id

## 3. Backend - Transactions Enhancement
- [ ] 3.1 Update Transaction model to include schema_id
- [ ] 3.2 Add filtering support to ListTransactions (status, schema_id, date range, amount range)
- [ ] 3.3 Create approval endpoint PATCH /api/v1/transactions/{id}/approve
- [ ] 3.4 Add validation for approval (only in_review status can be approved)
- [ ] 3.5 Update worker to preserve schema_id when processing transactions
- [ ] 3.6 Add Redis pub/sub for real-time transaction updates

## 4. Backend - Dashboard
- [ ] 4.1 Create dashboard service with aggregation queries
- [ ] 4.2 Implement status distribution aggregation
- [ ] 4.3 Implement temporal aggregation (24h, 7d, 30d)
- [ ] 4.4 Optimize queries with proper indexes and materialized views if needed
- [ ] 4.5 Create API endpoint GET /api/v1/dashboard/metrics
- [ ] 4.6 Add caching layer for dashboard metrics (Redis, TTL: 30s)

## 5. Frontend - Synthetic Events
- [ ] 5.1 Add "Generate Events" button/modal in SchemasView
- [ ] 5.2 Create event generation form (count, seed optional)
- [ ] 5.3 Show generation progress/status
- [ ] 5.4 Display generated events count

## 6. Frontend - Transactions Enhancement
- [ ] 6.1 Create TransactionsView component (separate from DashboardView)
- [ ] 6.2 Add filter UI (status, schema, date range, amount range)
- [ ] 6.3 Add "Approve" button for in_review transactions
- [ ] 6.4 Implement real-time updates using WebSocket or polling
- [ ] 6.5 Update transaction table to show schema_id/schema name
- [ ] 6.6 Add visual indicators for real-time updates

## 7. Frontend - Dashboard
- [ ] 7.1 Transform DashboardView to show aggregated metrics
- [ ] 7.2 Create status distribution visualization (pie/bar chart)
- [ ] 7.3 Create temporal visualization (line/bar chart for time series)
- [ ] 7.4 Add loading states and error handling
- [ ] 7.5 Implement auto-refresh (every 30s) with manual refresh option
- [ ] 7.6 Add performance monitoring to ensure <500ms load time

## 8. Testing & Validation
- [ ] 8.1 Test synthetic event generation with various schema types
- [ ] 8.2 Test transaction filtering with multiple combinations
- [ ] 8.3 Test approval workflow (only in_review can be approved)
- [ ] 8.4 Test real-time updates with multiple concurrent transactions
- [ ] 8.5 Load test dashboard with high volume data
- [ ] 8.6 Verify dashboard response time <500ms under load
- [ ] 8.7 Test schema_id preservation through worker processing

