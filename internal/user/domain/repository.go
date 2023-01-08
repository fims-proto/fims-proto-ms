package domain

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/user/domain/user"

	"github.com/google/uuid"
)

type Repository interface {
	UpsertUser(ctx context.Context, userId uuid.UUID, updateFn func(*user.User) (*user.User, error)) error

	Migrate(ctx context.Context) error
}
