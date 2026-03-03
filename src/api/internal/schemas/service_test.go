package schemas

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_Service_Create_WhenValidRequest_ThenCreatesSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &CreateSchemaRequest{
		Name:        "test-schema",
		Description: "Test schema",
		SampleJSON: map[string]any{
			"amount":   100.50,
			"currency": "USD",
		},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByName(gomock.Any(), req.Name).Return(nil, pgx.ErrNoRows)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
	service := NewService(mockRepo, nil)

	schema, err := service.Create(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, req.Name, schema.Name)
	assert.Equal(t, req.Description, schema.Description)
	assert.NotEmpty(t, schema.ExtractedFields)
}

func Test_Service_Create_WhenNameExists_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &CreateSchemaRequest{
		Name:        "existing-schema",
		Description: "Test schema",
		SampleJSON:  map[string]any{"amount": 100.50},
	}
	existing := &EventSchema{ID: uuid.New(), Name: req.Name}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByName(gomock.Any(), req.Name).Return(existing, nil)
	service := NewService(mockRepo, nil)

	schema, err := service.Create(context.Background(), req)

	assert.Nil(t, schema)
	assert.ErrorIs(t, err, ErrSchemaNameExists)
}

func Test_Service_Create_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &CreateSchemaRequest{
		Name:        "test-schema",
		Description: "Test schema",
		SampleJSON:  map[string]any{"amount": 100.50},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByName(gomock.Any(), req.Name).Return(nil, pgx.ErrNoRows)
	mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
	service := NewService(mockRepo, nil)

	schema, err := service.Create(context.Background(), req)

	assert.Nil(t, schema)
	assert.Error(t, err)
}

func Test_Service_GetByID_WhenSchemaExists_ThenReturnsSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	expected := &EventSchema{ID: id, Name: "test-schema"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(expected, nil)
	service := NewService(mockRepo, nil)

	schema, err := service.GetByID(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, expected, schema)
}

func Test_Service_GetByID_WhenSchemaNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, pgx.ErrNoRows)
	service := NewService(mockRepo, nil)

	schema, err := service.GetByID(context.Background(), id)

	assert.Nil(t, schema)
	assert.ErrorIs(t, err, ErrSchemaNotFound)
}

func Test_Service_List_WhenSuccess_ThenReturnsSchemas(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expected := []EventSchema{
		{ID: uuid.New(), Name: "schema1"},
		{ID: uuid.New(), Name: "schema2"},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().List(gomock.Any()).Return(expected, nil)
	service := NewService(mockRepo, nil)

	schemas, err := service.List(context.Background())

	require.NoError(t, err)
	assert.Equal(t, expected, schemas)
}

func Test_Service_List_WhenEmpty_ThenReturnsEmptySlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().List(gomock.Any()).Return(nil, nil)
	service := NewService(mockRepo, nil)

	schemas, err := service.List(context.Background())

	require.NoError(t, err)
	assert.Empty(t, schemas)
}

func Test_Service_Update_WhenValidRequest_ThenUpdatesSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	existing := &EventSchema{
		ID:          id,
		Name:        "old-name",
		Description: "old description",
		SampleJSON:  map[string]any{"field1": "value1"},
	}
	req := &UpdateSchemaRequest{
		Name:        "new-name",
		Description: "new description",
		SampleJSON:  map[string]any{"field2": "value2"},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	mockRepo.EXPECT().GetByName(gomock.Any(), req.Name).Return(nil, pgx.ErrNoRows)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
	service := NewService(mockRepo, nil)

	schema, err := service.Update(context.Background(), id, req)

	require.NoError(t, err)
	assert.Equal(t, req.Name, schema.Name)
	assert.Equal(t, req.Description, schema.Description)
}

func Test_Service_Update_WhenSchemaNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	req := &UpdateSchemaRequest{Name: "new-name"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, pgx.ErrNoRows)
	service := NewService(mockRepo, nil)

	schema, err := service.Update(context.Background(), id, req)

	assert.Nil(t, schema)
	assert.ErrorIs(t, err, ErrSchemaNotFound)
}

