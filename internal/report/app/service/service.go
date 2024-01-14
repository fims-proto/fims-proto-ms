package service

import (
	"context"

	userQuery "github/fims-proto/fims-proto-ms/internal/user/app/query"

	sobQuery "github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github.com/google/uuid"
)

type SobService interface {
	ReadById(ctx context.Context, sobId uuid.UUID) (sobQuery.Sob, error)
}

// type NumberingService interface {
// 	// GenerateIdentifier(ctx context.Context, periodId uuid.UUID, voucherType string) (string, error)
// 	// CreateIdentifierConfigurationForVoucher(ctx context.Context, periodId uuid.UUID) error
// }

type UserService interface {
	ReadUsersByIds(ctx context.Context, userIds []uuid.UUID) (map[uuid.UUID]userQuery.User, error)
}
