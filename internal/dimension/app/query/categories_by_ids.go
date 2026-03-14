package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CategoriesByIdsHandler struct {
	readModel DimensionReadModel
}

func NewCategoriesByIdsHandler(readModel DimensionReadModel) CategoriesByIdsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return CategoriesByIdsHandler{readModel: readModel}
}

func (h CategoriesByIdsHandler) Handle(ctx context.Context, categoryIds []uuid.UUID) ([]DimensionCategory, error) {
	categories, err := h.readModel.CategoriesByIds(ctx, categoryIds)
	if err != nil {
		return nil, fmt.Errorf("failed to read categories by ids: %w", err)
	}

	return categories, nil
}
