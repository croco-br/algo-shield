## ADDED Requirements

### Requirement: Synthetic Event Generation
The system SHALL provide the ability to generate synthetic transaction events from configured event schemas. Generated events SHALL have random values based on the field types defined in the schema.

#### Scenario: Generate events from schema
- **WHEN** an operator requests to generate N synthetic events from a configured schema
- **THEN** the system generates N events with random values based on each field's type
- **AND** each generated event is queued for processing as a transaction
- **AND** each transaction is linked to the source schema via schema_id

#### Scenario: Type-based randomization
- **WHEN** generating synthetic events
- **THEN** string fields receive random alphanumeric values
- **AND** number fields receive random numeric values in a reasonable range
- **AND** boolean fields receive random true/false values
- **AND** array fields receive 1-5 random elements based on array item type
- **AND** object fields are recursively generated with random nested values

#### Scenario: Event generation with seed
- **WHEN** an operator provides an optional seed value for event generation
- **THEN** the generated events are reproducible using the same seed
- **AND** the same seed with the same schema produces identical event sequences

#### Scenario: Generated events are processed
- **WHEN** synthetic events are generated
- **THEN** each event is immediately queued for worker processing
- **AND** the events are processed as normal transactions
- **AND** the schema_id is preserved through processing

