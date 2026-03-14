package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type PagingCategoriesHandler struct {
	readModel DimensionReadModel
}

func NewPagingCategoriesHandler(readModel DimensionReadModel) PagingCategoriesHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingCategoriesHandler{readModel: readModel}
}

func (h PagingCategoriesHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[DimensionCategory], error) {
	return h.readModel.SearchCategories(ctx, sobId, pageRequest)
}
