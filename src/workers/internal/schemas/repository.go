package schemas

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository defines the interface for schema read operations (used by worker)
type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*EventSchema, error)
	ListAll(ctx context.Context) ([]EventSchema, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*EventSchema, error) {
	query := `
		SELECT id, name, description, sample_json, extracted_fields, created_at, updated_at
		FROM event_schemas
		WHERE id = $1
	`

	var schema EventSchema
	var sampleJSON, extractedFields []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&schema.ID,
		&schema.Name,
		&schema.Description,
		&sampleJSON,
		&extractedFields,
		&schema.CreatedAt,
		&schema.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(sampleJSON, &schema.SampleJSON); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(extractedFields, &schema.ExtractedFields); err != nil {
		return nil, err
	}

	return &schema, nil
}

func (r *PostgresRepository) ListAll(ctx context.Context) ([]EventSchema, error) {
	query := `
		SELECT id, name, description, sample_json, extracted_fields, created_at, updated_at
		FROM event_schemas
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schemas []EventSchema
	for rows.Next() {
		var schema EventSchema
		var sampleJSON, extractedFields []byte

		if err := rows.Scan(
			&schema.ID,
			&schema.Name,
			&schema.Description,
			&sampleJSON,
			&extractedFields,
			&schema.CreatedAt,
			&schema.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(sampleJSON, &schema.SampleJSON); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(extractedFields, &schema.ExtractedFields); err != nil {
			return nil, err
		}

		schemas = append(schemas, schema)
	}

	return schemas, rows.Err()
}
