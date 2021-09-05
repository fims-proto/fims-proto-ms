package tenantservice

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"
	tenantport "github/fims-proto/fims-proto-ms/internal/tenant/port/private/intraprocess"
)

type TenantServiceImpl struct {
	tenantInterface tenantport.TenantInterface
}

func NewTenantService(tenantInterface tenantport.TenantInterface) TenantServiceImpl {
	return TenantServiceImpl{tenantInterface: tenantInterface}
}

func (s TenantServiceImpl) ReadTenantBySubdomain(ctx context.Context, subdomain string) (query.Tenant, error) {
	return s.tenantInterface.ReadTenantBySubdomain(ctx, subdomain)
}
