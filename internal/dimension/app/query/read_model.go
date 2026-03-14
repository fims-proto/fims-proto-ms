package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type DimensionReadModel interface {
	SearchCategories(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[DimensionCategory], error)
	SearchOptions(ctx context.Context, categoryId uuid.UUID, pageRequest data.PageRequest) (data.Page[DimensionOption], error)
	CategoryById(ctx context.Context, categoryId uuid.UUID) (DimensionCategory, error)
	CategoriesByIds(ctx context.Context, categoryIds []uuid.UUID) ([]DimensionCategory, error)
	OptionsByIds(ctx context.Context, optionIds []uuid.UUID) ([]DimensionOption, error)
}
