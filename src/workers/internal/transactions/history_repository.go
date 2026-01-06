package transactions

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionHistoryRepository defines the interface for transaction history queries
type TransactionHistoryRepository interface {
	// CountByAccountInTimeWindow counts transactions for an account within a time window
	CountByAccountInTimeWindow(ctx context.Context, account string, timeWindowSeconds int) (int, error)
	// SumAmountByAccountInTimeWindow sums transaction amounts for an account within a time window
	SumAmountByAccountInTimeWindow(ctx context.Context, account string, timeWindowSeconds int) (float64, error)
}

// PostgresHistoryRepository is the PostgreSQL implementation
type PostgresHistoryRepository struct {
	db *pgxpool.Pool
}

// NewPostgresHistoryRepository creates a new PostgreSQL history repository
func NewPostgresHistoryRepository(db *pgxpool.Pool) TransactionHistoryRepository {
	return &PostgresHistoryRepository{db: db}
}

func (r *PostgresHistoryRepository) CountByAccountInTimeWindow(ctx context.Context, account string, timeWindowSeconds int) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM transactions 
		WHERE origin = $1 
		AND created_at > NOW() - INTERVAL '1 second' * $2
	`

	var count int
	err := r.db.QueryRow(ctx, query, account, timeWindowSeconds).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresHistoryRepository) SumAmountByAccountInTimeWindow(ctx context.Context, account string, timeWindowSeconds int) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE origin = $1 
		AND created_at > NOW() - INTERVAL '1 second' * $2
	`

	var sum float64
	err := r.db.QueryRow(ctx, query, account, timeWindowSeconds).Scan(&sum)
	if err != nil {
		return 0.0, err
	}

	return sum, nil
}
