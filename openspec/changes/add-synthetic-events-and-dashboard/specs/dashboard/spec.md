## ADDED Requirements

### Requirement: Dashboard Aggregated Metrics
The system SHALL provide a dashboard view with aggregated metrics showing key indicators about transactions and system health. The dashboard SHALL respond within 500ms even under high data volume conditions.

#### Scenario: Dashboard loads quickly
- **WHEN** an operator accesses the dashboard
- **THEN** the dashboard metrics are displayed within 500ms
- **AND** cached metrics are used when available to ensure performance

#### Scenario: Status distribution displayed
- **WHEN** the dashboard loads
- **THEN** a visualization shows the distribution of transactions by status (approved, rejected, in_review, pending)
- **AND** the visualization includes counts and percentages for each status

#### Scenario: Temporal metrics displayed
- **WHEN** the dashboard loads
- **THEN** a visualization shows transaction volume over time
- **AND** the visualization supports viewing 24 hours, 7 days, or 30 days
- **AND** the time series data is aggregated efficiently

#### Scenario: Dashboard auto-refreshes
- **WHEN** the dashboard is displayed
- **THEN** the metrics automatically refresh every 30 seconds
- **AND** operators can manually refresh at any time
- **AND** the refresh maintains the <500ms response time requirement

#### Scenario: Dashboard performance under load
- **WHEN** the system has high transaction volume (millions of transactions)
- **THEN** the dashboard still responds within 500ms
- **AND** caching and optimized queries are used to achieve this performance
- **AND** database indexes support efficient aggregations

#### Scenario: Dashboard handles errors gracefully
- **WHEN** dashboard metrics cannot be loaded
- **THEN** an error message is displayed
- **AND** operators can retry loading the metrics
- **AND** partial data is displayed if available