func Test_Service_Update_WhenNewNameExists_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	existing := &EventSchema{ID: id, Name: "old-name"}
	req := &UpdateSchemaRequest{Name: "existing-name"}
	other := &EventSchema{ID: uuid.New(), Name: "existing-name"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	mockRepo.EXPECT().GetByName(gomock.Any(), req.Name).Return(other, nil)
	service := NewService(mockRepo, nil)

	schema, err := service.Update(context.Background(), id, req)

	assert.Nil(t, schema)
	assert.ErrorIs(t, err, ErrSchemaNameExists)
}

func Test_Service_Delete_WhenValid_ThenDeletesSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	existing := &EventSchema{ID: id, Name: "test-schema"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	mockRepo.EXPECT().HasRulesReferencing(gomock.Any(), id).Return(false, nil)
	mockRepo.EXPECT().Delete(gomock.Any(), id).Return(nil)
	service := NewService(mockRepo, nil)

	err := service.Delete(context.Background(), id)

	assert.NoError(t, err)
}

func Test_Service_Delete_WhenSchemaNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, pgx.ErrNoRows)
	service := NewService(mockRepo, nil)

	err := service.Delete(context.Background(), id)

	assert.ErrorIs(t, err, ErrSchemaNotFound)
}

func Test_Service_Delete_WhenHasRules_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	existing := &EventSchema{ID: id, Name: "test-schema"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	mockRepo.EXPECT().HasRulesReferencing(gomock.Any(), id).Return(true, nil)
	service := NewService(mockRepo, nil)

	err := service.Delete(context.Background(), id)

	assert.ErrorIs(t, err, ErrSchemaHasRules)
}

func Test_ExtractFields_WhenSimpleObject_ThenExtractsFields(t *testing.T) {
	data := map[string]any{
		"name":   "John",
		"age":    30,
		"active": true,
	}

	fields := ExtractFields(data, "", 0)

	assert.Len(t, fields, 3)
	assert.Contains(t, fields, ExtractedField{Path: "name", Type: FieldTypeString, Nullable: false, SampleValue: "John"})
	assert.Contains(t, fields, ExtractedField{Path: "age", Type: FieldTypeNumber, Nullable: false, SampleValue: 30})
	assert.Contains(t, fields, ExtractedField{Path: "active", Type: FieldTypeBoolean, Nullable: false, SampleValue: true})
}

func Test_ExtractFields_WhenNestedObject_ThenExtractsNested(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "John",
			"age":  30,
		},
	}

	fields := ExtractFields(data, "", 0)

	assert.Len(t, fields, 2)
	assert.Contains(t, fields, ExtractedField{Path: "user.name", Type: FieldTypeString, Nullable: false, SampleValue: "John"})
	assert.Contains(t, fields, ExtractedField{Path: "user.age", Type: FieldTypeNumber, Nullable: false, SampleValue: 30})
}

func Test_ExtractFields_WhenMaxDepth_ThenStopsRecursion(t *testing.T) {
	data := map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"level3": map[string]any{
					"level4": map[string]any{
						"level5": map[string]any{
							"level6": "too deep",
						},
					},
				},
			},
		},
	}

	fields := ExtractFields(data, "", 0)

	for _, field := range fields {
		depth := 0
		for _, c := range field.Path {
			if c == '.' {
				depth++
			}
		}
		assert.LessOrEqual(t, depth, MaxNestingDepth-1, "field %s exceeds max depth", field.Path)
	}
}

func Test_inferType_WhenString_ThenReturnsStringType(t *testing.T) {
	fieldType, nullable := inferType("test")

	assert.Equal(t, FieldTypeString, fieldType)
	assert.False(t, nullable)
}

func Test_inferType_WhenNumber_ThenReturnsNumberType(t *testing.T) {
	tests := []any{100, int32(100), int64(100), float32(100.5), float64(100.5)}

	for _, val := range tests {
		fieldType, nullable := inferType(val)

		assert.Equal(t, FieldTypeNumber, fieldType)
		assert.False(t, nullable)
	}
}

func Test_inferType_WhenBoolean_ThenReturnsBooleanType(t *testing.T) {
	fieldType, nullable := inferType(true)

	assert.Equal(t, FieldTypeBoolean, fieldType)
	assert.False(t, nullable)
}

func Test_inferType_WhenNull_ThenReturnsNullType(t *testing.T) {
	fieldType, nullable := inferType(nil)

	assert.Equal(t, FieldTypeNull, fieldType)
	assert.True(t, nullable)
}

