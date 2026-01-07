# Design: Adaptive Fraud Detection Platform

## Context
AlgoShield is currently a high-performance rule-based fraud engine optimized for deterministic decisions and low latency (<50ms). This design transforms it into an adaptive fraud detection platform that can detect structural fraud patterns, learn from feedback, and provide comprehensive AML compliance. The challenge is adding these capabilities while maintaining performance and architectural clarity.

## Goals / Non-Goals

### Goals
- Enable detection of network-level fraud (account rings, coordinated behavior)
- Create feedback loop from analyst decisions to improve rules over time
- Monitor signal health and detect distribution drift
- Generate structured alerts with intelligent severity and grouping
- Provide business-aware metrics (cost, risk, operational impact)
- Support complete investigation workflow with audit trail
- Enable regulatory reporting and compliance tracking
- Maintain <50ms transaction processing latency
- Keep detection and compliance concerns architecturally separated

### Non-Goals
- Real-time ML model training (adaptive rules, not ML)
- Full customer onboarding/KYC system (integration point only)
- Document management system
- Multi-tenant isolation (single deployment)
- Real-time graph visualization in UI (batch/cached views acceptable)
- Advanced graph analytics (keep queries simple and fast)

## Architectural Decisions

### Decision 1: Entity-Relationship Model

**What**: Introduce entity model with relationships graph.

**Entities**:
- `Customer` - individual or business entity
- `Account` - financial account linked to customer
- `Device` - device fingerprint or identifier
- `IPAddress` - IP address with geolocation
- `Session` - user session with metadata

**Relationships**:
- `Customer` ↔ `Account` (owns, many-to-many)
- `Customer` ↔ `Device` (uses, many-to-many)
- `Account` ↔ `Transaction` (originates, one-to-many)
- `Device` ↔ `Session` (creates, one-to-many)
- `IPAddress` ↔ `Session` (from, one-to-many)

**Why**:
- Most sophisticated fraud emerges at relationship level
- Account rings, device sharing, coordinated behavior require graph
- AML patterns (structuring, smurfing) require customer-level aggregation

**Alternatives Considered**:
- Keep transaction-only model → rejected, cannot detect structural fraud
- Use full graph database (Neo4j) → deferred, start with PostgreSQL + indexes
- Infer relationships on-the-fly → rejected, too slow for real-time detection

**Implementation**:
- PostgreSQL tables with foreign keys
- Adjacency list representation
- Indexes on relationship queries
- Consider graph database if performance insufficient

---

### Decision 2: Feedback Loop Architecture

**What**: Capture analyst decisions and calculate rule effectiveness metrics.

**Components**:
- `feedback` table - stores analyst decisions (approve, reject, false_positive)
- `rule_metrics` table - aggregated effectiveness metrics per rule
- `adaptive_rules` - rules with dynamic thresholds based on feedback
- Background job - recalculates metrics daily

**Metrics Calculated**:
- False positive rate per rule
- True positive rate per rule
- Precision and recall
- Alert-to-confirmation ratio
- Time-to-resolution

**Why**:
- Static rules become stale as fraud adapts
- Need quantitative measure of rule quality
- Analysts know which rules are noisy but system doesn't track it
- Enable data-driven rule tuning

**Alternatives Considered**:
- Manual rule adjustment → rejected, doesn't scale
- Real-time ML → rejected, too complex for MVP
- External analytics → rejected, feedback should be in-system

**Implementation**:
- API endpoint for feedback submission
- Background worker for metric aggregation
- Redis cache for real-time metrics access
- UI dashboard showing rule effectiveness

---

### Decision 3: Drift Detection Strategy

**What**: Monitor statistical properties of input signals and alert on drift.

**Monitored Signals**:
- Transaction amount distribution (mean, stddev, percentiles)
- Transaction frequency per customer
- Geographic distribution
- Device diversity
- Rule trigger rates

**Detection Methods**:
- Kolmogorov-Smirnov test for distribution changes
- Moving average for trend detection
- Percentile-based anomaly detection
- Window comparison (current week vs baseline)

**Why**:
- Input distribution changes indicate fraud adaptation or data quality issues
- Silent degradation is common in production fraud systems
- Early warning enables proactive response

**Alternatives Considered**:
- No monitoring → rejected, systems degrade silently
- ML-based anomaly detection → deferred, statistical methods sufficient for MVP
- Real-time drift detection → rejected, batch processing sufficient

