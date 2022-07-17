package command

import (
	"context"

	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"

	"github.com/google/uuid"
)

type SelfService interface {
	CreateLedgersForPeriod(ctx context.Context, periodId uuid.UUID) error
}

type AccountService interface {
	ReadAccountsBySobId(ctx context.Context, sobId uuid.UUID) ([]accountQuery.Account, error)
	ReadAccountsWithSuperiorsByIds(ctx context.Context, accountIds []uuid.UUID) ([]accountQuery.Account, error)
}
