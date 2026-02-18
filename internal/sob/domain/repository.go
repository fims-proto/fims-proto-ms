package domain

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/sob/domain/sob"

	"github.com/google/uuid"
)

type Repository interface {
	EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error

	CreateSob(ctx context.Context, sob *sob.Sob) error
	UpdateSob(ctx context.Context, sobId uuid.UUID, updateFn func(s *sob.Sob) (*sob.Sob, error)) error

	Migrate(ctx context.Context) error
}
