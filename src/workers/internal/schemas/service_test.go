package schemas

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// LoadSchemas Tests

func Test_SchemaService_LoadSchemas_WhenSuccess_ThenCachesSchemas(t *testing.T) {
	schemas := []EventSchema{
		{ID: uuid.New(), Name: "schema1"},
		{ID: uuid.New(), Name: "schema2"},
	}
	mockRepo := &mockSchemaRepository{schemas: schemas}
	service := NewSchemaService(mockRepo, nil)

	err := service.LoadSchemas(context.Background())

	require.NoError(t, err)
	assert.Len(t, service.schemas, 2)
}

func Test_SchemaService_LoadSchemas_WhenRepositoryFails_ThenReturnsError(t *testing.T) {
	mockRepo := &mockSchemaRepository{err: errors.New("database error")}
	service := NewSchemaService(mockRepo, nil)

	err := service.LoadSchemas(context.Background())

	assert.Error(t, err)
}

func Test_SchemaService_LoadSchemas_WhenCalledMultipleTimes_ThenReplacesCache(t *testing.T) {
	firstSchemas := []EventSchema{{ID: uuid.New(), Name: "schema1"}}
	secondSchemas := []EventSchema{{ID: uuid.New(), Name: "schema2"}, {ID: uuid.New(), Name: "schema3"}}
	mockRepo := &mockSchemaRepository{schemas: firstSchemas}
	service := NewSchemaService(mockRepo, nil)

	_ = service.LoadSchemas(context.Background())
	assert.Len(t, service.schemas, 1)

	mockRepo.schemas = secondSchemas
	_ = service.LoadSchemas(context.Background())

	assert.Len(t, service.schemas, 2)
}

func Test_SchemaService_LoadSchemas_WhenEmptyResult_ThenClearsCache(t *testing.T) {
	schemas := []EventSchema{{ID: uuid.New(), Name: "schema1"}}
	mockRepo := &mockSchemaRepository{schemas: schemas}
	service := NewSchemaService(mockRepo, nil)

	_ = service.LoadSchemas(context.Background())
	assert.Len(t, service.schemas, 1)

	mockRepo.schemas = []EventSchema{}
	_ = service.LoadSchemas(context.Background())

	assert.Empty(t, service.schemas)
}

// GetSchema Tests

func Test_SchemaService_GetSchema_WhenSchemaExists_ThenReturnsSchema(t *testing.T) {
	schemaID := uuid.New()
	schemas := []EventSchema{{ID: schemaID, Name: "test-schema"}}
	mockRepo := &mockSchemaRepository{schemas: schemas}
	service := NewSchemaService(mockRepo, nil)
	_ = service.LoadSchemas(context.Background())

	result := service.GetSchema(schemaID)

	require.NotNil(t, result)
	assert.Equal(t, schemaID, result.ID)
	assert.Equal(t, "test-schema", result.Name)
}

func Test_SchemaService_GetSchema_WhenSchemaNotExists_ThenReturnsNil(t *testing.T) {
	mockRepo := &mockSchemaRepository{schemas: []EventSchema{}}
	service := NewSchemaService(mockRepo, nil)
	_ = service.LoadSchemas(context.Background())

	result := service.GetSchema(uuid.New())

	assert.Nil(t, result)
}

func Test_SchemaService_GetSchema_WhenCacheEmpty_ThenReturnsNil(t *testing.T) {
	mockRepo := &mockSchemaRepository{schemas: []EventSchema{}}
	service := NewSchemaService(mockRepo, nil)

	result := service.GetSchema(uuid.New())

	assert.Nil(t, result)
}

func Test_SchemaService_GetSchema_WhenConcurrentReads_ThenHandlesSafely(t *testing.T) {
	schemaID := uuid.New()
	schemas := []EventSchema{{ID: schemaID, Name: "test-schema"}}
	mockRepo := &mockSchemaRepository{schemas: schemas}
	service := NewSchemaService(mockRepo, nil)
	_ = service.LoadSchemas(context.Background())

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := service.GetSchema(schemaID)
			assert.NotNil(t, result)
		}()
	}

	wg.Wait()
}

// InvalidateSchema Tests

