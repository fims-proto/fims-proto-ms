package intraprocess

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/tenant/app"
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"

	"github.com/google/uuid"
)

type TenantInterface struct {
	app *app.Application
}

func NewTenantInterface(app *app.Application) TenantInterface {
	return TenantInterface{app: app}
}

func (i TenantInterface) ReadTenantById(ctx context.Context, tenantId uuid.UUID) (query.Tenant, error) {
	return i.app.Queries.TenantById.Handle(ctx, tenantId)
}

func (i TenantInterface) ReadTenantBySubdomain(ctx context.Context, subdomain string) (query.Tenant, error) {
	return i.app.Queries.TenantBySubdomain.Handle(ctx, subdomain)
}
