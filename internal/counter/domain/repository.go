package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// maybe someday, reseting formatter in Counter is necessary
	CreateCounter(ctx context.Context, c *Counter) error
	DeleteCounter(ctx context.Context, id uuid.UUID) error

	UpdateCounter(ctx context.Context, id uuid.UUID, updateFn func(c *Counter) (*Counter, error)) error

	// I'll have to compromise here
	// So we can guarantee counter.Next() and counter.Identifier() are in same trasaction
	// and keep business logic away from repository implementation
	UpdateAndRead(ctx context.Context, id uuid.UUID, updateAndReadFn func(c *Counter) (*Counter, interface{}, error)) (interface{}, error)

	Migrate(ctx context.Context) error
}
