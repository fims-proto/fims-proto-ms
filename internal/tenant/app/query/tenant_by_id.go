package query

import (
	"context"

	"github.com/google/uuid"
)

type TenantByIdHandler struct {
	readModel TenantReadModel
}

func NewTenantByIdHandler(readModel TenantReadModel) TenantByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return TenantByIdHandler{readModel: readModel}
}

func (h TenantByIdHandler) Handle(ctx context.Context, tenantId uuid.UUID) (Tenant, error) {
	return h.readModel.TenantById(ctx, tenantId)
}
