package db

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestAdapter_PostgresRepository_ReadByUUID(t *testing.T) {
	t.Parallel()

	// GIVEN
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 mockDB,
		PreferSimpleProtocol: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	tenandId := uuid.New()

	row := sqlmock.NewRows([]string{"tenant_id", "subdomain", "db_conn_password"}).
		AddRow(tenandId.String(), "local_test", "password")
	mock.ExpectQuery("SELECT *").WithArgs(tenandId).WillReturnRows(row)

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

func TestAdapter_PostgresRepository_ReadBySubdomain(t *testing.T) {
	t.Parallel()

	// GIVEN
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 mockDB,
		PreferSimpleProtocol: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	row := sqlmock.NewRows([]string{"tenant_id", "subdomain", "db_conn_password"}).
		AddRow(uuid.New().String(), "local_test", "password")
	mock.ExpectQuery("SELECT *").WithArgs("local_domain").WillReturnRows(row)

	// WHEN
	repo := NewTenantPostgresRepository(db)
	tenant, err := repo.ReadBySubdomain(context.Background(), "local_domain")

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, "local_test", tenant.Subdomain)
	assert.Equal(t, "password", tenant.DBConnPassword)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