func Test_inferType_WhenArray_ThenReturnsArrayType(t *testing.T) {
	fieldType, nullable := inferType([]any{1, 2, 3})

	assert.Equal(t, FieldTypeArray, fieldType)
	assert.False(t, nullable)
}

func Test_inferType_WhenObject_ThenReturnsObjectType(t *testing.T) {
	fieldType, nullable := inferType(map[string]any{"key": "value"})

	assert.Equal(t, FieldTypeObject, fieldType)
	assert.False(t, nullable)
}

func Test_inferType_WhenDateTimeString_ThenReturnsDateTimeType(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"RFC3339", "2025-01-15T10:30:00Z"},
		{"RFC3339Nano", "2025-01-15T10:30:00.123456789Z"},
		{"RFC3339 with timezone", "2025-01-15T10:30:00-07:00"},
		{"ISO 8601 without timezone", "2025-01-15T10:30:00Z"},
		{"ISO 8601 with milliseconds", "2025-01-15T10:30:00.000Z"},
		{"Date only", "2025-01-15"},
		{"Date time space format", "2025-01-15 10:30:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldType, nullable := inferType(tt.value)

			assert.Equal(t, FieldTypeDateTime, fieldType)
			assert.False(t, nullable)
		})
	}
}

func Test_inferType_WhenRegularString_ThenReturnsStringType(t *testing.T) {
	tests := []string{
		"regular string",
		"not a date",
		"12345",
		"2025/01/15",
		"01-15-2025",
	}

	for _, val := range tests {
		fieldType, nullable := inferType(val)

		assert.Equal(t, FieldTypeString, fieldType)
		assert.False(t, nullable)
	}
}

func Test_Service_GetRulesReferencingSchema_WhenSchemaHasRules_ThenReturnsRuleNames(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	expectedRules := []string{"rule1", "rule2", "rule3"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetRulesReferencingSchema(gomock.Any(), id).Return(expectedRules, nil)
	service := NewService(mockRepo, nil)

	rules, err := service.GetRulesReferencingSchema(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, expectedRules, rules)
}

func Test_Service_GetRulesReferencingSchema_WhenNoRules_ThenReturnsEmptySlice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetRulesReferencingSchema(gomock.Any(), id).Return([]string{}, nil)
	service := NewService(mockRepo, nil)

	rules, err := service.GetRulesReferencingSchema(context.Background(), id)

	require.NoError(t, err)
	assert.Empty(t, rules)
}

func Test_Service_GetRulesReferencingSchema_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetRulesReferencingSchema(gomock.Any(), id).Return(nil, errors.New("database error"))
	service := NewService(mockRepo, nil)

	rules, err := service.GetRulesReferencingSchema(context.Background(), id)

	assert.Nil(t, rules)
	assert.Error(t, err)
}

func Test_Service_ParseSampleJSON_WhenSchemaExists_ThenReExtractsFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	existing := &EventSchema{
		ID:   id,
		Name: "test-schema",
		SampleJSON: map[string]any{
			"amount":   100.50,
			"currency": "USD",
		},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, schema *EventSchema) error {
			assert.NotEmpty(t, schema.ExtractedFields)
			return nil
		},
	)
	service := NewService(mockRepo, nil)

	schema, err := service.ParseSampleJSON(context.Background(), id)

	require.NoError(t, err)
	assert.NotNil(t, schema)
	assert.NotEmpty(t, schema.ExtractedFields)
}

func Test_Service_ParseSampleJSON_WhenSchemaNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, pgx.ErrNoRows)
	service := NewService(mockRepo, nil)

	schema, err := service.ParseSampleJSON(context.Background(), id)

	assert.Nil(t, schema)
	assert.ErrorIs(t, err, ErrSchemaNotFound)
}

func Test_Service_ParseSampleJSON_WhenUpdateFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	existing := &EventSchema{
		ID:         id,
		Name:       "test-schema",
		SampleJSON: map[string]any{"field": "value"},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
	service := NewService(mockRepo, nil)

	schema, err := service.ParseSampleJSON(context.Background(), id)

	assert.Nil(t, schema)
	assert.Error(t, err)
}

