package service

import (
	"context"

	sobQuery "github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github.com/google/uuid"
)

type SobService interface {
	ReadById(ctx context.Context, sobId uuid.UUID) (sobQuery.Sob, error)
}

type NumberingService interface {
	InitializeIdentifierConfigurationForVoucher(ctx context.Context, periodId uuid.UUID) error
}
