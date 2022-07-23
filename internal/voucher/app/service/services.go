package service

import (
	"context"
	"time"

	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"
	ledgerQuery "github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	userQuery "github/fims-proto/fims-proto-ms/internal/user/app/query"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
)

type LedgerService interface {
	ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, transactionTime time.Time) (ledgerQuery.Period, error)
	ReadPeriodsByIds(ctx context.Context, periodIds []uuid.UUID) (map[uuid.UUID]ledgerQuery.Period, error)
	PostVoucher(ctx context.Context, voucher domain.Voucher) error
}

type AccountService interface {
	ValidateExistenceAndGetId(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]uuid.UUID, error)
	ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]accountQuery.Account, error)
}

type UserService interface {
	ReadUsersByIds(ctx context.Context, userIds []uuid.UUID) (map[uuid.UUID]userQuery.User, error)
}

type NumberingService interface {
	GenerateIdentifier(ctx context.Context, periodId uuid.UUID, voucherType string) (string, error)
}
