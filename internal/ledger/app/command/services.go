package command

import (
	"context"

	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"

	"github.com/google/uuid"
)

type AccountService interface {
	ReadSuperiorAccountIds(ctx context.Context, accountId uuid.UUID) ([]uuid.UUID, error)
	ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]accountQuery.Account, error)
	ReadAllAccountsBySobId(ctx context.Context, sobId uuid.UUID) ([]accountQuery.Account, error)
}
