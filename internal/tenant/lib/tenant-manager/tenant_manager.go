package tenantmanager

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type dbConnector interface {
	Open(dsn string) (*gorm.DB, error)
}

type tenantService interface {
	ReadTenantBySubdomain(ctx context.Context, subdomain string) (query.Tenant, error)
}

type tenant struct {
	tenantId  uuid.UUID
	subdomain string
	dbConn    *gorm.DB
}

type syncData struct {
	data *tenant
	once *sync.Once
}

type TenantManagerImpl struct {
	tenants       *sync.Map
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
		tenants:       &sync.Map{},
		tenantService: tenantService,
		dbConnector:   dbConnector,
	}
}

func (t *TenantManagerImpl) GetDBConnBySubdomain(ctx context.Context, subdomain string) (*gorm.DB, error) {
	if subdomain == "" {
		return nil, errors.New("empty subdomain")
	}

	value, err := t.loadOrStore(ctx, subdomain)
	if err != nil {
		return nil, err
	}

	return value.dbConn, nil
}

// loadOrStore get tenant by subdomain if sync.Map has the value, other wise compute
func (t *TenantManagerImpl) loadOrStore(ctx context.Context, subdoamin string) (value *tenant, err error) {
	actual, _ := t.tenants.LoadOrStore(subdoamin, &syncData{
		data: nil,
		once: &sync.Once{},
	})

	d := actual.(*syncData)
	if d.data == nil {
		d.once.Do(func() {
			d.data, err = t.initiateTenant(ctx, subdoamin)
			if err != nil {
				// if failed, reset once
				d.once = &sync.Once{}
				log.Err(ctx, err, "failed to load tenant")
			}
		})
	}
	return d.data, err
}

func (t *TenantManagerImpl) initiateTenant(ctx context.Context, subdomain string) (*tenant, error) {
	queriedTenant, err := t.tenantService.ReadTenantBySubdomain(ctx, subdomain)
	if err != nil {
		return nil, errors.Wrap(err, "cannot load tenant")
	}

	log.Debug(ctx, "trying to open connection for schema %s", queriedTenant.TenantId.String())

	// DB connection
	db, err := t.dbConnector.Open(queriedTenant.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "open db connection failed")
	}

	return &tenant{
		tenantId:  queriedTenant.TenantId,
		subdomain: queriedTenant.Subdomain,
		dbConn:    db,
	}, nil
}
