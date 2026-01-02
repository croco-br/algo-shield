package schemas

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	// SchemaInvalidateChannel is the Redis pub/sub channel for schema invalidation
	SchemaInvalidateChannel = "schema:invalidate"
)

// SchemaService handles schema loading and caching for the worker
// Uses sync.RWMutex for thread-safe reads and writes
type SchemaService struct {
	repo    Repository
	redis   *redis.Client
	schemas map[uuid.UUID]*EventSchema
	mu      sync.RWMutex
}

// NewSchemaService creates a new schema service for the worker
func NewSchemaService(repo Repository, redisClient *redis.Client) *SchemaService {
	return &SchemaService{
		repo:    repo,
		redis:   redisClient,
		schemas: make(map[uuid.UUID]*EventSchema),
	}
}

// LoadSchemas loads all schemas from the repository into the cache
func (s *SchemaService) LoadSchemas(ctx context.Context) error {
	schemas, err := s.repo.ListAll(ctx)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear and rebuild cache
	s.schemas = make(map[uuid.UUID]*EventSchema)
	for i := range schemas {
		schema := schemas[i] // Create a copy to avoid pointer aliasing
		s.schemas[schema.ID] = &schema
	}

	log.Printf("Loaded %d schemas into cache", len(schemas))
	return nil
}

// GetSchema returns a cached schema by ID
// Returns nil if not found
func (s *SchemaService) GetSchema(id uuid.UUID) *EventSchema {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.schemas[id]
}

// InvalidateSchema removes a schema from the cache, forcing a reload on next access
func (s *SchemaService) InvalidateSchema(ctx context.Context, id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.schemas, id)

	// Reload from database if needed
	schema, err := s.repo.GetByID(ctx, id)
	if err != nil {
		// Schema was deleted, just remove from cache
		log.Printf("Schema %s was deleted or not found", id)
		return
	}

	s.schemas[id] = schema
	log.Printf("Reloaded schema %s", id)
}

// SubscribeToInvalidations subscribes to Redis pub/sub for schema invalidation events
// This should be called in a goroutine
func (s *SchemaService) SubscribeToInvalidations(ctx context.Context) {
	if s.redis == nil {
		log.Println("Redis not available, schema invalidation subscription disabled")
		return
	}

	pubsub := s.redis.Subscribe(ctx, SchemaInvalidateChannel)
	defer func() {
		if err := pubsub.Close(); err != nil {
			log.Printf("Error closing schema invalidation subscription: %v", err)
		}
	}()

	log.Println("Subscribed to schema invalidation channel")

	for {
		select {
		case <-ctx.Done():
			log.Println("Schema invalidation subscription stopped")
			return
		default:
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return // Context cancelled
				}
				log.Printf("Error receiving schema invalidation message: %v", err)
				continue
			}

			schemaID, err := uuid.Parse(msg.Payload)
			if err != nil {
				log.Printf("Invalid schema ID in invalidation message: %s", msg.Payload)
				continue
			}

			s.InvalidateSchema(ctx, schemaID)
		}
	}
}

// GetAllSchemas returns all cached schemas
func (s *SchemaService) GetAllSchemas() map[uuid.UUID]*EventSchema {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[uuid.UUID]*EventSchema)
	for k, v := range s.schemas {
		result[k] = v
	}
	return result
}
