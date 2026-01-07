package branding

import (
	"testing"

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

func Test_NewPostgresRepository_WhenCalledWithNilDB_ThenReturnsRepository(t *testing.T) {
	var redis *redis.Client

	repo := NewPostgresRepository(nil, redis)

	assert.NotNil(t, repo)
	assert.Implements(t, (*Repository)(nil), repo)
}
