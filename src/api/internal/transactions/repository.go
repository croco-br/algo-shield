package transactions

import (
	"context"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository is the PostgreSQL implementation of Repository
type PostgresRepository struct {
	db *pgxpool.Pool
}

// Repository defines the interface for transaction data access operations
type Repository interface {
	GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error)
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction

	query := `
		SELECT id, external_id, amount, currency, from_account, to_account, 
		       type, status, risk_score, risk_level, processing_time, 
		       matched_rules, metadata, created_at, processed_at
		FROM transactions
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.ExternalID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.FromAccount,
		&transaction.ToAccount,
		&transaction.Type,
		&transaction.Status,
		&transaction.RiskScore,
		&transaction.RiskLevel,
		&transaction.ProcessingTime,
		&transaction.MatchedRules,
		&transaction.Metadata,
		&transaction.CreatedAt,
		&transaction.ProcessedAt,
	)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *PostgresRepository) ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	query := `
		SELECT id, external_id, amount, currency, from_account, to_account, 
		       type, status, risk_score, risk_level, processing_time, 
		       matched_rules, metadata, created_at, processed_at
		FROM transactions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.ExternalID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.FromAccount,
			&transaction.ToAccount,
			&transaction.Type,
			&transaction.Status,
			&transaction.RiskScore,
			&transaction.RiskLevel,
			&transaction.ProcessingTime,
			&transaction.MatchedRules,
			&transaction.Metadata,
			&transaction.CreatedAt,
			&transaction.ProcessedAt,
		)
		if err != nil {
			continue
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
