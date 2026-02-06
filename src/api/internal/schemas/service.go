package schemas

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// MaxNestingDepth is the maximum depth for JSON field extraction
const MaxNestingDepth = 5

// Service errors
var (
	ErrSchemaNotFound    = errors.New("schema not found")
	ErrSchemaNameExists  = errors.New("schema with this name already exists")
	ErrSchemaHasRules    = errors.New("schema is referenced by rules and cannot be deleted")
	ErrInvalidSampleJSON = errors.New("sample_json must be a valid JSON object")
)

// ServiceInterface defines the interface for schema business logic
type ServiceInterface interface {
	Create(ctx context.Context, req *CreateSchemaRequest) (*EventSchema, error)
	GetByID(ctx context.Context, id uuid.UUID) (*EventSchema, error)
	List(ctx context.Context) ([]EventSchema, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateSchemaRequest) (*EventSchema, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetRulesReferencingSchema(ctx context.Context, id uuid.UUID) ([]string, error)
	ParseSampleJSON(ctx context.Context, id uuid.UUID) (*EventSchema, error)
	GenerateEvents(ctx context.Context, id uuid.UUID, req *GenerateEventsRequest) (*GenerateEventsResponse, error)
}

// TaskEnqueuer defines interface for enqueueing tasks (Asynq)
// Follows Dependency Inversion Principle
// Returns any to avoid coupling with asynq.TaskInfo
type TaskEnqueuer interface {
	EnqueueTransactionWithPriority(ctx context.Context, event map[string]any, priority string) (any, error)
}

// Service provides business logic for schema operations
type Service struct {
	repo     Repository
	enqueuer TaskEnqueuer
}

// NewService creates a new schema service with task enqueuer
func NewService(repo Repository, enqueuer TaskEnqueuer) *Service {
	return &Service{
		repo:     repo,
		enqueuer: enqueuer,
	}
}

// Create creates a new event schema from sample JSON
func (s *Service) Create(ctx context.Context, req *CreateSchemaRequest) (*EventSchema, error) {
	// Check for duplicate name
	existing, err := s.repo.GetByName(ctx, req.Name)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrSchemaNameExists
	}

	// Extract fields from sample JSON
	fields := ExtractFields(req.SampleJSON, "", 0)

	now := time.Now()
	schema := &EventSchema{
		ID:              uuid.New(),
		Name:            req.Name,
		Description:     req.Description,
		SampleJSON:      req.SampleJSON,
		ExtractedFields: fields,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.Create(ctx, schema); err != nil {
		return nil, err
	}

	return schema, nil
}

// GetByID retrieves a schema by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*EventSchema, error) {
	schema, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSchemaNotFound
		}
		return nil, err
	}
	return schema, nil
}

// List returns all schemas
func (s *Service) List(ctx context.Context) ([]EventSchema, error) {
	schemas, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	if schemas == nil {
		schemas = []EventSchema{}
	}
	return schemas, nil
}

// Update updates an existing schema
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateSchemaRequest) (*EventSchema, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSchemaNotFound
		}
		return nil, err
	}

	// Check for duplicate name if name is being changed
	if req.Name != "" && req.Name != existing.Name {
		other, err := s.repo.GetByName(ctx, req.Name)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		if other != nil {
			return nil, ErrSchemaNameExists
		}
		existing.Name = req.Name
	}

	if req.Description != "" {
		existing.Description = req.Description
	}

	// Re-extract fields if sample JSON is updated
	if req.SampleJSON != nil {
		existing.SampleJSON = req.SampleJSON
		existing.ExtractedFields = ExtractFields(req.SampleJSON, "", 0)
	}

	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete deletes a schema by ID
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if schema exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSchemaNotFound
		}
		return err
	}

	// Check if any rules reference this schema
	hasRules, err := s.repo.HasRulesReferencing(ctx, id)
	if err != nil {
		return err
	}
	if hasRules {
		return ErrSchemaHasRules
	}

	return s.repo.Delete(ctx, id)
}

// GetRulesReferencingSchema returns the names of rules that reference a schema
func (s *Service) GetRulesReferencingSchema(ctx context.Context, id uuid.UUID) ([]string, error) {
	return s.repo.GetRulesReferencingSchema(ctx, id)
}

