package service

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github.com/google/uuid"
)

type SobService interface {
	ReadById(ctx context.Context, sobId uuid.UUID) (query.Sob, error)
}
