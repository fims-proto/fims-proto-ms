package tenantservice

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"
	tenantport "github/fims-proto/fims-proto-ms/internal/tenant/port/private/intraprocess"

	"github.com/google/uuid"
)

type TenantServiceImpl struct {
	tenantInterface tenantport.TenantInterface
}

func NewTenantService(tenantInterface tenantport.TenantInterface) TenantServiceImpl {
	return TenantServiceImpl{tenantInterface: tenantInterface}
}

func (s TenantServiceImpl) ReadTenantByUUID(ctx context.Context, tenantId uuid.UUID) (query.Tenant, error) {
	return s.tenantInterface.ReadTenantByUUID(ctx, tenantId)
}

func (s TenantServiceImpl) ReadTenantIdBySubdomain(ctx context.Context, subdomain string) (uuid.UUID, error) {
	return s.tenantInterface.ReadTenantIdBySubdomain(ctx, subdomain)
}
