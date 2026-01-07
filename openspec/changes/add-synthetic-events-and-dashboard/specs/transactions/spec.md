## ADDED Requirements

### Requirement: Transaction Schema Association
Each transaction SHALL be associated with an event schema via schema_id. The schema_id SHALL be set when a transaction is created and SHALL be preserved through processing.

#### Scenario: Transaction created with schema
- **WHEN** a transaction is created from an event
- **THEN** the transaction's schema_id is set based on the event's schema
- **AND** the schema_id is stored in the database

#### Scenario: Schema ID preserved through processing
- **WHEN** a transaction is processed by the worker
- **THEN** the schema_id is preserved in the final transaction record
- **AND** the transaction can be filtered and analyzed by schema

### Requirement: Transaction Filtering
The system SHALL provide filtering capabilities for transactions based on multiple criteria including status, schema, date range, and amount range.

#### Scenario: Filter by status
- **WHEN** an operator filters transactions by status
- **THEN** only transactions matching the selected status are returned
- **AND** the filter can be combined with other filters

#### Scenario: Filter by schema
- **WHEN** an operator filters transactions by schema_id
- **THEN** only transactions associated with that schema are returned
- **AND** the filter supports selecting multiple schemas

#### Scenario: Filter by date range
- **WHEN** an operator filters transactions by created_at date range
- **THEN** only transactions within the specified date range are returned
- **AND** the date range filter uses indexed queries for performance

#### Scenario: Filter by amount range
- **WHEN** an operator filters transactions by amount range
- **THEN** only transactions with amounts within the specified range are returned
- **AND** the filter respects currency if specified

#### Scenario: Multiple filters combined
- **WHEN** an operator applies multiple filters simultaneously
- **THEN** transactions matching all filter criteria are returned
- **AND** the query is optimized using appropriate indexes

### Requirement: Transaction Approval
Operators SHALL be able to approve transactions that are in review status. Only transactions with status "in_review" SHALL be eligible for approval.

#### Scenario: Approve transaction in review
- **WHEN** an operator approves a transaction with status "in_review"
- **THEN** the transaction status is changed to "approved"
- **AND** the approval is recorded with timestamp and operator information
- **AND** the updated transaction is returned

#### Scenario: Cannot approve non-review transaction
- **WHEN** an operator attempts to approve a transaction that is not in "in_review" status
- **THEN** the system returns an error indicating the transaction cannot be approved
- **AND** the transaction status remains unchanged

### Requirement: Real-time Transaction Updates
The transactions view SHALL update automatically to reflect new transactions and status changes as they are processed by the worker.

#### Scenario: New transaction appears
- **WHEN** the worker processes a new transaction
- **THEN** the transaction appears in the transactions view without manual refresh
- **AND** the update occurs within 3-5 seconds of processing

#### Scenario: Transaction status updates
- **WHEN** a transaction's status changes (e.g., from pending to in_review)
- **THEN** the transactions view reflects the updated status automatically
- **AND** visual indicators show which transactions were recently updated

#### Scenario: Real-time updates pause when view inactive
- **WHEN** the transactions view is not visible or active
- **THEN** real-time updates are paused to conserve resources
- **AND** updates resume when the view becomes active again

