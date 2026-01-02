## ADDED Requirements

### Requirement: Dashboard Statistics API
The system SHALL provide a `/api/v1/dashboard/stats` endpoint that returns aggregated dashboard statistics in a single API call.

#### Scenario: Authenticated user requests dashboard stats
- **WHEN** an authenticated user sends GET `/api/v1/dashboard/stats`
- **THEN** the system returns HTTP 200 with JSON containing summary cards, performance metrics, and rules triggered data
- **AND** the response time SHALL be less than 300ms

#### Scenario: Unauthenticated request
- **WHEN** an unauthenticated user sends GET `/api/v1/dashboard/stats`
- **THEN** the system returns HTTP 401 Unauthorized

### Requirement: Dashboard Summary Cards
The dashboard statistics response SHALL include summary card data with transaction counts grouped by status and risk level.

#### Scenario: Summary cards data structure
- **WHEN** dashboard stats are retrieved
- **THEN** the response includes `summary` object with:
  - `total_transactions`: total count of all transactions
  - `by_status`: object with counts for `pending`, `approved`, `rejected`, `review`
  - `by_risk_level`: object with counts for `low`, `medium`, `high`

#### Scenario: Empty database
- **WHEN** no transactions exist in the database
- **THEN** the summary returns zero counts for all categories

### Requirement: Dashboard Performance Metrics
The dashboard statistics response SHALL include performance metrics showing processing time statistics.

#### Scenario: Performance metrics data structure
- **WHEN** dashboard stats are retrieved
- **THEN** the response includes `performance` object with:
  - `avg_processing_time_ms`: average processing time in milliseconds
  - `p95_processing_time_ms`: 95th percentile processing time in milliseconds
  - `target_latency_ms`: the target latency SLA (50ms)
  - `within_target_percent`: percentage of transactions processed within target latency

#### Scenario: No processed transactions
- **WHEN** no transactions have been processed (all pending)
- **THEN** performance metrics return null or zero values with appropriate indication

### Requirement: Dashboard Rules Triggered Statistics
The dashboard statistics response SHALL include statistics about which rules are triggered most frequently.

#### Scenario: Rules triggered data structure
- **WHEN** dashboard stats are retrieved
- **THEN** the response includes `rules_triggered` array with:
  - Each item containing `rule_name`, `trigger_count`, and `percentage` of total
  - Array sorted by `trigger_count` descending
  - Limited to top 10 most triggered rules

#### Scenario: No rules triggered
- **WHEN** no transactions have matched any rules
- **THEN** the `rules_triggered` array is empty

### Requirement: Dashboard Date Range Filtering
The dashboard statistics endpoint SHALL support filtering by date range via query parameters.

#### Scenario: Custom date range filter
- **WHEN** user sends GET `/api/v1/dashboard/stats?start_date=2026-01-01&end_date=2026-01-02`
- **THEN** the system returns statistics only for transactions created within the specified date range
- **AND** the response includes a `filter` object showing the applied date range

#### Scenario: Partial date range (start only)
- **WHEN** user sends GET `/api/v1/dashboard/stats?start_date=2026-01-01`
- **THEN** the system returns statistics for transactions from the start date to now

#### Scenario: Partial date range (end only)
- **WHEN** user sends GET `/api/v1/dashboard/stats?end_date=2026-01-02`
- **THEN** the system returns statistics for all transactions up to the end date

#### Scenario: No date range specified
- **WHEN** user sends GET `/api/v1/dashboard/stats` without date parameters
- **THEN** the system returns all-time statistics

#### Scenario: Invalid date format
- **WHEN** user sends GET `/api/v1/dashboard/stats?start_date=invalid`
- **THEN** the system returns HTTP 400 Bad Request with error message

### Requirement: Dashboard Statistics Caching
The system SHALL cache dashboard statistics in Redis to ensure consistent sub-300ms response times.

#### Scenario: Cache hit
- **WHEN** dashboard stats are requested and cache is valid for the requested date range
- **THEN** the system returns cached data without querying the database
- **AND** response includes `cached: true` indicator

#### Scenario: Cache miss or expired
- **WHEN** dashboard stats are requested and cache is empty or expired
- **THEN** the system queries the database, caches the result, and returns fresh data
- **AND** response includes `cached: false` indicator

