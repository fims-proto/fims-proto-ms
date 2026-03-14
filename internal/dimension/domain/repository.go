package domain

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/dimension/domain/category"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain/option"

	"github.com/google/uuid"
)

type Repository interface {
	Migrate(ctx context.Context) error
	EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error

	CreateCategory(ctx context.Context, c *category.DimensionCategory) error
	UpdateCategory(
		ctx context.Context,
		categoryId uuid.UUID,
		updateFn func(c *category.DimensionCategory) (*category.DimensionCategory, error),
	) error
	DeleteCategory(ctx context.Context, categoryId uuid.UUID) error
	ReadCategoryById(ctx context.Context, categoryId uuid.UUID) (*category.DimensionCategory, error)
	ExistsCategoryUsedByJournalLine(ctx context.Context, categoryId uuid.UUID) (bool, error)

	CreateOption(ctx context.Context, o *option.DimensionOption) error
	UpdateOption(
		ctx context.Context,
		optionId uuid.UUID,
		updateFn func(o *option.DimensionOption) (*option.DimensionOption, error),
	) error
	DeleteOption(ctx context.Context, optionId uuid.UUID) error
	DeleteOptionsByCategoryId(ctx context.Context, categoryId uuid.UUID) error
	ReadOptionsByIds(ctx context.Context, optionIds []uuid.UUID) ([]*option.DimensionOption, error)
	ExistsOptionUsedByJournalLine(ctx context.Context, optionId uuid.UUID) (bool, error)
}
