package query

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type PagingAccountConfigurationsHandler struct {
	readModel AccountReadModel
}

func NewPagingAccountConfigurationsHandler(readModel AccountReadModel) PagingAccountConfigurationsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingAccountConfigurationsHandler{readModel: readModel}
}

func (h PagingAccountConfigurationsHandler) Handle(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[AccountConfiguration], error) {
	return h.readModel.PagingAccountConfigurations(ctx, sobId, pageable)
}
