package intraprocess

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/dimension/app"
	"github/fims-proto/fims-proto-ms/internal/dimension/app/query"

	"github.com/google/uuid"
)

type DimensionInterface struct {
	app *app.Application
}

func NewDimensionInterface(app *app.Application) DimensionInterface {
	if app == nil {
		panic("nil dimension app")
	}

	return DimensionInterface{app: app}
}

// ValidateOptions validates that the provided optionIds satisfy all requiredCategoryIds.
// It enforces mandatory tagging, one-option-per-category, and disallowed category rules.
func (i DimensionInterface) ValidateOptions(
	ctx context.Context,
	requiredCategoryIds []uuid.UUID,
	optionIds []uuid.UUID,
) error {
	return i.app.Queries.ValidateOptions.Handle(ctx, requiredCategoryIds, optionIds)
}

func (i DimensionInterface) CategoriesByIds(
	ctx context.Context,
	categoryIds []uuid.UUID,
) ([]query.DimensionCategory, error) {
	return i.app.Queries.CategoriesByIds.Handle(ctx, categoryIds)
}

func (i DimensionInterface) OptionsByIds(
	ctx context.Context,
	optionIds []uuid.UUID,
) ([]query.DimensionOption, error) {
	return i.app.Queries.OptionsByIds.Handle(ctx, optionIds)
}
