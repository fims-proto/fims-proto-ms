package lib

import (
	"context"
	"errors"
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var tenant1, tenant2, tenant3 uuid.UUID = uuid.New(), uuid.New(), uuid.New()

func TestLib_TenantManager_GetDBConn(t *testing.T) {
	t.Parallel()

	tenantManager := NewTenantManager(mockTenantService{}, mockDBConnector{})

	db, err := tenantManager.GetDBConn(context.Background(), tenant1)
	assert.NotNil(t, db)
	assert.NoError(t, err)
}

func TestLib_TenantManager_GetDBConn_noTenant(t *testing.T) {
	t.Parallel()

	tenantManager := NewTenantManager(mockTenantService{}, mockDBConnector{})

	db, err := tenantManager.GetDBConn(context.Background(), tenant2)
	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestLib_TenantManager_GetDBConn_openConnFailed(t *testing.T) {
	t.Parallel()

	tenantManager := NewTenantManager(mockTenantService{}, mockDBConnector{})

	db, err := tenantManager.GetDBConn(context.Background(), tenant3)
	assert.Nil(t, db)
	assert.Error(t, err)
}

type mockTenantService struct{}

func (m mockTenantService) ReadTenantByUUID(ctx context.Context, tenantId uuid.UUID) (query.Tenant, error) {
	if tenantId == tenant1 || tenantId == tenant3 {
		return query.Tenant{
			TenantId:       tenantId,
			Subdomain:      "tenant",
			DBConnPassword: "password",
		}, nil
	}
	// tenant2
	return query.Tenant{}, errors.New("no tenant")
}

type mockDBConnector struct{}

func (m mockDBConnector) Open(username, password string) (*sqlx.DB, error) {
	if username == tenant1.String() {
		mockDB, _, _ := sqlmock.New()
		defer mockDB.Close()
		return sqlx.NewDb(mockDB, "sqlmock"), nil
	}
	// tenant3
	return nil, errors.New("open connection failed")
}
