package schemas

import (
	"time"

	"github.com/google/uuid"
)

// FieldType represents the inferred type of a JSON field
type FieldType string

const (
	FieldTypeString  FieldType = "string"
	FieldTypeNumber  FieldType = "number"
	FieldTypeBoolean FieldType = "boolean"
	FieldTypeArray   FieldType = "array"
	FieldTypeObject  FieldType = "object"
	FieldTypeNull    FieldType = "null"
)

// ExtractedField represents a field extracted from sample JSON
type ExtractedField struct {
	Path        string    `json:"path"`
	Type        FieldType `json:"type"`
	Nullable    bool      `json:"nullable"`
	SampleValue any       `json:"sample_value,omitempty"`
}

// EventSchema represents a user-defined event schema
type EventSchema struct {
	ID              uuid.UUID        `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description,omitempty"`
	SampleJSON      map[string]any   `json:"sample_json"`
	ExtractedFields []ExtractedField `json:"extracted_fields"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}
