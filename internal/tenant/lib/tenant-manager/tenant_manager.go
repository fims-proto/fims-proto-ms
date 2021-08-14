package tenantmanager

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"
	"sync"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type dbConnector interface {
	Open(username, password string) (*sqlx.DB, error)
}

type tenantService interface {
	ReadTenantByUUID(ctx context.Context, tenantId uuid.UUID) (query.Tenant, error)
}

type tenant struct {
	tenantId  uuid.UUID
	subdomain string
	dbConn    *sqlx.DB
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

func (t *TenantManagerImpl) GetDBConn(ctx context.Context, tenantId uuid.UUID) (db *sqlx.DB, err error) {
	defer func() {
		// change returing
		if r := recover(); r != nil {
			db = nil
			err = r.(error)
		}
	}()
	value, _ := t.tenants.LoadOrStore(tenantId, t.loadTenant(ctx, tenantId))
	return value.(*tenant).dbConn, nil
}

func (t *TenantManagerImpl) loadTenant(ctx context.Context, tenantId uuid.UUID) *tenant {
	queriedTenant, err := t.tenantService.ReadTenantByUUID(ctx, tenantId)
	if err != nil {
		panic(errors.Wrap(err, "cannot load tenant"))
	}
	db, err := t.dbConnector.Open(tenantId.String(), queriedTenant.DBConnPassword)
	if err != nil {
		panic(errors.Wrap(err, "open db connection failed"))
	}
	return &tenant{
		tenantId:  queriedTenant.TenantId,
		subdomain: queriedTenant.Subdomain,
		dbConn:    db,
	}
}
