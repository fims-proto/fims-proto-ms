package tenantmanager

import (
	"context"
	"errors"
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var tenant1, tenant3 uuid.UUID = uuid.New(), uuid.New()

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

func (m mockTenantService) ReadTenantBySubdomain(ctx context.Context, subdomain string) (query.Tenant, error) {
	if subdomain == "subdomain1" {
		return query.Tenant{
			TenantId:       tenant1,
			Subdomain:      subdomain,
			DBConnPassword: "password",
		}, nil
	}
	if subdomain == "subdomain3" {
		return query.Tenant{
			TenantId:       tenant3,
			Subdomain:      subdomain,
			DBConnPassword: "password",
		}, nil
	}
	// tenant2
	return query.Tenant{}, errors.New("no tenant")
}

type mockDBConnector struct{}

func (m mockDBConnector) Open(username, password string) (*gorm.DB, error) {
	if username == tenant1.String() {
		return &gorm.DB{}, nil
	}
	// tenant3
	return nil, errors.New("open connection failed")
}
