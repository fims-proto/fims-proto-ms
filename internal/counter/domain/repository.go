package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// maybe someday, reseting formatter in Counter is necessary
	CreateCounter(ctx context.Context, c *Counter) error
	DeleteCounter(ctx context.Context, counterUUID uuid.UUID) error

	UpdateCounter(
		ctx context.Context,
		counterUUID uuid.UUID,
		updateFn func(c *Counter) (*Counter, error),
	) error

	// I'll have to compromise here
	// So we can guarantee counter.Next() and counter.Identifier() are in same trasaction
	// and keep business logic away from repository implementation
	UpdateAndRead(
		ctx context.Context,
		counterUUID uuid.UUID,
		// 3 returns to repository impl: 1 updated counter to save in DB, 2 any value that needs to be read, 3 errors if any
		updateAndReadFn func(c *Counter) (*Counter, interface{}, error),
	) (interface{}, error) // 2 returns to command/query/service: 1 any value that needs to be read, 2 errors if any
}
