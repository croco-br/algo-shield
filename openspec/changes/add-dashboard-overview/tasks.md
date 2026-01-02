# Tasks: Add Dashboard Overview

## 1. Database Migration
- [ ] 1.1 Create `scripts/migrations/007_dashboard_indexes.sql` with index on `risk_level`
- [ ] 1.2 Run migration to apply indexes

## 2. Backend - Dashboard Vertical Slice
- [ ] 2.1 Create `src/api/internal/dashboard/models.go` with DashboardStats and filter structs
- [ ] 2.2 Create `src/api/internal/dashboard/repository.go` with optimized aggregate queries (date range support)
- [ ] 2.3 Create `src/api/internal/dashboard/service.go` with Redis caching logic (keyed by date range)
- [ ] 2.4 Create `src/api/internal/dashboard/handler.go` with GET /stats endpoint (date range query params)
- [ ] 2.5 Add SSE stream handler in `handler.go` for GET /stream endpoint
- [ ] 2.6 Add Redis pub/sub publisher in service for dashboard invalidation events
- [ ] 2.7 Register dashboard routes in `routes.go`

## 3. Worker - Publish Dashboard Invalidation
- [ ] 3.1 Update worker transaction processor to publish invalidation event after transaction processing

## 4. Frontend - Dashboard View
- [ ] 4.1 Rename `DashboardView.vue` to `TransactionsView.vue`
- [ ] 4.2 Add TypeScript types for dashboard API response in `src/ui/src/types/dashboard.ts`
- [ ] 4.3 Create new `DashboardView.vue` with summary cards, performance metrics, and rules triggered widgets
- [ ] 4.4 Add date range picker component to dashboard
- [ ] 4.5 Integrate EventSource for SSE real-time updates
- [ ] 4.6 Update `router/index.ts` to use `TransactionsView.vue` for `/transactions`
- [ ] 4.7 Ensure sidebar navigation works correctly for both routes

## 5. Validation
- [ ] 5.1 Verify dashboard API loads in < 300ms (use browser devtools)
- [ ] 5.2 Verify date range filtering works correctly
- [ ] 5.3 Verify SSE connection establishes and receives updates
- [ ] 5.4 Verify transactions page works independently
- [ ] 5.5 Verify sidebar navigation between dashboard and transactions