**Implementation**:
- Daily batch job calculates distribution statistics
- Compare with 30-day baseline
- Alert if KS test p-value < 0.05 or percentile shift > 20%
- Store metrics in `signal_health` table
- Dashboard shows trend lines

---

### Decision 4: Alert System Design

**What**: Transform transaction detections into structured alerts with severity.

**Alert Model**:
```
Alert {
  id: UUID
  severity: LOW | MEDIUM | HIGH | CRITICAL
  type: string (e.g., "high_value", "velocity", "ring_detected")
  status: NEW | TRIAGED | INVESTIGATING | CLOSED
  related_transactions: []UUID
  related_entities: []EntityRef
  triggered_rules: []string
  assigned_to: nullable UUID (analyst)
  created_at, updated_at
}
```

**Severity Assignment**:
- Rule priority + amount threshold = base severity
- Entity risk score modifies severity
- Historical pattern (first-time vs repeat) modifies severity
- Multi-rule triggers escalate severity

**Grouping Strategy**:
- Group alerts for same entity within time window
- Group alerts for related entities (graph-connected)
- Single alert for coordinated behavior across multiple entities

**Why**:
- Transaction-level detections create alert fatigue
- Severity helps prioritize analyst attention
- Grouping reduces noise and surfaces patterns
- Status tracking enables workflow management

**Alternatives Considered**:
- Keep transaction-only model → rejected, too noisy
- Generate alerts in UI only → rejected, alerts are domain objects
- Simple priority score → rejected, multi-factor severity more useful

**Implementation**:
- Worker generates alerts after transaction processing
- Alert service handles grouping and severity calculation
- Redis pub/sub notifies UI of new alerts
- API endpoints for alert management (assign, triage, close)

---

### Decision 5: Sociotechnical Metrics

**What**: Calculate business-aware metrics beyond technical accuracy.

**Metrics**:
- **Cost of False Positives**: blocked_transactions * avg_transaction_value * customer_churn_rate
- **Risk Exposure**: pending_review_transactions * amount * fraud_probability
- **Operational Load**: alerts_per_analyst, avg_investigation_time, backlog_size
- **SLA Compliance**: cases_within_sla / total_cases
- **Customer Impact**: friction_rate, false_decline_rate, support_ticket_rate

**Why**:
- Fraud detection is not just accuracy optimization
- False positives have real business cost
- Analysts need to prioritize based on risk and impact
- Compliance requires SLA tracking

**Alternatives Considered**:
- Technical metrics only → rejected, doesn't enable business decisions
- Manual calculation → rejected, should be automated
- Complex cost modeling → deferred, start with simple formulas

**Implementation**:
- Metrics calculation in dedicated service
- Configurable cost parameters (churn rate, fraud probability)
- Dashboard visualization with trends
- API endpoints for metric queries

---

### Decision 6: Event Correlation Architecture

**What**: Extend velocity checks to temporal correlation and behavior modeling.

**Current State**:
- `velocityCount(origin, time_window)` - count transactions
- `velocitySum(origin, time_window)` - sum amounts

**Enhancements**:
- `behaviorDeviation(entity_id, metric)` - compare to baseline
- `sequenceMatch(entity_id, pattern)` - detect specific sequences
- `correlationScore(entity_ids)` - measure coordinated behavior
- `temporalCluster(entity_id, time_window)` - detect bursts

**Why**:
- Single-event rules miss temporal patterns
- Baseline deviation is core AML requirement
- Coordinated behavior requires multi-entity correlation

**Alternatives Considered**:
- Keep simple velocity only → rejected, insufficient for AML
- Stream processing (Flink, Kafka Streams) → deferred, batch acceptable for MVP
- Complex sequence mining → rejected, keep patterns simple

**Implementation**:
- Pre-calculate baselines nightly (30-day window)
- Store in `entity_baselines` table
- Rule expressions can call new correlation functions
- Worker queries baselines during evaluation

---

### Decision 7: Case Management System

**What**: Full investigation workflow from alert to resolution.

**Case Model**:
```
Case {
  id: UUID
  alert_id: UUID
  status: OPEN | INVESTIGATING | ESCALATED | CLOSED | FALSE_POSITIVE
  priority: LOW | MEDIUM | HIGH | CRITICAL
  assigned_to: UUID (analyst)
  created_by: UUID
  resolution: CONFIRMED_FRAUD | FALSE_POSITIVE | REQUIRES_REPORTING | NO_ACTION
  notes: []CaseNote
  audit_trail: []AuditEntry
  sla_deadline: timestamp
  created_at, updated_at, closed_at
}

CaseNote {
  id: UUID
  case_id: UUID
  author_id: UUID
  content: text
  created_at
}

AuditEntry {
  id: UUID
  case_id: UUID
  user_id: UUID
  action: string (e.g., "status_changed", "assigned", "note_added")
  old_value, new_value: JSON
  timestamp
}
```

