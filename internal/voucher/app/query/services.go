package query

import (
	"context"

	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"
	userQuery "github/fims-proto/fims-proto-ms/internal/user/app/query"

	"github.com/google/uuid"
)

type AccountService interface {
	ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]accountQuery.Account, error)
}

type UserService interface {
	ReadUserByIds(ctx context.Context, userIds []uuid.UUID) (map[uuid.UUID]userQuery.User, error)
}
