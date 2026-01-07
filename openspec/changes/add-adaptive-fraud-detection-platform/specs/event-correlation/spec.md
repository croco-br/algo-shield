## ADDED Requirements

### Requirement: Baseline Calculation
The system SHALL calculate and maintain behavioral baselines for entities to enable deviation detection.

#### Scenario: Calculate entity baseline
- **WHEN** nightly baseline calculation runs
- **THEN** 30-day rolling baseline is calculated for each active entity
- **AND** baseline includes mean, standard deviation, and percentiles for key metrics
- **AND** metrics include transaction frequency, average amount, time-of-day distribution

#### Scenario: Baseline for new entities
- **WHEN** entity has insufficient history (< 30 days)
- **THEN** baseline uses available data with confidence factor
- **AND** baseline is marked as provisional
- **AND** baseline updates as more data becomes available

#### Scenario: Baseline storage
- **WHEN** baselines are calculated
- **THEN** baselines are stored in entity_baselines table
- **AND** baselines are cached in Redis for fast rule evaluation access
- **AND** baseline history is retained for trend analysis

### Requirement: Behavior Deviation Detection
The system SHALL detect when entity behavior deviates significantly from established baselines.

#### Scenario: Detect amount deviation
- **WHEN** rule calls `behaviorDeviation(entity_id, "transaction_amount")`
- **THEN** current transaction amount is compared to entity's baseline
- **AND** deviation is calculated as number of standard deviations from mean
- **AND** function returns true if deviation > 3 standard deviations

#### Scenario: Detect frequency deviation
- **WHEN** rule calls `behaviorDeviation(entity_id, "transaction_frequency")`
- **THEN** recent transaction frequency is compared to baseline
- **AND** function detects both sudden increases (velocity attacks) and decreases (account takeover)
- **AND** function returns deviation magnitude

#### Scenario: Detect time-of-day deviation
- **WHEN** analyzing temporal patterns
- **THEN** transaction time is compared to entity's typical activity hours
- **AND** unusual times (e.g., middle of night for business account) are flagged
- **AND** deviation considers day-of-week patterns

### Requirement: Sequence Pattern Detection
The system SHALL detect specific sequences of transactions that match known fraud patterns.

#### Scenario: Define sequence pattern
- **WHEN** administrator defines fraud sequence pattern
- **THEN** pattern specifies ordered sequence of transaction types or amounts
- **AND** pattern includes time constraints between steps
- **AND** pattern can include wildcards and ranges

#### Scenario: Detect structuring pattern
- **WHEN** rule calls `sequenceMatch(entity_id, "structuring")`
- **THEN** recent transactions are analyzed for structuring pattern
- **AND** pattern: multiple transactions just below reporting threshold within short timeframe
- **AND** function returns true if pattern matches

#### Scenario: Detect test-then-execute pattern
- **WHEN** analyzing transaction sequences
- **THEN** pattern of small test transaction followed by large transaction is detected
- **AND** time window between transactions is considered
- **AND** pattern is common in card testing fraud

#### Scenario: Sequence detection performance
- **WHEN** detecting sequences
- **THEN** detection completes within 50ms for 100-transaction history
- **AND** recent transactions are prioritized in search
- **AND** pattern matching uses efficient algorithms

### Requirement: Correlation Between Entities
The system SHALL measure correlation in behavior across multiple related entities.

#### Scenario: Calculate correlation score
- **WHEN** rule calls `correlationScore([entity_id_1, entity_id_2, entity_id_3])`
- **THEN** transaction timing correlation is calculated
- **AND** transaction amount correlation is calculated
- **AND** correlation score is returned (0.0 to 1.0)

#### Scenario: Detect coordinated behavior
- **WHEN** multiple graph-connected entities show highly correlated behavior
- **THEN** coordinated activity is flagged
- **AND** correlation strength and entity relationships are recorded
- **AND** potential account ring or organized fraud is indicated

