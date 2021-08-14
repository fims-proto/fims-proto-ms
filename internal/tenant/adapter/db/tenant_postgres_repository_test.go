package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAdapter_PostgresRepository_ReadByUUID(t *testing.T) {
	t.Parallel()

	// GIVEN
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	db := sqlx.NewDb(mockDB, "sqlmock")

	tenandId := uuid.New()

	row := sqlmock.NewRows([]string{"tenant_id", "subdomain", "db_conn_password"}).
		AddRow(tenandId.String(), "local_test", "password")
	mock.ExpectQuery("SELECT (.+) FROM tenant WHERE tenant_id = ?").WithArgs(tenandId).WillReturnRows(row)

	// WHEN
	repo := NewTenantPostgresRepository(db)
	tenant, err := repo.ReadByUUID(context.Background(), tenandId)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, "local_test", tenant.Subdomain)
	assert.Equal(t, "password", tenant.DBConnPassword)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
