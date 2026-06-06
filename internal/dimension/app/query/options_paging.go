package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type PagingOptionsHandler struct {
	readModel DimensionReadModel
}

func NewPagingOptionsHandler(readModel DimensionReadModel) PagingOptionsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingOptionsHandler{readModel: readModel}
}

func (h PagingOptionsHandler) Handle(ctx context.Context, categoryId uuid.UUID, pageRequest data.PageRequest) (data.Page[DimensionOption], error) {
	return h.readModel.SearchOptions(ctx, categoryId, pageRequest)
}