#### Scenario: Geographic correlation
- **WHEN** analyzing entity correlation
- **THEN** geographic patterns are considered
- **AND** entities transacting from same locations within short timeframe are flagged
- **AND** unusual geographic clustering is detected

### Requirement: Temporal Clustering
The system SHALL detect bursts of activity that deviate from normal patterns.

#### Scenario: Detect transaction burst
- **WHEN** rule calls `temporalCluster(entity_id, "1h")`
- **THEN** transaction frequency in last hour is compared to baseline hourly frequency
- **AND** burst is detected if frequency > 3x baseline
- **AND** burst magnitude and duration are returned

#### Scenario: Detect coordinated bursts
- **WHEN** multiple related entities show simultaneous bursts
- **THEN** coordinated burst is detected
- **AND** burst timing alignment is measured
- **AND** potential distributed attack is flagged

#### Scenario: Dormancy followed by activity
- **WHEN** entity shows long dormancy period followed by sudden activity
- **THEN** dormancy-activation pattern is detected
- **AND** pattern is common in account takeover
- **AND** reactivation is flagged for review

### Requirement: Multi-Event Context
The system SHALL provide transaction context including related historical events for rule evaluation.

#### Scenario: Include event history in context
- **WHEN** evaluating rules for transaction
- **THEN** context includes last N transactions for entity (configurable, default N=10)
- **AND** context includes entity baseline metrics
- **AND** context includes related entity activity

#### Scenario: Cross-entity context
- **WHEN** transaction involves multiple entities
- **THEN** context includes activity for all involved entities
- **AND** relationship paths between entities are included
- **AND** correlated activity across entities is highlighted

#### Scenario: Context performance
- **WHEN** loading transaction context
- **THEN** context loading completes within 20ms
- **AND** context data is cached for reuse in subsequent rules
- **AND** cache invalidation occurs after transaction evaluation

### Requirement: Correlation Functions in Rules
The system SHALL provide correlation functions that can be used in rule expressions.

#### Scenario: Behavior deviation function
- **WHEN** rule uses `behaviorDeviation(entity_id, "amount")` in expression
- **THEN** function returns deviation score (number of std devs)
- **AND** rule can use threshold like `behaviorDeviation(customer_id, "amount") > 3.0`
- **AND** function execution is optimized for rule engine

#### Scenario: Sequence match function
- **WHEN** rule uses `sequenceMatch(entity_id, pattern_name)` in expression
- **THEN** function returns boolean indicating pattern match
- **AND** rule can trigger on specific patterns
- **AND** pattern definitions are loaded from configuration

#### Scenario: Correlation function
- **WHEN** rule uses `correlationScore(entity_ids)` in expression
- **THEN** function returns correlation coefficient (0.0-1.0)
- **AND** rule can detect coordinated behavior
- **AND** function efficiently queries entity activity

#### Scenario: Temporal cluster function
- **WHEN** rule uses `temporalCluster(entity_id, time_window)` in expression
- **THEN** function returns burst magnitude
- **AND** rule can detect velocity attacks
- **AND** function uses cached frequency calculations

### Requirement: Correlation Analytics Dashboard
The system SHALL provide dashboard for analyzing correlation patterns and temporal trends.

#### Scenario: View entity behavior trends
- **WHEN** viewing entity detail page
- **THEN** behavior metrics are charted over time
- **AND** baseline is shown for comparison
- **AND** deviations are highlighted on timeline

#### Scenario: View correlation graph
- **WHEN** investigating coordinated behavior
- **THEN** graph visualization shows correlated entities
- **AND** edge thickness indicates correlation strength
- **AND** temporal alignment is visualized

#### Scenario: View sequence pattern matches
- **WHEN** analyzing pattern detection results
- **THEN** matched sequences are displayed with timeline
- **AND** pattern steps are highlighted
- **AND** similar historical patterns are shown

