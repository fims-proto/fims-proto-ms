package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateSob(ctx context.Context, sob *Sob) error
	UpdateSob(
		ctx context.Context,
		sobUUID uuid.UUID,
		updateFn func(s *Sob) (*Sob, error),
	)
}