**Workflow**:
1. Alert generated → Case created automatically for HIGH/CRITICAL
2. Case assigned to analyst (manual or auto-distribution)
3. Analyst investigates → adds notes, requests data
4. Decision made → case closed with resolution
5. Audit trail captured at every step

**Why**:
- Alerts need investigation workflow
- Compliance requires complete audit trail
- Analysts need collaboration tools (notes, assignment)
- SLA tracking requires case model

**Alternatives Considered**:
- No case management → rejected, compliance requirement
- Use external case system → possible, but integration overhead
- Simple status on alerts → rejected, insufficient for audit trail

**Implementation**:
- Dedicated case management service
- API endpoints for CRUD and workflow actions
- UI for case list, detail, assignment
- Audit trail captured via middleware

---

### Decision 8: Regulatory Reporting

**What**: Generate Suspicious Activity Reports (SARs) and compliance exports.

**Components**:
- `reportable_cases` - cases flagged for regulatory reporting
- `sar_generator` - template-based SAR generation
- `compliance_exports` - CSV/PDF exports for auditors
- `regulatory_obligations` - tracking of reporting deadlines

**SAR Generation**:
- Template with required fields per jurisdiction
- Auto-populate from case data
- Manual review and edit before submission
- Track submission status

**Compliance Exports**:
- CSV: transactions, alerts, cases, decisions
- PDF: investigation summary, timeline, supporting evidence
- Filters: date range, entity, type, jurisdiction

**Why**:
- AML compliance requires SAR submission
- Audits require exportable evidence
- Manual SAR creation is error-prone

**Alternatives Considered**:
- Manual reporting → rejected, error-prone and slow
- External reporting system → possible, provide API
- Full regulatory platform → rejected, out of scope

**Implementation**:
- Reporting service with template engine
- API endpoints for SAR CRUD
- Export endpoints with streaming for large datasets
- UI for SAR creation and submission tracking

---

## Database Schema Additions

### New Tables

```sql
-- Entities
CREATE TABLE entities (
  id UUID PRIMARY KEY,
  entity_type VARCHAR(50) NOT NULL, -- 'customer', 'account', 'device', 'ip', 'session'
  external_id VARCHAR(255),
  metadata JSONB,
  risk_score INTEGER DEFAULT 0,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);

CREATE TABLE entity_relationships (
  id UUID PRIMARY KEY,
  entity_a_id UUID REFERENCES entities(id),
  entity_b_id UUID REFERENCES entities(id),
  relationship_type VARCHAR(50), -- 'owns', 'uses', 'shares', etc.
  strength FLOAT, -- 0.0 to 1.0
  created_at TIMESTAMPTZ
);

-- Feedback
CREATE TABLE feedback (
  id UUID PRIMARY KEY,
  transaction_id UUID REFERENCES transactions(id),
  alert_id UUID,
  decision VARCHAR(50), -- 'approve', 'reject', 'false_positive'
  analyst_id UUID REFERENCES users(id),
  reason TEXT,
  created_at TIMESTAMPTZ
);

CREATE TABLE rule_metrics (
  id UUID PRIMARY KEY,
  rule_id UUID REFERENCES rules(id),
  date DATE,
  total_triggers INTEGER,
  true_positives INTEGER,
  false_positives INTEGER,
  precision FLOAT,
  recall FLOAT,
  calculated_at TIMESTAMPTZ
);

-- Drift Detection
CREATE TABLE signal_health (
  id UUID PRIMARY KEY,
  signal_name VARCHAR(100),
  date DATE,
  mean FLOAT,
  stddev FLOAT,
  p50 FLOAT,
  p95 FLOAT,
  p99 FLOAT,
  ks_statistic FLOAT,
  drift_detected BOOLEAN,
  calculated_at TIMESTAMPTZ
);

-- Alerts
CREATE TABLE alerts (
  id UUID PRIMARY KEY,
  severity VARCHAR(20),
  type VARCHAR(100),
  status VARCHAR(50),
  related_transaction_ids UUID[],
  related_entity_ids UUID[],
  triggered_rules TEXT[],
  assigned_to UUID REFERENCES users(id),
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);

-- Cases
CREATE TABLE cases (
  id UUID PRIMARY KEY,
  alert_id UUID REFERENCES alerts(id),
  status VARCHAR(50),
  priority VARCHAR(20),
  assigned_to UUID REFERENCES users(id),
  created_by UUID REFERENCES users(id),
  resolution VARCHAR(50),
  sla_deadline TIMESTAMPTZ,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ,
  closed_at TIMESTAMPTZ
);

CREATE TABLE case_notes (
  id UUID PRIMARY KEY,
  case_id UUID REFERENCES cases(id),
  author_id UUID REFERENCES users(id),
  content TEXT,
  created_at TIMESTAMPTZ
);

CREATE TABLE audit_trail (
  id UUID PRIMARY KEY,
  case_id UUID REFERENCES cases(id),
  user_id UUID REFERENCES users(id),
  action VARCHAR(100),
  old_value JSONB,
  new_value JSONB,
  created_at TIMESTAMPTZ
);

-- Baselines
CREATE TABLE entity_baselines (
  id UUID PRIMARY KEY,
  entity_id UUID REFERENCES entities(id),
  metric VARCHAR(100),
  window_days INTEGER,
  mean FLOAT,
  stddev FLOAT,
  calculated_at TIMESTAMPTZ
);
```

