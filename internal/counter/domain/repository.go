package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateCounter(ctx context.Context, c *Counter) error
	DeleteCounter(ctx context.Context, id uuid.UUID) error
	UpdateCounter(ctx context.Context, id uuid.UUID, updateFn func(c *Counter) (*Counter, error)) error
	UpdateAndRead(ctx context.Context, id uuid.UUID, updateAndReadFn func(c *Counter) (*Counter, interface{}, error)) (interface{}, error)
	Migrate(ctx context.Context) error
}
