package transactions

import (
	"context"
	"encoding/json"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for transaction data access operations
type Repository interface {
	// SaveTransaction saves a processed transaction to the database
	SaveTransaction(ctx context.Context, transaction *models.Transaction) error
}

// PostgresRepository is the PostgreSQL implementation of Repository
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL transaction repository
func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveTransaction(ctx context.Context, transaction *models.Transaction) error {
	matchedRulesJSON, _ := json.Marshal(transaction.MatchedRules)
	metadataJSON, _ := json.Marshal(transaction.Metadata)

	query := `
		INSERT INTO transactions (
			id, external_id, amount, currency, from_account, to_account, 
			type, status, risk_score, risk_level, processing_time, 
			matched_rules, metadata, created_at, processed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err := r.db.Exec(ctx, query,
		transaction.ID,
		transaction.ExternalID,
		transaction.Amount,
		transaction.Currency,
		transaction.FromAccount,
		transaction.ToAccount,
		transaction.Type,
		transaction.Status,
		transaction.RiskScore,
		transaction.RiskLevel,
		transaction.ProcessingTime,
		matchedRulesJSON,
		metadataJSON,
		transaction.CreatedAt,
		transaction.ProcessedAt,
	)

	return err
}
