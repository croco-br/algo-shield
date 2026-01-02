# Change: Add Dashboard Overview with Summary Cards, Performance Metrics, Rules Triggered, Real-time Updates, and Date Range Filtering

## Why

The current "Dashboard" page is just a transaction list, providing no at-a-glance insights for fraud analysts and compliance officers. A proper dashboard with summary cards, performance metrics, and rules triggered statistics enables faster decision-making and system health monitoring without navigating to individual transaction details. Real-time updates ensure analysts see the latest data without manual refresh, and date range filtering allows focusing on specific time periods for investigation.

## What Changes

### Backend (API)
- Add new `/api/v1/dashboard/stats` endpoint returning aggregated dashboard data
- Add new `/api/v1/dashboard/stream` SSE endpoint for real-time updates
- Create `dashboard` vertical slice (`internal/dashboard/`) with handler, service, repository
- Add optimized SQL queries for aggregate statistics with date range filtering
- Add Redis caching for dashboard stats (30s TTL) to ensure <300ms response time
- Add Redis pub/sub for broadcasting dashboard invalidation events
- Add database migration for `risk_level` index to optimize aggregations

### Frontend (UI)
- Create new `DashboardView.vue` with summary cards, performance metrics, and rules triggered widgets
- Add date range picker component for filtering dashboard data
- Integrate EventSource for SSE real-time updates
- Rename current `DashboardView.vue` to `TransactionsView.vue`
- Update router to point `/transactions` to the new `TransactionsView.vue`
- Add TypeScript types for dashboard API response

### Navigation
- `/dashboard` → New dashboard overview page
- `/transactions` → Existing transaction list (renamed component)
- Sidebar already has both entries configured

## Impact

- **Affected specs**: None (first spec for dashboard capability)
- **Affected code**:
  - `src/api/internal/dashboard/` (new vertical slice)
  - `src/api/internal/routes/routes.go` (register dashboard endpoints)
  - `src/ui/src/views/DashboardView.vue` (new implementation)
  - `src/ui/src/views/TransactionsView.vue` (renamed from DashboardView.vue)
  - `src/ui/src/router/index.ts` (update routes)
  - `scripts/migrations/007_dashboard_indexes.sql` (new migration)
- **Performance target**: Dashboard API response < 300ms
- **No breaking changes**: Existing `/api/v1/transactions` endpoint unchanged

