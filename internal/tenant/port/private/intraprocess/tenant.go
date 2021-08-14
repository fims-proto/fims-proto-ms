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

func (i TenantInterface) ReadTenantByUUID(ctx context.Context, tenantId uuid.UUID) (query.Tenant, error) {
	return i.app.Queries.ReadTenants.HandleReadByUUID(ctx, tenantId)
}

func (i TenantInterface) ReadTenantIdBySubdomain(ctx context.Context, subdomain string) (uuid.UUID, error) {
	t, err := i.app.Queries.ReadTenants.HandleReadBySubdomain(ctx, subdomain)
	if err != nil {
		return uuid.Nil, err
	}
	return t.TenantId, nil
}
