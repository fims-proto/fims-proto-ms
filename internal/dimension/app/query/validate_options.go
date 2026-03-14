package query

import (
	"context"
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

// ValidateOptionsHandler validates that a set of dimension option IDs is valid for a given
// set of required category IDs. It enforces:
//  1. All option IDs must exist.
//  2. Only one option per category.
//  3. If requiredCategoryIds is non-empty, exactly one option must be provided per category.
//  4. Options must not belong to categories outside requiredCategoryIds.
type ValidateOptionsHandler struct {
	readModel DimensionReadModel
}

func NewValidateOptionsHandler(readModel DimensionReadModel) ValidateOptionsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return ValidateOptionsHandler{readModel: readModel}
}

// Handle validates optionIds against requiredCategoryIds.
// requiredCategoryIds: the categories bound to the account (must all be covered).
// optionIds: the options provided on the journal line.
func (h ValidateOptionsHandler) Handle(
	ctx context.Context,
	requiredCategoryIds []uuid.UUID,
	optionIds []uuid.UUID,
) error {
	// If no required categories and no options provided, nothing to validate.
	if len(requiredCategoryIds) == 0 && len(optionIds) == 0 {
		return nil
	}

	// If account has no dimension categories but options are provided, reject.
	if len(requiredCategoryIds) == 0 && len(optionIds) > 0 {
		return commonErrors.NewSlugError("journalLine-disallowedDimensionCategory")
	}

	// If options are provided, fetch them to validate.
	if len(optionIds) == 0 && len(requiredCategoryIds) > 0 {
		return commonErrors.NewSlugError("journalLine-missingRequiredDimensionCategory")
	}

	options, err := h.readModel.OptionsByIds(ctx, optionIds)
	if err != nil {
		return fmt.Errorf("failed to fetch dimension options: %w", err)
	}

	if len(options) != len(optionIds) {
		return commonErrors.NewSlugError("journalLine-invalidDimensionOption")
	}

	// Build a set of allowed category IDs.
	allowedCategories := make(map[uuid.UUID]bool, len(requiredCategoryIds))
	for _, catId := range requiredCategoryIds {
		allowedCategories[catId] = true
	}

	// Track which categories are covered by the provided options.
	coveredCategories := make(map[uuid.UUID]bool, len(options))
	for _, opt := range options {
		if !allowedCategories[opt.CategoryId] {
			return commonErrors.NewSlugError("journalLine-disallowedDimensionCategory")
		}

		if coveredCategories[opt.CategoryId] {
			return commonErrors.NewSlugError("journalLine-duplicateDimensionCategory")
		}

		coveredCategories[opt.CategoryId] = true
	}

	// All required categories must be covered.
	for _, catId := range requiredCategoryIds {
		if !coveredCategories[catId] {
			return commonErrors.NewSlugError("journalLine-missingRequiredDimensionCategory")
		}
	}

	return nil
}
