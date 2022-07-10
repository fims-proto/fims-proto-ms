package query

import (
	"context"

	ledgerQuery "github/fims-proto/fims-proto-ms/internal/ledger/app/query"

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

type LedgerService interface {
	ReadPeriodsByIds(ctx context.Context, periodIds []uuid.UUID) (map[uuid.UUID]ledgerQuery.Period, error)
}
