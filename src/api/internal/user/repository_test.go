package user

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func Test_NewPostgresUserRepository_WhenCalled_ThenReturnsRepository(t *testing.T) {
	var db *pgxpool.Pool

	repo := NewPostgresUserRepository(db)

	assert.NotNil(t, repo)
	assert.Implements(t, (*UserRepository)(nil), repo)
}

func Test_NewPostgresUserRepository_WhenCalledWithNilDB_ThenReturnsRepository(t *testing.T) {
	repo := NewPostgresUserRepository(nil)

	assert.NotNil(t, repo)
	assert.Implements(t, (*UserRepository)(nil), repo)
}

func Test_NewPostgresTransactionManager_WhenCalled_ThenReturnsTransactionManager(t *testing.T) {
	var db *pgxpool.Pool

	manager := NewPostgresTransactionManager(db)

	assert.NotNil(t, manager)
	assert.Implements(t, (*TransactionManager)(nil), manager)
}

func Test_NewPostgresTransactionManager_WhenCalledWithNilDB_ThenReturnsTransactionManager(t *testing.T) {
	manager := NewPostgresTransactionManager(nil)

	assert.NotNil(t, manager)
	assert.Implements(t, (*TransactionManager)(nil), manager)
}
