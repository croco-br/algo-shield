//go:build integration

package schemas_test

import (
	"context"
	"testing"
	"time"

	"github.com/algo-shield/algo-shield/src/api/internal/schemas"
	"github.com/algo-shield/algo-shield/src/api/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_SchemasRepository_Create_StoresSchema(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	schemaID := uuid.New()
	now := time.Now()
	schema := &schemas.EventSchema{
		ID:          schemaID,
		Name:        "test_schema",
		Description: "Test schema",
		SampleJSON: map[string]any{
			"event_type": "test",
			"amount":     100.0,
		},
		ExtractedFields: []schemas.ExtractedField{
			{Path: "event_type", Type: schemas.FieldTypeString, Nullable: false},
			{Path: "amount", Type: schemas.FieldTypeNumber, Nullable: false},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := repo.Create(ctx, schema)

	require.NoError(t, err)

	var storedName string
	err = testDB.Postgres.QueryRow(ctx, `
		SELECT name FROM event_schemas WHERE id = $1
	`, schemaID).Scan(&storedName)
	require.NoError(t, err)
	assert.Equal(t, "test_schema", storedName)
}

func TestIntegration_SchemasRepository_Create_DuplicateName_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	schemaID1 := uuid.New()
	schemaID2 := uuid.New()
	now := time.Now()

	schema1 := &schemas.EventSchema{
		ID:              schemaID1,
		Name:            "duplicate_name",
		Description:     "First schema",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err := repo.Create(ctx, schema1)
	require.NoError(t, err)

	schema2 := &schemas.EventSchema{
		ID:              schemaID2,
		Name:            "duplicate_name",
		Description:     "Second schema",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = repo.Create(ctx, schema2)

	assert.Error(t, err)
}

func TestIntegration_SchemasRepository_GetByID_ReturnsSchema(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	schemaID := uuid.New()
	now := time.Now()
	schema := &schemas.EventSchema{
		ID:          schemaID,
		Name:        "get_by_id_schema",
		Description: "Get by ID test",
		SampleJSON: map[string]any{
			"test": "value",
		},
		ExtractedFields: []schemas.ExtractedField{
			{Path: "test", Type: schemas.FieldTypeString, Nullable: false},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
	err := repo.Create(ctx, schema)
	require.NoError(t, err)

	result, err := repo.GetByID(ctx, schemaID)

	require.NoError(t, err)
	assert.Equal(t, schemaID, result.ID)
	assert.Equal(t, "get_by_id_schema", result.Name)
	assert.Equal(t, "Get by ID test", result.Description)
	assert.NotEmpty(t, result.SampleJSON)
	assert.Len(t, result.ExtractedFields, 1)
}

func TestIntegration_SchemasRepository_GetByID_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	nonExistentID := uuid.New()

	result, err := repo.GetByID(ctx, nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_SchemasRepository_GetByName_ReturnsSchema(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	schemaID := uuid.New()
	now := time.Now()
	schema := &schemas.EventSchema{
		ID:              schemaID,
		Name:            "get_by_name_schema",
		Description:     "Get by name test",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err := repo.Create(ctx, schema)
	require.NoError(t, err)

	result, err := repo.GetByName(ctx, "get_by_name_schema")

	require.NoError(t, err)
	assert.Equal(t, schemaID, result.ID)
	assert.Equal(t, "get_by_name_schema", result.Name)
}

func TestIntegration_SchemasRepository_GetByName_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	result, err := repo.GetByName(ctx, "non_existent_schema")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_SchemasRepository_List_ReturnsAllSchemas(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	now := time.Now()
	schema1 := &schemas.EventSchema{
		ID:              uuid.New(),
		Name:            "schema_a",
		Description:     "Schema A",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	schema2 := &schemas.EventSchema{
		ID:              uuid.New(),
		Name:            "schema_b",
		Description:     "Schema B",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err := repo.Create(ctx, schema1)
	require.NoError(t, err)
	err = repo.Create(ctx, schema2)
	require.NoError(t, err)

	result, err := repo.List(ctx)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 2)
}

func TestIntegration_SchemasRepository_Update_UpdatesSchema(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	schemaID := uuid.New()
	now := time.Now()
	schema := &schemas.EventSchema{
		ID:              schemaID,
		Name:            "update_schema",
		Description:     "Original description",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err := repo.Create(ctx, schema)
	require.NoError(t, err)

	schema.Description = "Updated description"
	schema.UpdatedAt = time.Now()
	err = repo.Update(ctx, schema)

	require.NoError(t, err)

	result, err := repo.GetByID(ctx, schemaID)
	require.NoError(t, err)
	assert.Equal(t, "Updated description", result.Description)
}

func TestIntegration_SchemasRepository_Update_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	nonExistentID := uuid.New()
	schema := &schemas.EventSchema{
		ID:              nonExistentID,
		Name:            "non_existent",
		Description:     "Test",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err := repo.Update(ctx, schema)

	assert.Error(t, err)
}

func TestIntegration_SchemasRepository_Delete_DeletesSchema(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	schemaID := uuid.New()
	now := time.Now()
	schema := &schemas.EventSchema{
		ID:              schemaID,
		Name:            "delete_schema",
		Description:     "To be deleted",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err := repo.Create(ctx, schema)
	require.NoError(t, err)

	err = repo.Delete(ctx, schemaID)

	require.NoError(t, err)

	result, err := repo.GetByID(ctx, schemaID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestIntegration_SchemasRepository_Delete_NotFound_ReturnsError(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	nonExistentID := uuid.New()

	err := repo.Delete(ctx, nonExistentID)

	assert.Error(t, err)
}

func TestIntegration_SchemasRepository_HasRulesReferencing_WithRules_ReturnsTrue(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	schemaID := uuid.New()
	now := time.Now()
	schema := &schemas.EventSchema{
		ID:              schemaID,
		Name:            "referenced_schema",
		Description:     "Referenced schema",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err := repo.Create(ctx, schema)
	require.NoError(t, err)

	ruleID := uuid.New()
	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, ruleID, "test_rule", "Test rule", "allow", 0, true, "{}", schemaID, time.Now(), time.Now())
	require.NoError(t, err)

	hasRules, err := repo.HasRulesReferencing(ctx, schemaID)

	require.NoError(t, err)
	assert.True(t, hasRules)
}

func TestIntegration_SchemasRepository_HasRulesReferencing_NoRules_ReturnsFalse(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	schemaID := uuid.New()
	now := time.Now()
	schema := &schemas.EventSchema{
		ID:              schemaID,
		Name:            "unreferenced_schema",
		Description:     "Unreferenced schema",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err := repo.Create(ctx, schema)
	require.NoError(t, err)

	hasRules, err := repo.HasRulesReferencing(ctx, schemaID)

	require.NoError(t, err)
	assert.False(t, hasRules)
}

func TestIntegration_SchemasRepository_GetRulesReferencingSchema_ReturnsRuleNames(t *testing.T) {
	testDB := testutil.SetupTestDB(t)
	repo := schemas.NewPostgresRepository(testDB.Postgres, testDB.Redis)
	ctx := context.Background()

	schemaID := uuid.New()
	now := time.Now()
	schema := &schemas.EventSchema{
		ID:              schemaID,
		Name:            "referenced_schema",
		Description:     "Referenced schema",
		SampleJSON:      map[string]any{"key": "value"},
		ExtractedFields: []schemas.ExtractedField{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	err := repo.Create(ctx, schema)
	require.NoError(t, err)

	ruleID1 := uuid.New()
	ruleID2 := uuid.New()
	_, err = testDB.Postgres.Exec(ctx, `
		INSERT INTO rules (id, name, description, action, priority, enabled, conditions, schema_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10),
		       ($11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`, ruleID1, "rule1", "Rule 1", "allow", 0, true, "{}", schemaID, time.Now(), time.Now(),
		ruleID2, "rule2", "Rule 2", "block", 0, true, "{}", schemaID, time.Now(), time.Now())
	require.NoError(t, err)

	ruleNames, err := repo.GetRulesReferencingSchema(ctx, schemaID)

	require.NoError(t, err)
	assert.Len(t, ruleNames, 2)
	assert.Contains(t, ruleNames, "rule1")
	assert.Contains(t, ruleNames, "rule2")
}
