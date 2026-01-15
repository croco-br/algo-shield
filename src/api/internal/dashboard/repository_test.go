package dashboard

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

// Test_pgxIdentifier_Sanitize_WhenValidTableName_ThenQuotesCorrectly verifies that
// pgx.Identifier correctly sanitizes table names to prevent SQL injection
func Test_pgxIdentifier_Sanitize_WhenValidTableName_ThenQuotesCorrectly(t *testing.T) {
	tests := []struct {
		name          string
		tableName     string
		expectedQuote string
	}{
		{
			name:          "simple table name",
			tableName:     "transactions",
			expectedQuote: `"transactions"`,
		},
		{
			name:          "synthetic table name",
			tableName:     "transactions_synthetic",
			expectedQuote: `"transactions_synthetic"`,
		},
		{
			name:          "table name with special chars (sanitized)",
			tableName:     "transactions; DROP TABLE users--",
			expectedQuote: `"transactions; DROP TABLE users--"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pgx.Identifier{tt.tableName}.Sanitize()

			assert.Equal(t, tt.expectedQuote, result, "pgx.Identifier should properly quote table names")
		})
	}
}

// Test_pgxIdentifier_Sanitize_WhenUsedInQuery_ThenPreventsSQLInjection demonstrates
// that even with malicious input, pgx.Identifier safely quotes the identifier
func Test_pgxIdentifier_Sanitize_WhenUsedInQuery_ThenPreventsSQLInjection(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "SQL injection with semicolon",
			input:    "transactions; DROP TABLE users;--",
			expected: `"transactions; DROP TABLE users;--"`,
		},
		{
			name:     "SQL injection with quotes",
			input:    "transactions' OR '1'='1",
			expected: `"transactions' OR '1'='1"`,
		},
		{
			name:     "SQL injection with escaped quotes",
			input:    `transactions"; DELETE FROM users WHERE "1"="1`,
			expected: `"transactions""; DELETE FROM users WHERE ""1""=""1"`, // pgx doubles internal quotes
		},
		{
			name:     "SQL injection with backticks",
			input:    "transactions`; DROP DATABASE;--",
			expected: "\"transactions`; DROP DATABASE;--\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sanitize the table name
			sanitized := pgx.Identifier{tt.input}.Sanitize()

			// Verify it matches expected output
			assert.Equal(t, tt.expected, sanitized, "pgx.Identifier should properly escape and quote")

			// Verify it's properly quoted (starts with " and ends with ")
			assert.True(t, len(sanitized) >= 2, "sanitized output should have at least opening and closing quotes")
			assert.Equal(t, '"', rune(sanitized[0]), "should start with double quote")
			assert.Equal(t, '"', rune(sanitized[len(sanitized)-1]), "should end with double quote")
		})
	}
}
