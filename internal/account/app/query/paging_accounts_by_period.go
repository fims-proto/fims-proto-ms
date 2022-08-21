package query

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type PagingAccountsByPeriodHandler struct {
	readModel AccountReadModel
}

func NewPagingAccountsByPeriodHandler(readModel AccountReadModel) PagingAccountsByPeriodHandler {
	if readModel == nil {
		panic("nil account read model")
	}

	return PagingAccountsByPeriodHandler{readModel: readModel}
}

func (h PagingAccountsByPeriodHandler) Handle(ctx context.Context, sobId, periodId uuid.UUID, pageable data.Pageable) (data.Page[Account], error) {
	return h.readModel.PagingAccountsByPeriod(ctx, sobId, periodId, pageable)
}