func Test_SchemaService_InvalidateSchema_WhenSchemaExists_ThenReloadsFromDB(t *testing.T) {
	schemaID := uuid.New()
	schemas := []EventSchema{{ID: schemaID, Name: "old-name"}}
	updatedSchema := &EventSchema{ID: schemaID, Name: "new-name"}
	mockRepo := &mockSchemaRepository{schemas: schemas, getByIDResult: updatedSchema}
	service := NewSchemaService(mockRepo, nil)
	_ = service.LoadSchemas(context.Background())

	service.InvalidateSchema(context.Background(), schemaID)

	result := service.GetSchema(schemaID)
	require.NotNil(t, result)
	assert.Equal(t, "new-name", result.Name)
}

func Test_SchemaService_InvalidateSchema_WhenSchemaDeleted_ThenRemovesFromCache(t *testing.T) {
	schemaID := uuid.New()
	schemas := []EventSchema{{ID: schemaID, Name: "test-schema"}}
	mockRepo := &mockSchemaRepository{schemas: schemas, getByIDErr: errors.New("not found")}
	service := NewSchemaService(mockRepo, nil)
	_ = service.LoadSchemas(context.Background())

	service.InvalidateSchema(context.Background(), schemaID)

	result := service.GetSchema(schemaID)
	assert.Nil(t, result)
}

func Test_SchemaService_InvalidateSchema_WhenSchemaNotInCache_ThenAddsToCache(t *testing.T) {
	schemaID := uuid.New()
	newSchema := &EventSchema{ID: schemaID, Name: "new-schema"}
	mockRepo := &mockSchemaRepository{schemas: []EventSchema{}, getByIDResult: newSchema}
	service := NewSchemaService(mockRepo, nil)
	_ = service.LoadSchemas(context.Background())

	service.InvalidateSchema(context.Background(), schemaID)

	result := service.GetSchema(schemaID)
	require.NotNil(t, result)
	assert.Equal(t, "new-schema", result.Name)
}

// GetAllSchemas Tests

func Test_SchemaService_GetAllSchemas_WhenSchemasExist_ThenReturnsAll(t *testing.T) {
	schemas := []EventSchema{
		{ID: uuid.New(), Name: "schema1"},
		{ID: uuid.New(), Name: "schema2"},
	}
	mockRepo := &mockSchemaRepository{schemas: schemas}
	service := NewSchemaService(mockRepo, nil)
	_ = service.LoadSchemas(context.Background())

	result := service.GetAllSchemas()

	assert.Len(t, result, 2)
}

func Test_SchemaService_GetAllSchemas_WhenCacheEmpty_ThenReturnsEmpty(t *testing.T) {
	mockRepo := &mockSchemaRepository{schemas: []EventSchema{}}
	service := NewSchemaService(mockRepo, nil)

	result := service.GetAllSchemas()

	assert.Empty(t, result)
}

func Test_SchemaService_GetAllSchemas_WhenModifyingReturned_ThenDoesNotAffectCache(t *testing.T) {
	schemaID := uuid.New()
	schemas := []EventSchema{{ID: schemaID, Name: "test-schema"}}
	mockRepo := &mockSchemaRepository{schemas: schemas}
	service := NewSchemaService(mockRepo, nil)
	_ = service.LoadSchemas(context.Background())

	result := service.GetAllSchemas()
	delete(result, schemaID)

	cached := service.GetSchema(schemaID)
	assert.NotNil(t, cached, "cache should not be affected by external modification")
}

// SubscribeToInvalidations Tests

func Test_SchemaService_SubscribeToInvalidations_WhenRedisNil_ThenReturnsImmediately(t *testing.T) {
	mockRepo := &mockSchemaRepository{}
	service := NewSchemaService(mockRepo, nil)

	done := make(chan bool, 1)
	go func() {
		service.SubscribeToInvalidations(context.Background())
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("SubscribeToInvalidations did not return immediately when redis is nil")
	}
}

// Note: Additional SubscribeToInvalidations tests with Redis pub/sub are better suited
// as integration tests since mocking redis.Client and redis.PubSub is complex

// Mock implementations

type mockSchemaRepository struct {
	schemas       []EventSchema
	err           error
	getByIDResult *EventSchema
	getByIDErr    error
}

func (m *mockSchemaRepository) GetByID(ctx context.Context, id uuid.UUID) (*EventSchema, error) {
	if m.getByIDErr != nil {
		return nil, m.getByIDErr
	}
	if m.getByIDResult != nil {
		return m.getByIDResult, nil
	}
	for i := range m.schemas {
		if m.schemas[i].ID == id {
			return &m.schemas[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (m *mockSchemaRepository) ListAll(ctx context.Context) ([]EventSchema, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.schemas, nil
}
