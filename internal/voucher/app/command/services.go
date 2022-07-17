package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	ledgerQuery "github/fims-proto/fims-proto-ms/internal/ledger/app/query"

	"github.com/google/uuid"
)

type LedgerService interface {
	ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, transactionTime time.Time) (ledgerQuery.Period, error)
	PostVoucher(ctx context.Context, voucher domain.Voucher) error
}

type CounterService interface {
	GetNextIdentifier(ctx context.Context, businessObjects ...string) (string, error)
}

type AccountService interface {
	ValidateExistenceAndGetId(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]uuid.UUID, error)
}