### Modified Tables

```sql
ALTER TABLE transactions ADD COLUMN entity_id UUID REFERENCES entities(id);
ALTER TABLE transactions ADD COLUMN customer_id UUID REFERENCES entities(id);
CREATE INDEX idx_transactions_entity_id ON transactions(entity_id);
CREATE INDEX idx_transactions_customer_id ON transactions(customer_id);
```

---

## Performance Considerations

### Graph Query Optimization
- Index on `entity_relationships(entity_a_id, entity_b_id)`
- Index on `entity_relationships(relationship_type)`
- Limit graph traversal depth (max 3 hops)
- Cache frequently queried subgraphs in Redis
- Consider materialized views for common patterns

### Real-Time Processing
- Entity resolution happens async after transaction decision
- Alert generation happens async
- Case creation happens async
- Keep transaction processing path unchanged (<50ms)

### Dashboard Performance
- Pre-aggregate metrics in background jobs
- Cache dashboard data in Redis (30s TTL)
- Use indexed queries for all filters
- Pagination for all lists
- Streaming for large exports

---

## Migration Strategy

### Phase 1 Migration: Entities
1. Create entity tables
2. Backfill entities from existing transactions
   - Extract unique origins/destinations → accounts
   - Infer customers from account patterns
   - Extract device fingerprints from metadata
3. Link transactions to entities
4. Validate entity graph integrity

### Phase 2 Migration: Feedback & Metrics
1. Create feedback tables
2. No backfill needed (start fresh)
3. Begin capturing analyst decisions
4. Calculate baseline metrics after 30 days

### Phase 3 Migration: Alerts & Cases
1. Create alert and case tables
2. Optionally migrate in_review transactions → alerts
3. No case backfill (historical data stays in transactions)

### Rollback Plan
- Each phase is additive (no destructive changes)
- Can disable features via feature flags
- Can drop new tables if rollback needed
- Transaction processing path unchanged (safe rollback)

---

## Security Considerations

- **Audit Trail**: All case actions logged with user_id and timestamp
- **RBAC**: New roles - `analyst`, `senior_analyst`, `compliance_officer`
- **Data Privacy**: Entity PII encrypted at rest
- **Access Control**: Entity data access logged
- **Export Controls**: Only compliance officers can export full datasets

---

## Open Questions

1. **Graph Database**: Should we use Neo4j or stick with PostgreSQL?
   - **Decision**: Start with PostgreSQL, migrate if performance insufficient
   
2. **Real-time vs Batch**: Should correlation be real-time or batch?
   - **Decision**: Batch for MVP (nightly baseline calculation), real-time later if needed

3. **ML Integration**: When should we add ML models?
   - **Decision**: After feedback loop has 6+ months of data, use for scoring not rule replacement

4. **Separate Service**: Should Phase 3 (compliance) be separate microservice?
   - **Decision**: Keep in monolith initially, extract if complexity grows

5. **Graph Visualization**: What library for entity graph UI?
   - **Decision**: D3.js for simple visualizations, defer complex interactive graphs

