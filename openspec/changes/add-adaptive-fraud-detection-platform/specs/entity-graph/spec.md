## ADDED Requirements

### Requirement: Entity Model
The system SHALL provide an entity model that represents customers, accounts, devices, IPs, and sessions as first-class objects with unique identifiers and metadata.

#### Scenario: Create customer entity
- **WHEN** a new customer is identified from transaction data
- **THEN** a customer entity is created with unique ID
- **AND** entity metadata includes external identifiers and attributes
- **AND** entity is assigned an initial risk score of 0

#### Scenario: Create account entity
- **WHEN** a new account origin or destination is detected
- **THEN** an account entity is created
- **AND** account is linked to owning customer (if known)
- **AND** account entity includes account type and status

#### Scenario: Create device entity
- **WHEN** a new device fingerprint is detected in transaction metadata
- **THEN** a device entity is created
- **AND** device metadata includes fingerprint, type, and first-seen timestamp

### Requirement: Entity Relationships
The system SHALL model relationships between entities as typed edges with relationship strength and creation timestamp.

#### Scenario: Customer owns account
- **WHEN** a customer entity is linked to an account entity
- **THEN** an "owns" relationship is created
- **AND** relationship has strength value between 0.0 and 1.0
- **AND** relationship timestamp is recorded

#### Scenario: Device used by customer
- **WHEN** a device is used by a customer in a transaction
- **THEN** a "uses" relationship is created or updated
- **AND** relationship strength increases with repeated usage
- **AND** last-used timestamp is updated

#### Scenario: Accounts share device
- **WHEN** multiple accounts use the same device
- **THEN** "shares" relationships are created between accounts
- **AND** relationship strength reflects frequency of sharing

### Requirement: Graph Traversal
The system SHALL provide graph traversal capabilities to query connected entities up to a configurable maximum depth.

#### Scenario: Find related accounts
- **WHEN** querying for accounts related to a customer
- **THEN** all directly owned accounts are returned
- **AND** accounts sharing devices with customer's accounts are included
- **AND** traversal depth is limited to maximum 3 hops

#### Scenario: Detect account ring
- **WHEN** querying for potential account rings
- **THEN** accounts connected through shared devices or IPs are identified
- **AND** ring size (number of connected accounts) is calculated
- **AND** ring strength (average relationship strength) is calculated

#### Scenario: Graph query performance
- **WHEN** executing graph traversal queries
- **THEN** queries complete within 100ms for graphs with up to 1000 connected entities
- **AND** queries use indexed lookups for performance
- **AND** frequently accessed subgraphs are cached in Redis

### Requirement: Entity Resolution
The system SHALL resolve and match entities across transactions to maintain a consistent entity graph.

#### Scenario: Match existing entity
- **WHEN** processing a transaction with known account identifier
- **THEN** existing account entity is matched and reused
- **AND** entity metadata is updated if new information is available

#### Scenario: Merge duplicate entities
- **WHEN** two entities are determined to represent the same real-world entity
- **THEN** entities are merged into a single entity
- **AND** all relationships are transferred to the merged entity
- **AND** merge is recorded in audit trail

#### Scenario: Entity resolution during transaction processing
- **WHEN** worker processes a transaction
- **THEN** entities are resolved asynchronously (not blocking transaction decision)
- **AND** entity resolution completes within 200ms
- **AND** transaction includes entity_id and customer_id after resolution

### Requirement: Graph Query Functions for Rules
The system SHALL provide graph query functions that can be used in rule expressions.

#### Scenario: Get related entities in rule
- **WHEN** a rule expression calls `getRelatedEntities(entity_id, "device")`
- **THEN** all device entities related to the specified entity are returned
- **AND** rule evaluation continues with returned entity list

#### Scenario: Detect ring in rule
- **WHEN** a rule expression calls `detectRing(account_id, 2)`
- **THEN** connected accounts within 2 hops are analyzed
- **AND** function returns true if ring size > 3 accounts
- **AND** ring detection completes within 50ms

#### Scenario: Count related entities in rule
- **WHEN** a rule expression calls `countRelated(customer_id, "account", "7d")`
- **THEN** number of accounts linked to customer in last 7 days is returned
- **AND** count includes only active relationships

### Requirement: Entity Risk Score
The system SHALL maintain a risk score for each entity that reflects historical fraud patterns and current risk indicators.

#### Scenario: Calculate entity risk score
- **WHEN** entity is involved in high-risk transactions
- **THEN** entity risk score increases proportionally
- **AND** risk score is capped at 100
- **AND** risk score decays over time (30-day half-life)

#### Scenario: Risk score influences alert severity
- **WHEN** generating alert for transaction involving high-risk entity
- **THEN** alert severity is elevated
- **AND** entity risk score is included in alert metadata

#### Scenario: Display entity risk in UI
- **WHEN** viewing entity details in UI
- **THEN** current risk score is displayed with visual indicator
- **AND** risk score history is shown as trend line
- **AND** contributing factors to risk score are listed

### Requirement: Entity Data Privacy
The system SHALL protect entity PII through encryption and access controls.

#### Scenario: Encrypt PII at rest
- **WHEN** entity data includes PII (names, addresses, etc.)
- **THEN** PII fields are encrypted in database
- **AND** encryption uses AES-256
- **AND** encryption keys are managed separately from application

#### Scenario: Access logging
- **WHEN** user accesses entity data
- **THEN** access is logged with user_id, entity_id, and timestamp
- **AND** access logs are retained for audit purposes
- **AND** unauthorized access attempts trigger alerts

