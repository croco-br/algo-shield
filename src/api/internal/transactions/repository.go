package transactions

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/shared"
	"github.com/algo-shield/algo-shield/src/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TransactionFilter contains filter criteria for listing transactions
type TransactionFilter struct {
	Status    *models.TransactionStatus
	SchemaID  *uuid.UUID
	StartDate *time.Time
	EndDate   *time.Time
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
	RejectTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	CountTransactions(ctx context.Context) (int, error)
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	var schemaName *string
	tableName := shared.GetTransactionsTable(ctx)

	// Use pgx.Identifier to safely quote table name and prevent SQL injection
	query := fmt.Sprintf(`
		SELECT t.id, t.schema_id, es.name as schema_name,
		       t.status, t.processing_time,
		       t.matched_rules, t.metadata, t.created_at, t.processed_at
		FROM %s t
		LEFT JOIN event_schemas es ON t.schema_id = es.id
		WHERE t.id = $1
	`, pgx.Identifier{tableName}.Sanitize())

	err := r.db.QueryRow(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.SchemaID,
		&schemaName,
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

	transaction.SchemaName = schemaName
	return &transaction, nil
}

func (r *PostgresRepository) ListTransactions(ctx context.Context, limit, offset int) ([]models.Transaction, error) {
	tableName := shared.GetTransactionsTable(ctx)
	// Use pgx.Identifier to safely quote table name and prevent SQL injection
	query := fmt.Sprintf(`
		SELECT t.id, t.schema_id, es.name as schema_name,
		       t.status, t.processing_time,
		       t.matched_rules, t.metadata, t.created_at, t.processed_at
		FROM %s t
		LEFT JOIN event_schemas es ON t.schema_id = es.id
		ORDER BY t.created_at DESC
		LIMIT $1 OFFSET $2
	`, pgx.Identifier{tableName}.Sanitize())

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		var transaction models.Transaction
		var schemaName *string
		err := rows.Scan(
			&transaction.ID,
			&transaction.SchemaID,
			&schemaName,
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
		transaction.SchemaName = schemaName
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *PostgresRepository) ListTransactionsWithFilter(ctx context.Context, filter TransactionFilter, limit, offset int) ([]models.Transaction, int, error) {
	tableName := shared.GetTransactionsTable(ctx)
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

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count query - use pgx.Identifier to safely quote table name and prevent SQL injection
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", pgx.Identifier{tableName}.Sanitize(), whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Data query with LEFT JOIN to get schema name - use pgx.Identifier to safely quote table name
	query := fmt.Sprintf(`
		SELECT t.id, t.schema_id, es.name as schema_name,
		       t.status, t.processing_time,
		       t.matched_rules, t.metadata, t.created_at, t.processed_at
		FROM %s t
		LEFT JOIN event_schemas es ON t.schema_id = es.id
		%s
		ORDER BY t.created_at DESC
		LIMIT $%d OFFSET $%d
	`, pgx.Identifier{tableName}.Sanitize(), whereClause, argNum, argNum+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	transactions := make([]models.Transaction, 0)
	for rows.Next() {
		var transaction models.Transaction
		var schemaName *string
		err := rows.Scan(
			&transaction.ID,
			&transaction.SchemaID,
			&schemaName,
			&transaction.Status,
			&transaction.ProcessingTime,
			&transaction.MatchedRules,
			&transaction.Metadata,
			&transaction.CreatedAt,
			&transaction.ProcessedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		transaction.SchemaName = schemaName
		transactions = append(transactions, transaction)
	}

	return transactions, total, nil
}

func (r *PostgresRepository) ApproveTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	tableName := shared.GetTransactionsTable(ctx)
	now := time.Now()
	// Use pgx.Identifier to safely quote table name and prevent SQL injection
	updateQuery := fmt.Sprintf(`
		UPDATE %s
		SET status = $1, processed_at = $2
		WHERE id = $3 AND (status = $4 OR status = $5)
	`, pgx.Identifier{tableName}.Sanitize())

	result, err := r.db.Exec(ctx, updateQuery, models.StatusApproved, now, id, models.StatusPending, models.StatusInReview)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	selectQuery := fmt.Sprintf(`
		SELECT t.id, t.schema_id, es.name as schema_name,
		       t.status, t.processing_time, t.matched_rules, t.metadata, t.created_at, t.processed_at
		FROM %s t
		LEFT JOIN event_schemas es ON t.schema_id = es.id
		WHERE t.id = $1
	`, pgx.Identifier{tableName}.Sanitize())

	var transaction models.Transaction
	var schemaName *string
	err = r.db.QueryRow(ctx, selectQuery, id).Scan(
		&transaction.ID,
		&transaction.SchemaID,
		&schemaName,
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

	transaction.SchemaName = schemaName
	return &transaction, nil
}

func (r *PostgresRepository) RejectTransaction(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	tableName := shared.GetTransactionsTable(ctx)
	now := time.Now()
	// Use pgx.Identifier to safely quote table name and prevent SQL injection
	updateQuery := fmt.Sprintf(`
		UPDATE %s
		SET status = $1, processed_at = $2
		WHERE id = $3 AND (status = $4 OR status = $5)
	`, pgx.Identifier{tableName}.Sanitize())

	result, err := r.db.Exec(ctx, updateQuery, models.StatusRejected, now, id, models.StatusPending, models.StatusInReview)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		return nil, pgx.ErrNoRows
	}

	selectQuery := fmt.Sprintf(`
		SELECT t.id, t.schema_id, es.name as schema_name,
		       t.status, t.processing_time, t.matched_rules, t.metadata, t.created_at, t.processed_at
		FROM %s t
		LEFT JOIN event_schemas es ON t.schema_id = es.id
		WHERE t.id = $1
	`, pgx.Identifier{tableName}.Sanitize())

	var transaction models.Transaction
	var schemaName *string
	err = r.db.QueryRow(ctx, selectQuery, id).Scan(
		&transaction.ID,
		&transaction.SchemaID,
		&schemaName,
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

	transaction.SchemaName = schemaName
	return &transaction, nil
}

func (r *PostgresRepository) CountTransactions(ctx context.Context) (int, error) {
	tableName := shared.GetTransactionsTable(ctx)
	var count int
	// Use pgx.Identifier to safely quote table name and prevent SQL injection
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", pgx.Identifier{tableName}.Sanitize())
	err := r.db.QueryRow(ctx, query).Scan(&count)
	return count, err
}
