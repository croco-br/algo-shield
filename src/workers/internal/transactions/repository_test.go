package transactions

// Note: The PostgresRepository methods (SaveTransaction, SaveSyntheticTransaction)
// are thin wrappers around PostgreSQL INSERT queries with JSON marshaling. These are
// best tested with integration tests using real database connections (testcontainers-go)
// rather than complex mocking of pgxpool.Pool which provides limited value.
//
// Integration tests for this repository should cover:
// - SaveTransaction with complete transaction data
// - SaveSyntheticTransaction with complete transaction data
// - JSON marshaling of matched_rules and metadata fields
// - Database constraint violations (e.g., duplicate IDs)
// - Database connection errors
// - NULL handling for optional fields (schema_id, processed_at)
//
// The TransactionHistoryRepository methods are tested for SQL injection prevention
// in history_repository_test.go. The actual query execution requires integration tests.
//
// See docs/agents/integration-test.md for integration testing guidelines.
