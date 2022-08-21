package query

import (
	"context"

	"github.com/google/uuid"
)

type UserReadModel interface {
	UserById(ctx context.Context, id uuid.UUID) (User, error)
	UsersByIds(ctx context.Context, ids []uuid.UUID) ([]User, error)
}
