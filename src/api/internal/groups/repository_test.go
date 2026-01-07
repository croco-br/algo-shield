package groups

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func Test_NewPostgresRepository_WhenCalled_ThenReturnsRepository(t *testing.T) {
	var db *pgxpool.Pool

	repo := NewPostgresRepository(db)

	assert.NotNil(t, repo)
	assert.Implements(t, (*Repository)(nil), repo)
}

func Test_NewPostgresRepository_WhenCalledWithNilDB_ThenReturnsRepository(t *testing.T) {
	repo := NewPostgresRepository(nil)

	assert.NotNil(t, repo)
	assert.Implements(t, (*Repository)(nil), repo)
}
