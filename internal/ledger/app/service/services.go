package service

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"

	"github.com/google/uuid"
)

type AccountService interface {
	ReadSuperiorAccountIds(ctx context.Context, accountId uuid.UUID) ([]uuid.UUID, error)
	ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]query.Account, error)
	ReadAllAccountIdsBySobId(ctx context.Context, sobId uuid.UUID) ([]uuid.UUID, error)
}
