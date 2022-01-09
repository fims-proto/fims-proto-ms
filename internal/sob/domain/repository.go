package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateSob(ctx context.Context, sob *Sob) error
	UpdateSob(
		ctx context.Context,
		sobId uuid.UUID,
		updateFn func(s *Sob) (*Sob, error),
	) error
	Migrate(ctx context.Context) error
}
