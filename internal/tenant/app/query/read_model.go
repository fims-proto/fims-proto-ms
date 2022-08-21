package query

import (
	"context"

	"github.com/google/uuid"
)

type TenantReadModel interface {
	TenantById(ctx context.Context, tenantId uuid.UUID) (Tenant, error)
	TenantBySubdomain(ctx context.Context, subdomain string) (Tenant, error)
}