#### Scenario: Cache TTL
- **WHEN** dashboard stats are cached
- **THEN** the cache expires after 30 seconds to balance freshness with performance

#### Scenario: Cache key includes date range
- **WHEN** dashboard stats are cached with a date range filter
- **THEN** the cache key includes the date range to prevent serving incorrect cached data for different filters

### Requirement: Dashboard Real-time Updates via SSE
The system SHALL provide a `/api/v1/dashboard/stream` endpoint using Server-Sent Events to push real-time update notifications.

#### Scenario: SSE connection established
- **WHEN** an authenticated user connects to GET `/api/v1/dashboard/stream`
- **THEN** the server establishes an SSE connection and keeps it open
- **AND** the response Content-Type is `text/event-stream`

#### Scenario: Stats update notification
- **WHEN** a transaction is processed by the worker
- **THEN** the system publishes a dashboard invalidation event via Redis pub/sub
- **AND** all connected SSE clients receive a `stats_updated` event

#### Scenario: SSE heartbeat
- **WHEN** an SSE connection is open
- **THEN** the server sends periodic heartbeat events (every 30 seconds) to keep the connection alive

#### Scenario: Unauthenticated SSE request
- **WHEN** an unauthenticated user attempts to connect to `/api/v1/dashboard/stream`
- **THEN** the system returns HTTP 401 Unauthorized

#### Scenario: Client reconnection
- **WHEN** an SSE connection is lost
- **THEN** the frontend automatically reconnects using EventSource's built-in retry mechanism

### Requirement: Dashboard UI Overview Page
The frontend SHALL display a dashboard overview page at `/dashboard` with visual widgets for summary cards, performance metrics, and rules triggered.

#### Scenario: Dashboard page loads
- **WHEN** user navigates to `/dashboard`
- **THEN** the page displays loading state, fetches dashboard stats, and renders widgets

#### Scenario: Summary cards display
- **WHEN** dashboard stats are loaded
- **THEN** summary cards show total transactions, counts by status, and counts by risk level with appropriate visual styling (colors for risk levels, icons for status)

#### Scenario: Performance metrics display
- **WHEN** dashboard stats are loaded
- **THEN** performance section shows average processing time, P95 latency, and percentage within SLA target with visual indicator (green if within target, red if not)

#### Scenario: Rules triggered display
- **WHEN** dashboard stats are loaded
- **THEN** rules triggered section shows a ranked list of top rules with trigger counts and percentages

### Requirement: Dashboard Date Range Picker
The frontend SHALL provide a date range picker to filter dashboard statistics.

#### Scenario: Date range picker UI
- **WHEN** dashboard page loads
- **THEN** a date range picker is visible allowing selection of start and end dates

#### Scenario: Date range selection
- **WHEN** user selects a date range
- **THEN** the dashboard re-fetches statistics with the selected date range filter
- **AND** widgets update to show filtered data

#### Scenario: Clear date filter
- **WHEN** user clears the date range filter
- **THEN** the dashboard shows all-time statistics

### Requirement: Dashboard Real-time Updates Integration
The frontend SHALL integrate with the SSE stream to receive real-time update notifications.

#### Scenario: SSE connection on load
- **WHEN** dashboard page loads
- **THEN** the frontend establishes an SSE connection to `/api/v1/dashboard/stream`

#### Scenario: Auto-refresh on update
- **WHEN** frontend receives a `stats_updated` event
- **THEN** the dashboard automatically re-fetches statistics with current date range filter

#### Scenario: Connection status indicator
- **WHEN** SSE connection is active
- **THEN** a visual indicator shows the dashboard is receiving live updates

#### Scenario: SSE cleanup on navigation
- **WHEN** user navigates away from the dashboard page
- **THEN** the SSE connection is closed to free resources

### Requirement: Transactions Page Separation
The system SHALL provide a separate transactions list page accessible at `/transactions`.

#### Scenario: Transactions page navigation
- **WHEN** user clicks "Transactions" in sidebar
- **THEN** user is navigated to `/transactions` showing the transaction list table

#### Scenario: Dashboard and transactions are independent
- **WHEN** user navigates between `/dashboard` and `/transactions`
- **THEN** each page loads independently with its own data and UI
