package transactions

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionHistoryRepository defines the interface for transaction history queries
type TransactionHistoryRepository interface {
	// CountByFieldInTimeWindow counts transactions where the specified field matches the value within a time window
	// groupFieldPath: path to the field in metadata (e.g., "origin", "user.id", "customer_id")
	// fieldValue: value to match in the groupFieldPath
	CountByFieldInTimeWindow(ctx context.Context, groupFieldPath string, fieldValue string, timeWindowSeconds int) (int, error)
	// SumFieldByFieldInTimeWindow sums a numeric field where another field matches the value within a time window
	// groupFieldPath: path to the field used for grouping (e.g., "origin", "user.id")
	// groupFieldValue: value to match in the groupFieldPath
	// sumFieldPath: path to the numeric field to sum (e.g., "amount", "value", "total")
	SumFieldByFieldInTimeWindow(ctx context.Context, groupFieldPath string, groupFieldValue string, sumFieldPath string, timeWindowSeconds int) (float64, error)
}

// PostgresHistoryRepository is the PostgreSQL implementation
type PostgresHistoryRepository struct {
	db *pgxpool.Pool
}

// NewPostgresHistoryRepository creates a new PostgreSQL history repository
func NewPostgresHistoryRepository(db *pgxpool.Pool) TransactionHistoryRepository {
	return &PostgresHistoryRepository{db: db}
}

var (
	// validFieldPathPattern validates that field paths only contain alphanumeric, underscore, and dot characters
	validFieldPathPattern = regexp.MustCompile(`^[a-zA-Z0-9_.]+$`)
)

// buildJSONPathQuery builds a safe JSONB path query from a field path
// e.g., "user.id" extracts metadata->'user'->>'id'
// Validates the field path to prevent SQL injection
func buildJSONPathQuery(fieldPath string) (string, error) {
	// Validate field path to prevent SQL injection
	if !validFieldPathPattern.MatchString(fieldPath) {
		return "", fmt.Errorf("invalid field path: %s (only alphanumeric, underscore, and dot characters allowed)", fieldPath)
	}

	parts := strings.Split(fieldPath, ".")
	if len(parts) == 1 {
		// Escape single quotes in the field name (though validation should prevent them)
		escaped := strings.ReplaceAll(parts[0], "'", "''")
		return fmt.Sprintf("metadata->>'%s'", escaped), nil
	}

	// Build path: metadata->'part1'->'part2'->>'lastPart'
	path := "metadata"
	for i := 0; i < len(parts)-1; i++ {
		// Escape single quotes in each part
		escaped := strings.ReplaceAll(parts[i], "'", "''")
		path = fmt.Sprintf("%s->'%s'", path, escaped)
	}
	// Escape single quotes in the last part
	escaped := strings.ReplaceAll(parts[len(parts)-1], "'", "''")
	path = fmt.Sprintf("%s->>'%s'", path, escaped)
	return path, nil
}

func (r *PostgresHistoryRepository) CountByFieldInTimeWindow(ctx context.Context, groupFieldPath string, fieldValue string, timeWindowSeconds int) (int, error) {
	fieldQuery, err := buildJSONPathQuery(groupFieldPath)
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM transactions 
		WHERE %s = $1
		AND created_at > NOW() - INTERVAL '1 second' * $2
	`, fieldQuery)

	var count int
	err = r.db.QueryRow(ctx, query, fieldValue, timeWindowSeconds).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresHistoryRepository) SumFieldByFieldInTimeWindow(ctx context.Context, groupFieldPath string, groupFieldValue string, sumFieldPath string, timeWindowSeconds int) (float64, error) {
	groupFieldQuery, err := buildJSONPathQuery(groupFieldPath)
	if err != nil {
		return 0.0, err
	}

	sumFieldQuery, err := buildJSONPathQuery(sumFieldPath)
	if err != nil {
		return 0.0, err
	}

	query := fmt.Sprintf(`
		SELECT COALESCE(SUM((%s)::numeric), 0) 
		FROM transactions 
		WHERE %s = $1
		AND created_at > NOW() - INTERVAL '1 second' * $2
		AND %s IS NOT NULL
	`, sumFieldQuery, groupFieldQuery, sumFieldQuery)

	var sum float64
	err = r.db.QueryRow(ctx, query, groupFieldValue, timeWindowSeconds).Scan(&sum)
	if err != nil {
		return 0.0, err
	}

	return sum, nil
}
