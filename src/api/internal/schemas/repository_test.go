package schemas

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func Test_NewPostgresRepository_WhenCalled_ThenReturnsRepository(t *testing.T) {
	var db *pgxpool.Pool
	var redis *redis.Client

	repo := NewPostgresRepository(db, redis)

	assert.NotNil(t, repo)
	assert.Implements(t, (*Repository)(nil), repo)
}

func Test_NewPostgresRepository_WhenCalledWithNilRedis_ThenReturnsRepository(t *testing.T) {
	var db *pgxpool.Pool

	repo := NewPostgresRepository(db, nil)

	assert.NotNil(t, repo)
	assert.Implements(t, (*Repository)(nil), repo)
}

func Test_PostgresRepository_publishInvalidation_WhenRedisIsNil_ThenDoesNothing(t *testing.T) {
	repo := &PostgresRepository{
		db:    nil,
		redis: nil,
	}

	schemaID := uuid.New()

	repo.publishInvalidation(context.Background(), schemaID)

	assert.NotNil(t, repo)
}
