package tenantmanager

import (
	"context"
	"errors"
	"testing"

	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var tenant1, tenant3 = uuid.New(), uuid.New()

func TestLib_TenantManager_GetDBConn(t *testing.T) {
	t.Parallel()

	tenantManager := NewTenantManager(mockTenantService{}, mockDBConnector{})

	db, err := tenantManager.GetDBConnBySubdomain(context.Background(), "subdomain1")
	assert.NotNil(t, db)
	assert.NoError(t, err)
}

func TestLib_TenantManager_GetDBConn_noTenant(t *testing.T) {
	t.Parallel()

	tenantManager := NewTenantManager(mockTenantService{}, mockDBConnector{})

	db, err := tenantManager.GetDBConnBySubdomain(context.Background(), "subdomain2")
	assert.Nil(t, db)
	assert.Error(t, err)
}

func TestLib_TenantManager_GetDBConn_openConnFailed(t *testing.T) {
	t.Parallel()

	tenantManager := NewTenantManager(mockTenantService{}, mockDBConnector{})

	db, err := tenantManager.GetDBConnBySubdomain(context.Background(), "subdomain3")
	assert.Nil(t, db)
	assert.Error(t, err)
}

type mockTenantService struct{}

func (m mockTenantService) ReadTenantBySubdomain(_ context.Context, subdomain string) (query.Tenant, error) {
	if subdomain == "subdomain1" {
		return query.Tenant{
			TenantId:  tenant1,
			Subdomain: subdomain,
			DSN:       tenant1.String(),
		}, nil
	}
	if subdomain == "subdomain3" {
		return query.Tenant{
			TenantId:  tenant3,
			Subdomain: subdomain,
			DSN:       tenant3.String(),
		}, nil
	}
	// tenant2
	return query.Tenant{}, errors.New("no tenant")
}

type mockDBConnector struct{}

func (m mockDBConnector) Open(dsn string) (*gorm.DB, error) {
	if dsn == tenant1.String() {
		return &gorm.DB{}, nil
	}
	// tenant3
	return nil, errors.New("open connection failed")
}
