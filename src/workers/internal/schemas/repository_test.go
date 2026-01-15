package schemas

// Note: The PostgresRepository methods (GetByID, ListAll) are thin wrappers around
// PostgreSQL queries with JSON unmarshaling. These are best tested with integration
// tests using real database connections (testcontainers-go) rather than complex
// mocking of pgxpool.Pool which provides limited value.
//
// Integration tests for this repository should cover:
// - GetByID with valid UUID
// - GetByID with invalid/non-existent UUID
// - ListAll returning multiple schemas
// - ListAll with empty result
// - JSON unmarshaling of sample_json and extracted_fields
// - Database connection errors
//
// See docs/agents/integration-test.md for integration testing guidelines.
