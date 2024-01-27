package query

import (
	"context"
)

type TenantBySubdomainHandler struct {
	readModel TenantReadModel
}

func NewTenantBySubdomain(readModel TenantReadModel) TenantBySubdomainHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return TenantBySubdomainHandler{readModel: readModel}
}

func (h TenantBySubdomainHandler) Handle(ctx context.Context, subdomain string) (Tenant, error) {
	return h.readModel.TenantBySubdomain(ctx, subdomain)
}
