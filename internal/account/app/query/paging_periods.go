package query

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type PagingPeriodsHandler struct {
	readModel AccountReadModel
}

func NewPagingPeriodsHandler(readModel AccountReadModel) PagingPeriodsHandler {
	if readModel == nil {
		panic("nil account read model")
	}
	return PagingPeriodsHandler{readModel: readModel}
}

func (h PagingPeriodsHandler) Handle(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Period], error) {
	return h.readModel.PagingPeriods(ctx, sobId, pageable)
}