func Test_Service_Create_WhenGetByNameFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	req := &CreateSchemaRequest{
		Name:        "test-schema",
		Description: "Test schema",
		SampleJSON:  map[string]any{"amount": 100.50},
	}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByName(gomock.Any(), req.Name).Return(nil, errors.New("database error"))
	service := NewService(mockRepo, nil)

	schema, err := service.Create(context.Background(), req)

	assert.Nil(t, schema)
	assert.Error(t, err)
}

func Test_Service_GetByID_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("database error"))
	service := NewService(mockRepo, nil)

	schema, err := service.GetByID(context.Background(), id)

	assert.Nil(t, schema)
	assert.Error(t, err)
}

func Test_Service_List_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().List(gomock.Any()).Return(nil, errors.New("database error"))
	service := NewService(mockRepo, nil)

	schemas, err := service.List(context.Background())

	assert.Nil(t, schemas)
	assert.Error(t, err)
}

func Test_Service_Update_WhenGetByIDFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	req := &UpdateSchemaRequest{Name: "new-name"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("database error"))
	service := NewService(mockRepo, nil)

	schema, err := service.Update(context.Background(), id, req)

	assert.Nil(t, schema)
	assert.Error(t, err)
}

func Test_Service_Update_WhenGetByNameFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	existing := &EventSchema{ID: id, Name: "old-name"}
	req := &UpdateSchemaRequest{Name: "new-name"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	mockRepo.EXPECT().GetByName(gomock.Any(), req.Name).Return(nil, errors.New("database error"))
	service := NewService(mockRepo, nil)

	schema, err := service.Update(context.Background(), id, req)

	assert.Nil(t, schema)
	assert.Error(t, err)
}

func Test_Service_Update_WhenUpdateFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	existing := &EventSchema{ID: id, Name: "old-name", Description: "old"}
	req := &UpdateSchemaRequest{Description: "new description"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
	service := NewService(mockRepo, nil)

	schema, err := service.Update(context.Background(), id, req)

	assert.Nil(t, schema)
	assert.Error(t, err)
}

func Test_Service_Delete_WhenGetByIDFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("database error"))
	service := NewService(mockRepo, nil)

	err := service.Delete(context.Background(), id)

	assert.Error(t, err)
}

func Test_Service_Delete_WhenHasRulesReferencingFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	existing := &EventSchema{ID: id, Name: "test-schema"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	mockRepo.EXPECT().HasRulesReferencing(gomock.Any(), id).Return(false, errors.New("database error"))
	service := NewService(mockRepo, nil)

	err := service.Delete(context.Background(), id)

	assert.Error(t, err)
}

func Test_Service_Delete_WhenDeleteFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	existing := &EventSchema{ID: id, Name: "test-schema"}
	mockRepo := NewMockRepository(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(existing, nil)
	mockRepo.EXPECT().HasRulesReferencing(gomock.Any(), id).Return(false, nil)
	mockRepo.EXPECT().Delete(gomock.Any(), id).Return(errors.New("database error"))
	service := NewService(mockRepo, nil)

	err := service.Delete(context.Background(), id)

	assert.Error(t, err)
}

func Test_ExtractFields_WhenEmptyObject_ThenReturnsEmptySlice(t *testing.T) {
	data := map[string]any{}

	fields := ExtractFields(data, "", 0)

	assert.Empty(t, fields)
}

func Test_ExtractFields_WhenArrayField_ThenExtractsArrayType(t *testing.T) {
	data := map[string]any{
		"tags": []any{"tag1", "tag2"},
	}

	fields := ExtractFields(data, "", 0)

	assert.Len(t, fields, 1)
	assert.Equal(t, "tags", fields[0].Path)
	assert.Equal(t, FieldTypeArray, fields[0].Type)
}

func Test_ExtractFields_WhenNullField_ThenExtractsNullType(t *testing.T) {
	data := map[string]any{
		"optional": nil,
	}

	fields := ExtractFields(data, "", 0)

	assert.Len(t, fields, 1)
	assert.Equal(t, "optional", fields[0].Path)
	assert.Equal(t, FieldTypeNull, fields[0].Type)
	assert.True(t, fields[0].Nullable)
}

