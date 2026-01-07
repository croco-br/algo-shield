## ADDED Requirements

### Requirement: Case Creation
The system SHALL automatically create cases from high-severity alerts and support manual case creation.

#### Scenario: Auto-create case from critical alert
- **WHEN** alert with CRITICAL or HIGH severity is generated
- **THEN** case is automatically created and linked to alert
- **AND** case priority is set based on alert severity
- **AND** case SLA deadline is calculated from creation time and priority

#### Scenario: Create case manually
- **WHEN** analyst manually creates case
- **THEN** case is created with analyst as creator
- **AND** case can be linked to one or more alerts
- **AND** case can be created without associated alert (investigation-initiated)

#### Scenario: Set case priority
- **WHEN** creating case
- **THEN** priority is set to LOW, MEDIUM, HIGH, or CRITICAL
- **AND** priority determines SLA deadline and assignment routing
- **AND** priority can be manually adjusted with justification

### Requirement: Case Assignment
The system SHALL support manual and automatic assignment of cases to investigators with workload balancing.

#### Scenario: Auto-assign case
- **WHEN** case is created with auto-assignment enabled
- **THEN** case is assigned using round-robin or load-based algorithm
- **AND** analyst availability and skills are considered
- **AND** assigned analyst receives notification

#### Scenario: Manual assignment
- **WHEN** supervisor manually assigns case
- **THEN** case assigned_to field is updated
- **AND** assignment reason can be provided
- **AND** previous assignee (if any) is notified of reassignment

#### Scenario: Workload balancing
- **WHEN** auto-assigning cases
- **THEN** current workload of each analyst is calculated
- **AND** workload considers case count and weighted priority
- **AND** cases are assigned to maintain balanced workload

#### Scenario: Skill-based routing
- **WHEN** case requires specific expertise
- **THEN** case is routed to analysts with matching skills
- **AND** skills are configured per analyst profile
- **AND** skill match is considered alongside workload

### Requirement: Case Status Workflow
The system SHALL enforce case status workflow with valid state transitions and required actions.

#### Scenario: Open new case
- **WHEN** case is created
- **THEN** case status is set to OPEN
- **AND** case appears in open cases queue
- **AND** SLA timer starts

#### Scenario: Begin investigation
- **WHEN** analyst begins investigating case
- **THEN** case status changes to INVESTIGATING
- **AND** first investigation note should be added
- **AND** status change is logged in audit trail

#### Scenario: Escalate case
- **WHEN** case requires escalation
- **THEN** case status changes to ESCALATED
- **AND** escalation reason is required
- **AND** case is reassigned to senior analyst or supervisor
- **AND** escalation notification is sent

#### Scenario: Close case
- **WHEN** investigation is complete
- **THEN** case status changes to CLOSED
- **AND** resolution type is required (confirmed_fraud, false_positive, requires_reporting, no_action)
- **AND** closing notes are required
- **AND** closed_at timestamp is recorded

#### Scenario: Reopen closed case
- **WHEN** new information emerges about closed case
- **THEN** case can be reopened by supervisor
- **AND** reopen reason is required
- **AND** reopen is logged in audit trail
- **AND** case SLA is recalculated

### Requirement: Case Notes and Collaboration
The system SHALL provide collaborative workspace for case investigation with notes and attachments.

#### Scenario: Add investigation note
- **WHEN** investigator adds note to case
- **THEN** note is stored with author_id and timestamp
- **AND** note content supports rich text formatting
- **AND** note is visible to all assigned investigators

#### Scenario: Mention colleagues in notes
- **WHEN** note mentions other user (@username)
- **THEN** mentioned user receives notification
- **AND** mentioned user can quickly navigate to case
- **AND** mention creates implicit collaboration link

#### Scenario: Attach evidence to case
- **WHEN** investigator uploads evidence file
- **THEN** file is stored securely with case association
- **AND** file type and size are validated
- **AND** file access is restricted to authorized users

#### Scenario: View note history
- **WHEN** viewing case notes
- **THEN** notes are displayed chronologically
- **AND** note author and timestamp are shown
- **AND** edited notes show edit history
- **AND** deleted notes are marked but retained

### Requirement: Case SLA Tracking
The system SHALL track SLA compliance for cases with deadline enforcement and breach notifications.

#### Scenario: Calculate SLA deadline
- **WHEN** case is created
- **THEN** SLA deadline is calculated based on priority
- **AND** deadlines: CRITICAL=2h, HIGH=8h, MEDIUM=24h, LOW=72h
- **AND** deadline is adjusted for business hours (configurable)

#### Scenario: Warn before SLA breach
- **WHEN** case is approaching SLA deadline
- **THEN** warning notification is sent to assigned investigator (2 hours before)
- **AND** second warning is sent to supervisor (30 minutes before)
- **AND** warnings are escalated if case still open

