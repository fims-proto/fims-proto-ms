package tenantmanager

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type dbConnector interface {
	Open(username, password string) (*gorm.DB, error)
}

type tenantService interface {
	ReadTenantBySubdomain(ctx context.Context, subdomain string) (query.Tenant, error)
}

type tenant struct {
	tenantId  uuid.UUID
	subdomain string
	dbConn    *gorm.DB
}

type TenantManagerImpl struct {
	tenants       sync.Map
	tenantService tenantService
	dbConnector   dbConnector
}

func NewTenantManager(tenantService tenantService, dbConnector dbConnector) *TenantManagerImpl {
	if tenantService == nil {
		panic("nil tenant service")
	}
	if dbConnector == nil {
		panic("nil dbConnector")
	}
	return &TenantManagerImpl{
		tenants:       sync.Map{},
		tenantService: tenantService,
		dbConnector:   dbConnector,
	}
}

func (t *TenantManagerImpl) GetDBConnBySubdomain(ctx context.Context, subdomain string) (db *gorm.DB, err error) {
	defer func() {
		// change returing
		if r := recover(); r != nil {
			db = nil
			err = r.(error)
		}
	}()

	if subdomain == "" {
		return nil, errors.New("empty subdomain")
	}

	value, _ := t.tenants.LoadOrStore(subdomain, t.loadTenant(ctx, subdomain))
	return value.(*tenant).dbConn, nil
}

func (t *TenantManagerImpl) loadTenant(ctx context.Context, subdomain string) *tenant {
	queriedTenant, err := t.tenantService.ReadTenantBySubdomain(ctx, subdomain)
	if err != nil {
		panic(errors.Wrap(err, "cannot load tenant"))
	}
	db, err := t.dbConnector.Open(queriedTenant.TenantId.String(), queriedTenant.DBConnPassword)
	if err != nil {
		panic(errors.Wrap(err, "open db connection failed"))
	}
	return &tenant{
		tenantId:  queriedTenant.TenantId,
		subdomain: queriedTenant.Subdomain,
		dbConn:    db,
	}
}
