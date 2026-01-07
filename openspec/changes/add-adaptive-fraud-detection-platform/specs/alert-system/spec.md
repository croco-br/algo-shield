## ADDED Requirements

### Requirement: Alert Generation
The system SHALL generate structured alerts from transaction detections with severity, type, and related entities.

#### Scenario: Generate alert from high-risk transaction
- **WHEN** worker detects high-risk transaction with matched rules
- **THEN** alert is generated with unique ID
- **AND** alert includes severity calculated from rules and entity risk
- **AND** alert includes alert type based on primary matched rule
- **AND** alert is linked to triggering transaction and related entities

#### Scenario: Calculate alert severity
- **WHEN** generating alert
- **THEN** base severity is determined by highest-priority matched rule
- **AND** severity is elevated if entity risk score > 70
- **AND** severity is elevated if multiple rules matched
- **AND** final severity is one of LOW, MEDIUM, HIGH, CRITICAL

#### Scenario: Assign alert type
- **WHEN** generating alert
- **THEN** alert type is assigned based on primary rule category
- **AND** types include "high_value", "velocity", "ring_detected", "structuring", "unusual_pattern"
- **AND** custom alert types can be defined

#### Scenario: Link related entities
- **WHEN** generating alert
- **THEN** all entities involved in transaction are linked to alert
- **AND** entity relationships relevant to alert are recorded
- **AND** entity risk scores at time of alert are captured

### Requirement: Alert Grouping
The system SHALL intelligently group related alerts to reduce noise and surface coordinated behavior patterns.

#### Scenario: Group alerts by entity
- **WHEN** multiple alerts involve same entity within 24-hour window
- **THEN** alerts are grouped into single alert group
- **AND** group severity is set to highest individual alert severity
- **AND** all grouped transactions are accessible from group

#### Scenario: Group alerts by pattern
- **WHEN** multiple alerts match same pattern (e.g., multiple velocity violations)
- **THEN** alerts are grouped by pattern type
- **AND** group shows aggregated statistics (total amount, count)
- **AND** individual alerts remain accessible within group

#### Scenario: Group alerts for related entities
- **WHEN** alerts involve graph-connected entities
- **THEN** alerts are grouped if entities are within 2 hops in graph
- **AND** coordinated behavior is surfaced through grouping
- **AND** group metadata includes relationship paths

#### Scenario: Prevent over-grouping
- **WHEN** grouping alerts
- **THEN** maximum group size is enforced (default: 20 alerts)
- **AND** time window for grouping is limited (default: 24 hours)
- **AND** dissimilar alerts are not forced into same group

### Requirement: Alert Status Management
The system SHALL track alert lifecycle through defined statuses with state transitions.

#### Scenario: New alert created
- **WHEN** alert is generated
- **THEN** alert status is set to NEW
- **AND** created_at timestamp is recorded
- **AND** alert appears in new alerts queue

#### Scenario: Triage alert
- **WHEN** analyst reviews alert and determines it needs investigation
- **THEN** alert status changes from NEW to TRIAGED
- **AND** analyst_id is recorded
- **AND** triage_notes can be added

#### Scenario: Assign alert to case
- **WHEN** alert is escalated to formal investigation
- **THEN** alert status changes to INVESTIGATING
- **AND** case is created and linked to alert
- **AND** alert is removed from triage queue

#### Scenario: Close alert
- **WHEN** alert investigation is complete
- **THEN** alert status changes to CLOSED
- **AND** resolution type is recorded (confirmed_fraud, false_positive, no_action)
- **AND** closed_at timestamp is recorded

#### Scenario: Invalid state transitions
- **WHEN** attempting invalid status transition (e.g., NEW → CLOSED without TRIAGED)
- **THEN** transition is rejected with error
- **AND** valid next states are provided
- **AND** status remains unchanged

### Requirement: Alert Assignment
The system SHALL support manual and automatic assignment of alerts to analysts for investigation.

#### Scenario: Manual alert assignment
- **WHEN** supervisor assigns alert to specific analyst
- **THEN** alert assigned_to field is set to analyst user_id
- **AND** analyst receives notification
- **AND** alert appears in analyst's work queue

#### Scenario: Auto-assign alerts
- **WHEN** auto-assignment is enabled
- **THEN** new alerts are automatically assigned using round-robin or load-based strategy
- **AND** analyst workload is balanced
- **AND** high-priority alerts are assigned first

#### Scenario: Reassign alert
- **WHEN** alert is reassigned to different analyst
- **THEN** assigned_to is updated
- **AND** both old and new analyst receive notifications
- **AND** reassignment is logged in audit trail

#### Scenario: Unassign alert
- **WHEN** alert is unassigned
- **THEN** assigned_to is set to null
- **AND** alert returns to unassigned queue
- **AND** unassignment reason is recorded

### Requirement: Alert Filtering and Search
The system SHALL provide comprehensive filtering and search capabilities for alerts.

#### Scenario: Filter by severity
- **WHEN** filtering alerts by severity
- **THEN** only alerts matching selected severities are displayed
- **AND** filter supports multi-select (e.g., HIGH and CRITICAL)
- **AND** severity distribution is shown in filter UI

