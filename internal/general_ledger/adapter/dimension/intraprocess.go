package dimension

import (
	"context"
	"fmt"

	dimensionQuery "github/fims-proto/fims-proto-ms/internal/dimension/app/query"
	dimensionIntraPort "github/fims-proto/fims-proto-ms/internal/dimension/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	dimensionInterface dimensionIntraPort.DimensionInterface
}

func NewIntraProcessAdapter(dimensionInterface dimensionIntraPort.DimensionInterface) IntraProcessAdapter {
	return IntraProcessAdapter{dimensionInterface: dimensionInterface}
}

func (a IntraProcessAdapter) ValidateOptions(
	ctx context.Context,
	requiredCategoryIds []uuid.UUID,
	optionIds []uuid.UUID,
) error {
	return a.dimensionInterface.ValidateOptions(ctx, requiredCategoryIds, optionIds)
}

func (a IntraProcessAdapter) FetchCategoriesByIds(
	ctx context.Context,
	ids []uuid.UUID,
) (map[uuid.UUID]dimensionQuery.DimensionCategory, error) {
	categories, err := a.dimensionInterface.CategoriesByIds(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dimension categories: %w", err)
	}

	result := make(map[uuid.UUID]dimensionQuery.DimensionCategory, len(categories))
	for _, cat := range categories {
		result[cat.Id] = cat
	}

	return result, nil
}

func (a IntraProcessAdapter) FetchOptionsByIds(
	ctx context.Context,
	ids []uuid.UUID,
) (map[uuid.UUID]dimensionQuery.DimensionOption, error) {
	options, err := a.dimensionInterface.OptionsByIds(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dimension options: %w", err)
	}

	result := make(map[uuid.UUID]dimensionQuery.DimensionOption, len(options))
	for _, opt := range options {
		result[opt.Id] = opt
	}

	return result, nil
}