// ParseSampleJSON re-extracts fields from a schema's sample JSON
func (s *Service) ParseSampleJSON(ctx context.Context, id uuid.UUID) (*EventSchema, error) {
	schema, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSchemaNotFound
		}
		return nil, err
	}

	schema.ExtractedFields = ExtractFields(schema.SampleJSON, "", 0)
	schema.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, schema); err != nil {
		return nil, err
	}

	return schema, nil
}

// enqueueConcurrency controls the number of parallel Redis enqueue operations
const enqueueConcurrency = 20

// GenerateEvents generates synthetic events from a schema and queues them for processing
func (s *Service) GenerateEvents(ctx context.Context, id uuid.UUID, req *GenerateEventsRequest) (*GenerateEventsResponse, error) {
	if s.enqueuer == nil {
		return nil, errors.New("task enqueuer not configured")
	}

	schema, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSchemaNotFound
		}
		return nil, err
	}

	generator := NewEventGenerator()

	// Generate all events first (CPU-bound, fast)
	events := make([]map[string]any, req.Count)
	for i := 0; i < req.Count; i++ {
		event := generator.GenerateEvent(schema)
		event["_schema_id"] = id.String()
		event["_synthetic"] = true
		events[i] = event
	}

	// Enqueue in parallel with bounded concurrency
	var generated atomic.Int64
	var failed atomic.Int64
	sem := make(chan struct{}, enqueueConcurrency)
	var wg sync.WaitGroup

	for _, event := range events {
		if ctx.Err() != nil {
			break
		}

		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore slot

		go func(ev map[string]any) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore slot

			if ctx.Err() != nil {
				failed.Add(1)
				return
			}

			_, err := s.enqueuer.EnqueueTransactionWithPriority(ctx, ev, "default")
			if err != nil {
				failed.Add(1)
				return
			}
			generated.Add(1)
		}(event)
	}

	wg.Wait()

	gen := int(generated.Load())
	fail := int(failed.Load())

	var message string
	if fail == 0 {
		message = fmt.Sprintf("Successfully generated and queued %d events", gen)
	} else {
		message = fmt.Sprintf("Generated %d events (%d succeeded, %d failed)", req.Count, gen, fail)
	}

	return &GenerateEventsResponse{
		SchemaID:       id,
		GeneratedCount: gen,
		FailedCount:    fail,
		Message:        message,
	}, nil
}

// ExtractFields recursively extracts fields from a JSON object
func ExtractFields(data map[string]any, prefix string, depth int) []ExtractedField {
	if depth >= MaxNestingDepth {
		return nil
	}

	var fields []ExtractedField

	for key, value := range data {
		path := key
		if prefix != "" {
			path = fmt.Sprintf("%s.%s", prefix, key)
		}

		fieldType, nullable := inferType(value)
		field := ExtractedField{
			Path:        path,
			Type:        fieldType,
			Nullable:    nullable,
			SampleValue: value,
		}

		// For objects, recurse into nested fields
		if fieldType == FieldTypeObject {
			if nested, ok := value.(map[string]any); ok {
				nestedFields := ExtractFields(nested, path, depth+1)
				fields = append(fields, nestedFields...)
			}
		} else {
			fields = append(fields, field)
		}
	}

	return fields
}

// inferType determines the field type from a JSON value
func inferType(value any) (FieldType, bool) {
	if value == nil {
		return FieldTypeNull, true
	}

	switch v := value.(type) {
	case bool:
		return FieldTypeBoolean, false
	case float64, float32, int, int32, int64:
		return FieldTypeNumber, false
	case string:
		if isDateTimeString(v) {
			return FieldTypeDateTime, false
		}
		return FieldTypeString, false
	case []any:
		return FieldTypeArray, false
	case map[string]any:
		return FieldTypeObject, false
	default:
		// Handle JSON numbers that come as json.Number
		if _, ok := v.(interface{ Float64() (float64, error) }); ok {
			return FieldTypeNumber, false
		}
		return FieldTypeString, false
	}
}

// isDateTimeString checks if a string represents a valid datetime/timestamp
// Supports common formats: RFC3339, RFC3339Nano, and other ISO 8601 variants
func isDateTimeString(s string) bool {
	if s == "" {
		return false
	}

	// Common datetime formats to try
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05.000-07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, s); err == nil {
			return true
		}
	}

	return false
}
