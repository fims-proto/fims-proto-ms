package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type AccountReadModel interface {
	AllAccounts(ctx context.Context, sobId uuid.UUID) ([]Account, error)
	PagingAccounts(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Account], error)
	AccountsByIds(ctx context.Context, accountIds []uuid.UUID) ([]Account, error)
	AccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]Account, error)
	SuperiorAccounts(ctx context.Context, accountId uuid.UUID) ([]Account, error)

	LedgersInPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) ([]Ledger, error)
	PagingLedgersByPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID, pageable data.Pageable) (data.Page[Ledger], error)

	PagingPeriods(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Period], error)
	PeriodById(ctx context.Context, periodId uuid.UUID) (Period, error)
	PeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (Period, error)
	PeriodsByIds(ctx context.Context, periodIds []uuid.UUID) ([]Period, error)
}
