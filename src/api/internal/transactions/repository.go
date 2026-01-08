package transactions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionFilter contains filter criteria for listing transactions
type TransactionFilter struct {
	Status    *models.TransactionStatus
	SchemaID  *uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
	MinAmount *float64
	MaxAmount *float64
}

// PostgresRepository is the PostgreSQL implementation of Repository
type PostgresRepository struct {
	db *pgxpool.Pool
}

// Repository defines the interface for transaction data access operations
type Repository interface {
	GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error)
	ListTransactionsWithFilter(ctx context.Context, filter TransactionFilter, limit, offset int) ([]models.Transaction, int, error)
	ApproveTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	CountTransactions(ctx context.Context) (int, error)
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction

	query := `
		SELECT id, external_id, schema_id, amount, currency, origin, destination, 
		       type, status, processing_time, 
		       matched_rules, metadata, created_at, processed_at
		FROM transactions
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.ExternalID,
		&transaction.SchemaID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Origin,
		&transaction.Destination,
		&transaction.Type,
		&transaction.Status,
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
		SELECT id, external_id, schema_id, amount, currency, origin, destination, 
		       type, status, processing_time, 
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
			&transaction.SchemaID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Origin,
			&transaction.Destination,
			&transaction.Type,
			&transaction.Status,
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

func (r *PostgresRepository) ListTransactionsWithFilter(ctx context.Context, filter TransactionFilter, limit, offset int) ([]models.Transaction, int, error) {
	var conditions []string
	var args []any
	argNum := 1

	if filter.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argNum))
		args = append(args, string(*filter.Status))
		argNum++
	}
	if filter.SchemaID != nil {
		conditions = append(conditions, fmt.Sprintf("schema_id = $%d", argNum))
		args = append(args, *filter.SchemaID)
		argNum++
	}
	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argNum))
		args = append(args, *filter.StartDate)
		argNum++
	}
	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argNum))
		args = append(args, *filter.EndDate)
		argNum++
	}
	if filter.MinAmount != nil {
		conditions = append(conditions, fmt.Sprintf("amount >= $%d", argNum))
		args = append(args, *filter.MinAmount)
		argNum++
	}
	if filter.MaxAmount != nil {
		conditions = append(conditions, fmt.Sprintf("amount <= $%d", argNum))
		args = append(args, *filter.MaxAmount)
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM transactions %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query
	query := fmt.Sprintf(`
		SELECT id, external_id, schema_id, amount, currency, origin, destination, 
		       type, status, processing_time, 
		       matched_rules, metadata, created_at, processed_at
		FROM transactions
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.ExternalID,
			&transaction.SchemaID,
			&transaction.Amount,
			&transaction.Currency,
			&transaction.Origin,
			&transaction.Destination,
			&transaction.Type,
			&transaction.Status,
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

	return transactions, total, nil
}

func (r *PostgresRepository) ApproveTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	now := time.Now()
	query := `
		UPDATE transactions 
		SET status = $1, processed_at = $2
		WHERE id = $3 AND status = $4
		RETURNING id, external_id, schema_id, amount, currency, origin, destination,
		          type, status, processing_time, matched_rules, metadata, created_at, processed_at
	`

	var transaction models.Transaction
	err := r.db.QueryRow(ctx, query, models.StatusApproved, now, id, models.StatusInReview).Scan(
		&transaction.ID,
		&transaction.ExternalID,
		&transaction.SchemaID,
		&transaction.Amount,
		&transaction.Currency,
		&transaction.Origin,
		&transaction.Destination,
		&transaction.Type,
		&transaction.Status,
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

func (r *PostgresRepository) CountTransactions(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM transactions").Scan(&count)
	return count, err
}
