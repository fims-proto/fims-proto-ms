package query

import (
	"context"

	"github.com/google/uuid"
)

type CategoryByIdHandler struct {
	readModel DimensionReadModel
}

func NewCategoryByIdHandler(readModel DimensionReadModel) CategoryByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return CategoryByIdHandler{readModel: readModel}
}

func (h CategoryByIdHandler) Handle(ctx context.Context, categoryId uuid.UUID) (DimensionCategory, error) {
	return h.readModel.CategoryById(ctx, categoryId)
}