func Test_ExtractFields_WhenDateTimeField_ThenExtractsDateTimeType(t *testing.T) {
	data := map[string]any{
		"created_at": "2025-01-15T10:30:00Z",
		"updated_at": "2025-01-15T10:30:00.123456789Z",
		"timestamp":  "2025-01-15T10:30:00-07:00",
		"date":       "2025-01-15",
	}

	fields := ExtractFields(data, "", 0)

	assert.Len(t, fields, 4)
	for _, field := range fields {
		assert.Equal(t, FieldTypeDateTime, field.Type)
		assert.False(t, field.Nullable)
	}
}

func Test_ExtractFields_WhenNestedDateTimeField_ThenExtractsDateTimeType(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"created_at": "2025-01-15T10:30:00Z",
			"name":       "John",
		},
	}

	fields := ExtractFields(data, "", 0)

	assert.Len(t, fields, 2)
	dateTimeField := findFieldByPath(fields, "user.created_at")
	require.NotNil(t, dateTimeField)
	assert.Equal(t, FieldTypeDateTime, dateTimeField.Type)
}

func Test_Service_GenerateEvents_WhenSuccess_ThenReturnsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	schema := &EventSchema{
		ID:   id,
		Name: "test-schema",
		ExtractedFields: []ExtractedField{
			{Path: "amount", Type: FieldTypeNumber},
			{Path: "currency", Type: FieldTypeString},
		},
	}
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(schema, nil)
	mockEnqueuer.EXPECT().
		EnqueueTransactionWithPriority(gomock.Any(), gomock.Any(), "default").
		Return(nil, nil).
		Times(3)

	service := NewService(mockRepo, mockEnqueuer)
	req := &GenerateEventsRequest{Count: 3}

	resp, err := service.GenerateEvents(context.Background(), id, req)

	require.NoError(t, err)
	assert.Equal(t, id, resp.SchemaID)
	assert.Equal(t, 3, resp.GeneratedCount)
	assert.Equal(t, 0, resp.FailedCount)
	assert.Contains(t, resp.Message, "Successfully generated")
}

func Test_Service_GenerateEvents_WhenNoEnqueuer_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockRepository(ctrl)
	service := NewService(mockRepo, nil)

	id := uuid.New()
	req := &GenerateEventsRequest{Count: 1}

	resp, err := service.GenerateEvents(context.Background(), id, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task enqueuer not configured")
}

func Test_Service_GenerateEvents_WhenSchemaNotFound_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, pgx.ErrNoRows)
	service := NewService(mockRepo, mockEnqueuer)

	req := &GenerateEventsRequest{Count: 1}

	resp, err := service.GenerateEvents(context.Background(), id, req)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrSchemaNotFound)
}

func Test_Service_GenerateEvents_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(nil, errors.New("database error"))
	service := NewService(mockRepo, mockEnqueuer)

	req := &GenerateEventsRequest{Count: 1}

	resp, err := service.GenerateEvents(context.Background(), id, req)

	assert.Nil(t, resp)
	assert.Error(t, err)
}

func Test_Service_GenerateEvents_WhenEnqueueFails_ThenReportsFailures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	id := uuid.New()
	schema := &EventSchema{
		ID:   id,
		Name: "test-schema",
		ExtractedFields: []ExtractedField{
			{Path: "amount", Type: FieldTypeNumber},
		},
	}
	mockRepo := NewMockRepository(ctrl)
	mockEnqueuer := NewMockTaskEnqueuer(ctrl)
	mockRepo.EXPECT().GetByID(gomock.Any(), id).Return(schema, nil)
	mockEnqueuer.EXPECT().
		EnqueueTransactionWithPriority(gomock.Any(), gomock.Any(), "default").
		Return(nil, errors.New("queue full")).
		Times(2)

	service := NewService(mockRepo, mockEnqueuer)
	req := &GenerateEventsRequest{Count: 2}

	resp, err := service.GenerateEvents(context.Background(), id, req)

	require.NoError(t, err)
	assert.Equal(t, id, resp.SchemaID)
	assert.Equal(t, 0, resp.GeneratedCount)
	assert.Equal(t, 2, resp.FailedCount)
	assert.Contains(t, resp.Message, "failed")
}

func findFieldByPath(fields []ExtractedField, path string) *ExtractedField {
	for i := range fields {
		if fields[i].Path == path {
			return &fields[i]
		}
	}
	return nil
}
