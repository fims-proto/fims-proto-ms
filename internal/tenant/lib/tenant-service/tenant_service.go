package tenantservice

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"

	tenantPort "github/fims-proto/fims-proto-ms/internal/tenant/port/private/intraprocess"
)

type Impl struct {
	tenantInterface tenantPort.TenantInterface
}

func NewTenantService(tenantInterface tenantPort.TenantInterface) Impl {
	return Impl{tenantInterface: tenantInterface}
}

func (s Impl) ReadTenantBySubdomain(ctx context.Context, subdomain string) (query.Tenant, error) {
	return s.tenantInterface.ReadTenantBySubdomain(ctx, subdomain)
}
