package schemas

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// Repository defines the interface for schema persistence operations
type Repository interface {
	Create(ctx context.Context, schema *EventSchema) error
	GetByID(ctx context.Context, id uuid.UUID) (*EventSchema, error)
	GetByName(ctx context.Context, name string) (*EventSchema, error)
	List(ctx context.Context) ([]EventSchema, error)
	Update(ctx context.Context, schema *EventSchema) error
	Delete(ctx context.Context, id uuid.UUID) error
	HasRulesReferencing(ctx context.Context, schemaID uuid.UUID) (bool, error)
	GetRulesReferencingSchema(ctx context.Context, schemaID uuid.UUID) ([]string, error)
}

// PostgresRepository implements Repository using PostgreSQL
type PostgresRepository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *pgxpool.Pool, redis *redis.Client) *PostgresRepository {
	return &PostgresRepository{
		db:    db,
		redis: redis,
	}
}

const (
	schemaInvalidateChannel = "schema:invalidate"
)

// publishInvalidation publishes a schema invalidation event to Redis
func (r *PostgresRepository) publishInvalidation(ctx context.Context, schemaID uuid.UUID) {
	if r.redis != nil {
		r.redis.Publish(ctx, schemaInvalidateChannel, schemaID.String())
	}
}

func (r *PostgresRepository) Create(ctx context.Context, schema *EventSchema) error {
	sampleJSON, err := json.Marshal(schema.SampleJSON)
	if err != nil {
		return err
	}

	extractedFields, err := json.Marshal(schema.ExtractedFields)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO event_schemas (id, name, description, sample_json, extracted_fields, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = r.db.Exec(ctx, query,
		schema.ID,
		schema.Name,
		schema.Description,
		sampleJSON,
		extractedFields,
		schema.CreatedAt,
		schema.UpdatedAt,
	)

	if err == nil {
		r.publishInvalidation(ctx, schema.ID)
	}

	return err
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

func (r *PostgresRepository) GetByName(ctx context.Context, name string) (*EventSchema, error) {
	query := `
		SELECT id, name, description, sample_json, extracted_fields, created_at, updated_at
		FROM event_schemas
		WHERE name = $1
	`

	var schema EventSchema
	var sampleJSON, extractedFields []byte

	err := r.db.QueryRow(ctx, query, name).Scan(
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

func (r *PostgresRepository) List(ctx context.Context) ([]EventSchema, error) {
	query := `
		SELECT id, name, description, sample_json, extracted_fields, created_at, updated_at
		FROM event_schemas
		ORDER BY name ASC
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

func (r *PostgresRepository) Update(ctx context.Context, schema *EventSchema) error {
	sampleJSON, err := json.Marshal(schema.SampleJSON)
	if err != nil {
		return err
	}

	extractedFields, err := json.Marshal(schema.ExtractedFields)
	if err != nil {
		return err
	}

	query := `
		UPDATE event_schemas
		SET name = $2, description = $3, sample_json = $4, extracted_fields = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query,
		schema.ID,
		schema.Name,
		schema.Description,
		sampleJSON,
		extractedFields,
		schema.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	r.publishInvalidation(ctx, schema.ID)
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM event_schemas WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	r.publishInvalidation(ctx, id)
	return nil
}

func (r *PostgresRepository) HasRulesReferencing(ctx context.Context, schemaID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM rules WHERE schema_id = $1)`

	var exists bool
	err := r.db.QueryRow(ctx, query, schemaID).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) GetRulesReferencingSchema(ctx context.Context, schemaID uuid.UUID) ([]string, error) {
	query := `SELECT name FROM rules WHERE schema_id = $1`

	rows, err := r.db.Query(ctx, query, schemaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}

	return names, rows.Err()
}
