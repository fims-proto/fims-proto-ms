package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	UpsertUser(ctx context.Context, id uuid.UUID, updateFn func(*User) (*User, error)) error
	Migrate(ctx context.Context) error
}
