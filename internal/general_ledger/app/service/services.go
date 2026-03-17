package service

import (
	"context"

	dimensionQuery "github/fims-proto/fims-proto-ms/internal/dimension/app/query"
	userQuery "github/fims-proto/fims-proto-ms/internal/user/app/query"

	sobQuery "github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github.com/google/uuid"
)

type SobService interface {
	ReadById(ctx context.Context, sobId uuid.UUID) (sobQuery.Sob, error)
}

type NumberingService interface {
	GenerateIdentifier(ctx context.Context, periodId uuid.UUID) (string, error)
	CreateIdentifierConfigurationForJournal(ctx context.Context, periodId uuid.UUID) error
}

type UserService interface {
	ReadUsersByIds(ctx context.Context, userIds []uuid.UUID) (map[uuid.UUID]userQuery.User, error)
}

type DimensionService interface {
	// ValidateOptions validates that optionIds cover all requiredCategoryIds (mandatory one-per-category),
	// with no duplicates and no options from disallowed categories.
	ValidateOptions(ctx context.Context, requiredCategoryIds []uuid.UUID, optionIds []uuid.UUID) error
	FetchCategoriesByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]dimensionQuery.DimensionCategory, error)
	FetchOptionsByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]dimensionQuery.DimensionOption, error)
}
