package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
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

func (h PagingPeriodsHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Period], error) {
	return h.readModel.SearchPeriods(ctx, sobId, pageRequest)
}
