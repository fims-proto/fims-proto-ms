package domain

import (
	"context"
)

type Repository interface {
	CreateSob(ctx context.Context, sob *Sob) error
	UpdateSob(
		ctx context.Context,
		sobId string,
		updateFn func(s *Sob) (*Sob, error),
	) error
	Migrate(ctx context.Context) error
}
