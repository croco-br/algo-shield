## ADDED Requirements

### Requirement: Feedback Capture
The system SHALL capture analyst decisions on transactions and alerts with reason codes and timestamps.

#### Scenario: Submit feedback on transaction
- **WHEN** analyst reviews a transaction and marks it as false positive
- **THEN** feedback is recorded with decision ("false_positive")
- **AND** feedback includes analyst_id and timestamp
- **AND** feedback includes optional reason text
- **AND** feedback is linked to transaction_id and any associated alert_id

#### Scenario: Feedback on approved transaction
- **WHEN** analyst confirms a blocked transaction was correctly identified as fraud
- **THEN** feedback is recorded with decision ("confirmed_fraud")
- **AND** feedback validates the triggered rules

#### Scenario: Feedback with reason code
- **WHEN** submitting feedback
- **THEN** analyst can select from predefined reason codes
- **AND** analyst can provide free-text explanation
- **AND** reason codes include common categories (timing, amount, behavior, etc.)

#### Scenario: View feedback history
- **WHEN** viewing transaction details
- **THEN** all feedback entries for that transaction are displayed
- **AND** feedback shows analyst name, decision, reason, and timestamp
- **AND** feedback is displayed chronologically

### Requirement: Rule Effectiveness Metrics
The system SHALL calculate and track effectiveness metrics for each rule based on analyst feedback.

#### Scenario: Calculate false positive rate
- **WHEN** daily metrics calculation runs
- **THEN** false positive rate is calculated per rule as (false_positives / total_triggers)
- **AND** metrics include 7-day, 30-day, and 90-day windows
- **AND** metrics are stored in rule_metrics table

#### Scenario: Calculate precision and recall
- **WHEN** calculating rule metrics
- **THEN** precision is calculated as (true_positives / (true_positives + false_positives))
- **AND** recall is calculated as (true_positives / (true_positives + false_negatives))
- **AND** metrics handle cases where denominators are zero

#### Scenario: Track alert-to-confirmation ratio
- **WHEN** calculating rule metrics
- **THEN** ratio of alerts that resulted in confirmed fraud is calculated
- **AND** ratio indicates rule precision in real-world usage

#### Scenario: Metrics calculation performance
- **WHEN** daily metrics calculation job runs
- **THEN** calculation completes within 5 minutes for 100,000 feedback entries
- **AND** calculation is non-blocking and runs in background

### Requirement: Rule Effectiveness Dashboard
The system SHALL provide a dashboard displaying rule effectiveness metrics for operators and analysts.

#### Scenario: View rule effectiveness summary
- **WHEN** accessing rule metrics dashboard
- **THEN** all rules are listed with key metrics (FP rate, precision, trigger count)
- **AND** rules can be sorted by any metric
- **AND** rules are color-coded by effectiveness (green=good, yellow=review, red=poor)

#### Scenario: Compare rules
- **WHEN** selecting multiple rules for comparison
- **THEN** side-by-side comparison shows all metrics
- **AND** chart visualizes metric trends over time
- **AND** statistical significance of differences is indicated

#### Scenario: Drill down into rule performance
- **WHEN** clicking on a specific rule
- **THEN** detailed view shows metric trends over time
- **AND** recent feedback entries for that rule are displayed
- **AND** transactions that triggered the rule are listed with outcomes

#### Scenario: Filter rules by performance
- **WHEN** filtering rules by false positive rate
- **THEN** only rules meeting filter criteria are displayed
- **AND** filters include thresholds for FP rate, precision, trigger count
- **AND** filtered results can be exported

### Requirement: False Positive Tracking
The system SHALL track false positive patterns to identify noisy rules and improve rule quality.

#### Scenario: Identify high-noise rules
- **WHEN** analyzing rule metrics
- **THEN** rules with false positive rate > 50% are flagged
- **AND** flagged rules trigger notifications to rule administrators
- **AND** recommendations for rule adjustment are generated

#### Scenario: False positive trends
- **WHEN** viewing false positive metrics
- **THEN** trend over time is displayed as line chart
- **AND** sudden increases in FP rate are highlighted
- **AND** alerts are generated if FP rate increases > 20% week-over-week

#### Scenario: False positive reasons analysis
- **WHEN** analyzing false positive feedback
- **THEN** common reason codes are aggregated
- **AND** most frequent reasons are displayed
- **AND** insights help identify rule improvement opportunities

### Requirement: Adaptive Rule Thresholds
The system SHALL support rules with adaptive thresholds that adjust based on feedback and performance metrics.

#### Scenario: Enable adaptive threshold for rule
- **WHEN** administrator enables adaptive mode for a rule
- **THEN** rule threshold is adjusted weekly based on feedback
- **AND** adjustment aims to maintain target false positive rate (e.g., 10%)
- **AND** threshold changes are logged in audit trail

#### Scenario: Threshold adjustment calculation
- **WHEN** calculating threshold adjustment
- **THEN** if FP rate > target, threshold increases (makes rule less sensitive)
- **AND** if FP rate < target and recall is low, threshold decreases (makes rule more sensitive)
- **AND** adjustment is bounded to prevent extreme changes (max 20% per week)

#### Scenario: Manual override of adaptive threshold
- **WHEN** administrator manually sets threshold for adaptive rule
- **THEN** manual override disables adaptive adjustment temporarily
- **AND** override is logged with reason
- **AND** administrator can re-enable adaptive mode

#### Scenario: Monitor adaptive rule behavior
- **WHEN** viewing adaptive rule dashboard
- **THEN** threshold history is displayed as timeline
- **AND** correlation between threshold changes and FP rate is shown
- **AND** adaptive algorithm effectiveness metrics are displayed

### Requirement: Feedback API Integration
The system SHALL provide API endpoints for external systems to submit feedback and query metrics.

#### Scenario: Submit feedback via API
- **WHEN** external case management system submits feedback via API
- **THEN** feedback is validated and stored
- **AND** API returns 201 Created with feedback ID
- **AND** invalid feedback returns 400 Bad Request with error details

#### Scenario: Query rule metrics via API
- **WHEN** external analytics system queries rule metrics
- **THEN** metrics are returned in JSON format
- **AND** API supports filtering by rule ID, date range, metric type
- **AND** response includes pagination for large result sets

#### Scenario: Webhook for metric updates
- **WHEN** rule metrics are recalculated
- **THEN** webhook is triggered to notify external systems
- **AND** webhook payload includes updated metrics
- **AND** webhook retries on failure with exponential backoff

### Requirement: Feedback Data Retention
The system SHALL retain feedback data for audit and analysis purposes with configurable retention periods.

#### Scenario: Configure retention policy
- **WHEN** administrator configures feedback retention
- **THEN** retention period can be set between 90 days and 7 years
- **AND** retention applies to feedback entries and calculated metrics
- **AND** expired data is archived before deletion

#### Scenario: Archive old feedback
- **WHEN** feedback data exceeds retention period
- **THEN** data is moved to archive storage
- **AND** archived data remains queryable but read-only
- **AND** archive process runs nightly and does not impact performance

#### Scenario: Export feedback for compliance
- **WHEN** compliance officer exports feedback data
- **THEN** export includes all feedback with analyst decisions
- **AND** export format is CSV or JSON
- **AND** export is encrypted and access-controlled

