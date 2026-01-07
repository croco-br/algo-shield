# Change: Add Synthetic Events Generation and Enhanced Dashboard

## Why
Operators need the ability to generate synthetic transaction events for testing and development purposes based on configured schemas. Additionally, the transactions view needs enhanced filtering capabilities, real-time updates, and the ability to approve transactions in review. Finally, a dashboard with aggregated metrics is needed to provide quick insights into system health and transaction patterns, with strict performance requirements (<500ms response time).

## What Changes
- **ADDED**: Synthetic event generation capability - generate random transaction events from configured schemas
- **ADDED**: Transaction filtering system - multiple filter options for finding specific transactions
- **ADDED**: Transaction approval workflow - operators can approve transactions in review status
- **ADDED**: Real-time transaction updates - transactions view updates automatically as worker processes transactions
- **ADDED**: Dashboard with aggregated metrics - overview of key indicators with 3 basic visualization options
- **MODIFIED**: Transaction model - add `schema_id` field to link transactions to their event schema
- **MODIFIED**: Transactions API - add filtering, approval endpoint, and real-time update support
- **MODIFIED**: Dashboard view - transform from simple transaction list to aggregated metrics dashboard

## Impact
- **Affected specs**: 
  - `synthetic-events` (new capability)
  - `transactions` (enhanced with filters, approval, real-time)
  - `dashboard` (new capability)
- **Affected code**:
  - `src/api/internal/transactions/` - add filters, approval endpoint
  - `src/api/internal/schemas/` - add synthetic event generation
  - `src/pkg/models/transaction.go` - add schema_id field
  - `src/ui/src/views/DashboardView.vue` - transform to metrics dashboard
  - `src/ui/src/views/TransactionsView.vue` - add filters, real-time updates, approval
  - Database migration - add schema_id to transactions table, add indexes for performance
  - Worker - ensure schema_id is preserved when processing transactions
- **Performance considerations**:
  - Dashboard queries must be optimized with proper indexes and aggregation
  - Real-time updates must use efficient pub/sub mechanism
  - Filter queries must use indexed columns