#### Scenario: Handle SLA breach
- **WHEN** case exceeds SLA deadline without resolution
- **THEN** breach is recorded and flagged
- **AND** breach notification is sent to management
- **AND** case is highlighted in dashboards
- **AND** breach does not prevent case closure

#### Scenario: Pause SLA for external dependencies
- **WHEN** case is waiting on external information
- **THEN** investigator can pause SLA timer
- **AND** pause reason is required
- **AND** pause duration is tracked separately
- **AND** SLA can be resumed when information received

### Requirement: Case Search and Filtering
The system SHALL provide comprehensive search and filtering for cases.

#### Scenario: Search by case ID or keywords
- **WHEN** searching cases by text
- **THEN** search includes case ID, entity IDs, and note content
- **AND** search results are ranked by relevance
- **AND** search supports partial matches

#### Scenario: Filter by status
- **WHEN** filtering cases by status
- **THEN** cases in selected statuses are displayed
- **AND** default view shows OPEN and INVESTIGATING cases
- **AND** status counts are shown in filter UI

#### Scenario: Filter by priority
- **WHEN** filtering by priority
- **THEN** cases with selected priorities are displayed
- **AND** HIGH and CRITICAL cases are prioritized in views

#### Scenario: Filter by assignee
- **WHEN** filtering by assigned investigator
- **THEN** cases assigned to selected investigators are shown
- **AND** "My Cases" view shows current user's assignments
- **AND** "Unassigned" filter shows cases needing assignment

#### Scenario: Filter by SLA status
- **WHEN** filtering by SLA status
- **THEN** cases can be filtered by: within_sla, at_risk, breached
- **AND** at_risk shows cases with <25% time remaining
- **AND** breached shows cases past deadline

### Requirement: Case Resolution and Feedback
The system SHALL capture case resolution with outcomes that feed back into fraud detection system.

#### Scenario: Resolve as confirmed fraud
- **WHEN** case is closed with resolution "confirmed_fraud"
- **THEN** resolution confirms fraud was correctly detected
- **AND** positive feedback is recorded for triggered rules
- **AND** entity risk scores are updated

#### Scenario: Resolve as false positive
- **WHEN** case is closed with resolution "false_positive"
- **THEN** resolution indicates incorrect fraud detection
- **AND** negative feedback is recorded for triggered rules
- **AND** false positive is included in rule effectiveness metrics

#### Scenario: Resolve as requiring regulatory reporting
- **WHEN** case is closed with resolution "requires_reporting"
- **THEN** case is flagged for SAR generation
- **AND** case appears in regulatory reporting queue
- **AND** case evidence is compiled for report

#### Scenario: Resolve with no action
- **WHEN** case is closed with resolution "no_action"
- **THEN** case indicates investigation found no fraud
- **AND** neutral feedback is recorded
- **AND** case is available for future reference

### Requirement: Case Audit Trail
The system SHALL maintain complete audit trail of all case actions and changes for regulatory compliance.

#### Scenario: Log case creation
- **WHEN** case is created
- **THEN** creation event is logged with creator, timestamp, and initial state

#### Scenario: Log all status changes
- **WHEN** case status changes
- **THEN** change is logged with old status, new status, user_id, reason, and timestamp
- **AND** state transition validation is enforced

#### Scenario: Log assignment changes
- **WHEN** case is assigned or reassigned
- **THEN** assignment change is logged with previous and new assignee, reason, and timestamp

#### Scenario: Log all notes and actions
- **WHEN** investigator adds notes or takes actions
- **THEN** all actions are logged with user_id and timestamp
- **AND** note edits and deletions are logged

#### Scenario: Generate audit report
- **WHEN** compliance officer generates audit report for case
- **THEN** complete timeline of all actions is provided
- **AND** report includes all participants, decisions, and evidence
- **AND** report is exportable in PDF format

### Requirement: Case Dashboard
The system SHALL provide dashboard showing case metrics and operational insights.

#### Scenario: View open cases summary
- **WHEN** accessing case dashboard
- **THEN** total open cases count is displayed
- **AND** open cases by priority are shown
- **AND** cases by status distribution is charted

#### Scenario: View SLA compliance
- **WHEN** viewing SLA metrics
- **THEN** overall SLA compliance percentage is displayed
- **AND** cases at risk of SLA breach are highlighted
- **AND** breached cases count is shown

#### Scenario: View team workload
- **WHEN** viewing team metrics
- **THEN** case count per investigator is displayed
- **AND** weighted workload (considering priority) is shown
- **AND** workload imbalance is highlighted

#### Scenario: View resolution metrics
- **WHEN** viewing resolution statistics
- **THEN** average time-to-resolution by priority is displayed
- **AND** resolution distribution (confirmed/FP/reporting/no-action) is shown
- **AND** case reopening rate is displayed

