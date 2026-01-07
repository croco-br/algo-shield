# Change: Transform AlgoShield into Adaptive Fraud Detection Platform

## Why
AlgoShield is currently a high-quality rule-based fraud engine, but it lacks critical capabilities needed for modern fraud detection: entity relationship modeling, adaptive learning from human feedback, signal health monitoring, and comprehensive AML compliance features. The current static rule approach limits the system's ability to detect sophisticated fraud patterns (account rings, coordinated behavior) and adapt to evolving fraud tactics. This change transforms AlgoShield from a rule engine into an adaptive fraud detection platform that can detect structural fraud patterns, learn from analyst decisions, and provide comprehensive AML compliance capabilities.

## What Changes
- **ADDED**: Entity Graph Layer - model relationships between customers, accounts, devices, IPs to detect coordinated fraud
- **ADDED**: Feedback Loop System - capture analyst decisions, calculate rule effectiveness, enable adaptive rules
- **ADDED**: Drift Detection & Signal Health - monitor input distributions, detect statistical drift, alert on signal degradation
- **ADDED**: Alert System - structured alerts with severity, triaging, and false positive tracking
- **ADDED**: Sociotechnical Metrics - business cost of false positives, risk exposure, operational impact, SLA tracking
- **ADDED**: Advanced Event Correlation - temporal correlation, behavior modeling, sequence detection
- **ADDED**: Case Management System - investigation workflow, analyst assignment, case tracking, audit trail
- **ADDED**: Regulatory Reporting - SAR generation, compliance reporting, regulatory obligation tracking
- **MODIFIED**: Transaction model - add entity_id, customer_id to link transactions to graph entities
- **MODIFIED**: Dashboard - separate Detection Dashboard (system health) from Compliance Dashboard (investigation)
- **MODIFIED**: Rules engine - enable adaptive thresholds based on feedback
- **MODIFIED**: Database schema - add tables for entities, relationships, feedback, alerts, cases

## Impact
- **Affected specs**: 
  - `entity-graph` (new capability)
  - `feedback-loop` (new capability)
  - `drift-detection` (new capability)
  - `alert-system` (new capability)
  - `sociotechnical-metrics` (new capability)
  - `event-correlation` (new capability)
  - `case-management` (new capability)
  - `regulatory-reporting` (new capability)
  - `transactions` (enhanced with entity relationships)
  - `dashboard` (split into detection and compliance views)
  - `rules` (enhanced with adaptive capabilities)

- **Affected code**:
  - `src/api/internal/entities/` - new entity management service
  - `src/api/internal/graph/` - new graph query service
  - `src/api/internal/feedback/` - new feedback tracking service
  - `src/api/internal/drift/` - new drift detection service
  - `src/api/internal/alerts/` - new alert management service
  - `src/api/internal/metrics/` - new sociotechnical metrics service
  - `src/api/internal/cases/` - new case management service
  - `src/api/internal/reporting/` - new regulatory reporting service
  - `src/api/internal/transactions/` - enhanced with entity links
  - `src/workers/internal/rules/` - enhanced with adaptive thresholds
  - `src/workers/internal/graph/` - new graph-aware rule evaluation
  - `src/ui/src/views/DetectionDashboard.vue` - new detection monitoring view
  - `src/ui/src/views/ComplianceDashboard.vue` - new compliance view
  - `src/ui/src/views/AlertsView.vue` - new alerts management view
  - `src/ui/src/views/CasesView.vue` - new case management view
  - `src/ui/src/views/EntityGraphView.vue` - new graph visualization
  - Database migrations - comprehensive schema additions

- **Architecture considerations**:
  - **BREAKING**: Major architectural shift from transaction-only to entity-relationship model
  - **Performance**: Graph queries require careful indexing and query optimization
  - **Scalability**: Entity graph may require graph database (Neo4j) or GraphQL layer
  - **Complexity**: Significant increase in system complexity across all layers
  - **Migration**: Requires backfilling entity data from existing transactions
  - **Integration**: New APIs for external case management and reporting systems

## Implementation Phases

### Phase 1: Adaptive Detection Layer (6-8 weeks)
Core detection capabilities that make AlgoShield adaptive and graph-aware.

### Phase 2: Orchestration & Enrichment Layer (4-6 weeks)
Alert generation, sociotechnical metrics, and advanced correlation.

### Phase 3: AML Compliance Layer (8-12 weeks)
Full case management and regulatory reporting capabilities.

## Success Criteria
- Entity graph can model relationships and detect coordinated fraud
- Feedback loop reduces false positive rate by 20%+ within 3 months
- Drift detection alerts on signal degradation before accuracy drops
- Alert system reduces analyst noise through intelligent severity and grouping
- Sociotechnical metrics enable cost-aware decision making
- Case management provides complete audit trail for regulatory compliance
- Detection dashboard <500ms, Compliance dashboard <2s response times
- System maintains <50ms transaction processing latency despite added complexity

## Risks
- **Complexity**: Massive increase in system complexity may impact maintainability
- **Performance**: Graph queries and correlation may impact processing latency
- **Data Quality**: Entity resolution and relationship inference require high-quality data
- **Migration**: Backfilling entity graph from historical data may be challenging
- **Adoption**: Analysts need training on new workflows and capabilities
- **Scope Creep**: Phase 3 (compliance) may be better as separate system

