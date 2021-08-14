package query

import (
	"context"

	"github.com/google/uuid"
)

type TenantsReadModel interface {
	ReadByUUID(ctx context.Context, tenantId uuid.UUID) (Tenant, error)
	ReadBySubdomain(ctx context.Context, subdomain string) (Tenant, error)
}

type ReadTenantsHandler struct {
	readModel TenantsReadModel
}

func NewReadTenantsHandler(readModel TenantsReadModel) ReadTenantsHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return ReadTenantsHandler{readModel: readModel}
}

func (h ReadTenantsHandler) HandleReadByUUID(ctx context.Context, tenantId uuid.UUID) (Tenant, error) {
	return h.readModel.ReadByUUID(ctx, tenantId)
}

func (h ReadTenantsHandler) HandleReadBySubdomain(ctx context.Context, subdomain string) (Tenant, error) {
	return h.readModel.ReadBySubdomain(ctx, subdomain)
}