#### Scenario: Filter by status
- **WHEN** filtering alerts by status
- **THEN** alerts in selected statuses are displayed
- **AND** default view shows NEW and TRIAGED alerts only
- **AND** status counts are displayed

#### Scenario: Filter by assigned analyst
- **WHEN** filtering by analyst
- **THEN** alerts assigned to selected analysts are displayed
- **AND** "My Alerts" shortcut shows current user's assignments
- **AND** "Unassigned" filter shows alerts with no owner

#### Scenario: Filter by time range
- **WHEN** filtering by date range
- **THEN** alerts created within range are displayed
- **AND** common ranges (today, last 7 days, last 30 days) are provided
- **AND** custom date range can be specified

#### Scenario: Filter by alert type
- **WHEN** filtering by alert type
- **THEN** alerts of selected types are displayed
- **AND** type distribution is shown
- **AND** multiple types can be selected

#### Scenario: Combined filters
- **WHEN** applying multiple filters simultaneously
- **THEN** alerts matching ALL filter criteria are displayed (AND logic)
- **AND** filter combination is saved as view
- **AND** active filters are clearly displayed with ability to remove individual filters

### Requirement: Alert Dashboard
The system SHALL provide dashboard showing alert metrics and operational insights.

#### Scenario: View alert summary
- **WHEN** accessing alert dashboard
- **THEN** total alert count is displayed
- **AND** alert count by severity is shown
- **AND** alert count by status is shown
- **AND** trend line shows alerts over time

#### Scenario: View alert backlog
- **WHEN** viewing backlog metrics
- **THEN** number of NEW alerts is displayed
- **AND** number of TRIAGED alerts pending investigation is shown
- **AND** oldest alert age is highlighted
- **AND** backlog trend over 30 days is charted

#### Scenario: View alert resolution metrics
- **WHEN** viewing resolution metrics
- **THEN** average time-to-resolution is displayed
- **AND** resolution distribution (confirmed/FP/no-action) is shown
- **AND** analyst productivity metrics are displayed

#### Scenario: Alert dashboard performance
- **WHEN** loading alert dashboard
- **THEN** dashboard loads within 500ms
- **AND** metrics are cached with 1-minute TTL
- **AND** dashboard auto-refreshes every 30 seconds

### Requirement: Alert Notifications
The system SHALL notify analysts of new alerts and status changes through multiple channels.

#### Scenario: Notify on new high-severity alert
- **WHEN** CRITICAL or HIGH severity alert is generated
- **THEN** assigned analyst receives immediate notification
- **AND** notification includes alert ID, type, severity, and summary
- **AND** notification includes link to alert detail

#### Scenario: Notify on alert assignment
- **WHEN** alert is assigned to analyst
- **THEN** analyst receives notification
- **AND** notification is delivered via configured channels (email, in-app, Slack)
- **AND** notification can be dismissed

#### Scenario: Batch notifications for low-severity alerts
- **WHEN** LOW or MEDIUM severity alerts are generated
- **THEN** notifications are batched and sent every hour
- **AND** batch includes summary of all new alerts
- **AND** analyst can opt out of batch notifications

#### Scenario: Notify on approaching SLA breach
- **WHEN** alert is approaching SLA deadline
- **THEN** assigned analyst receives warning notification
- **AND** notification is sent 2 hours and 30 minutes before breach
- **AND** supervisor is notified 15 minutes before breach

### Requirement: Alert API
The system SHALL provide REST API for alert management and integration with external systems.

#### Scenario: List alerts via API
- **WHEN** external system queries alerts via API
- **THEN** paginated list of alerts is returned
- **AND** filters can be applied via query parameters
- **AND** response includes alert metadata and links to related resources

#### Scenario: Get alert details via API
- **WHEN** querying specific alert by ID
- **THEN** full alert details are returned including all related entities and transactions
- **AND** response includes status history and assignment history
- **AND** 404 returned if alert not found

#### Scenario: Update alert via API
- **WHEN** external system updates alert status or assignment
- **THEN** update is validated and applied
- **AND** audit trail is recorded with API client ID
- **AND** webhooks are triggered for status changes

#### Scenario: Create alert via API
- **WHEN** external system creates alert
- **THEN** alert is validated and stored
- **AND** alert follows same workflow as system-generated alerts
- **AND** 201 Created returned with alert ID and location header

### Requirement: Alert Audit Trail
The system SHALL maintain complete audit trail of all alert actions and changes.

#### Scenario: Log alert creation
- **WHEN** alert is created
- **THEN** creation event is logged with timestamp and source
- **AND** initial severity and status are recorded

#### Scenario: Log status changes
- **WHEN** alert status changes
- **THEN** change is logged with old status, new status, user_id, and timestamp
- **AND** reason for change is captured if provided

#### Scenario: Log assignment changes
- **WHEN** alert is assigned or reassigned
- **THEN** assignment change is logged with previous and new assignee
- **AND** timestamp and acting user are recorded

#### Scenario: View alert history
- **WHEN** viewing alert detail
- **THEN** complete history of all changes is displayed as timeline
- **AND** timeline includes all status changes, assignments, and notes
- **AND** timeline is sorted chronologically with most recent first

