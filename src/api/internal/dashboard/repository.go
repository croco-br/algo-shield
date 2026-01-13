package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/shared"
	"github.com/jackc/pgx/v5/pgxpool"
)

// StatusCount represents count of transactions per status
type StatusCount struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

// TemporalCount represents count of transactions per time bucket
type TemporalCount struct {
	Bucket time.Time `json:"bucket"`
	Count  int64     `json:"count"`
}

// Repository defines the interface for dashboard data access
type Repository interface {
	GetStatusDistribution(ctx context.Context) ([]StatusCount, error)
	GetTemporalCounts(ctx context.Context, period string) ([]TemporalCount, error)
	GetTotalCount(ctx context.Context) (int64, error)
}

// PostgresRepository is the PostgreSQL implementation
type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetStatusDistribution(ctx context.Context) ([]StatusCount, error) {
	tableName := shared.GetTransactionsTable(ctx)
	query := fmt.Sprintf(`
		SELECT status, COUNT(*) as count
		FROM %s
		GROUP BY status
		ORDER BY count DESC
	`, tableName)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []StatusCount
	for rows.Next() {
		var sc StatusCount
		if err := rows.Scan(&sc.Status, &sc.Count); err != nil {
			continue
		}
		result = append(result, sc)
	}

	return result, nil
}

func (r *PostgresRepository) GetTemporalCounts(ctx context.Context, period string) ([]TemporalCount, error) {
	tableName := shared.GetTransactionsTable(ctx)
	var interval, bucketSize string
	switch period {
	case "24h":
		interval = "24 hours"
		bucketSize = "1 hour"
	case "7d":
		interval = "7 days"
		bucketSize = "1 day"
	case "30d":
		interval = "30 days"
		bucketSize = "1 day"
	default:
		interval = "24 hours"
		bucketSize = "1 hour"
	}

	query := fmt.Sprintf(`
		SELECT date_trunc($1, created_at) as bucket, COUNT(*) as count
		FROM %s
		WHERE created_at >= NOW() - $2::interval
		GROUP BY bucket
		ORDER BY bucket ASC
	`, tableName)

	// Map bucket size to PostgreSQL date_trunc argument
	truncArg := "hour"
	if bucketSize == "1 day" {
		truncArg = "day"
	}

	rows, err := r.db.Query(ctx, query, truncArg, interval)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []TemporalCount
	for rows.Next() {
		var tc TemporalCount
		if err := rows.Scan(&tc.Bucket, &tc.Count); err != nil {
			continue
		}
		result = append(result, tc)
	}

	return result, nil
}

func (r *PostgresRepository) GetTotalCount(ctx context.Context) (int64, error) {
	tableName := shared.GetTransactionsTable(ctx)
	var count int64
	err := r.db.QueryRow(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
	return count, err
}
