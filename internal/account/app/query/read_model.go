package query

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type AccountReadModel interface {
	SearchAccounts(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Account], error)
	SearchLedgers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Ledger], error)
	SearchPeriods(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Period], error)

	AllAccounts(ctx context.Context, sobId uuid.UUID) ([]Account, error)
	AccountsByIds(ctx context.Context, accountIds []uuid.UUID) ([]Account, error)
	AccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]Account, error)
	SuperiorAccounts(ctx context.Context, accountId uuid.UUID) ([]Account, error)

	LedgersInPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) ([]Ledger, error)
	PagingLedgersByPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID, pageRequest data.PageRequest) (data.Page[Ledger], error)

	OpenPeriod(ctx context.Context, sobId uuid.UUID) (Period, error)
	PeriodById(ctx context.Context, periodId uuid.UUID) (Period, error)
	PeriodByFiscalYearAndNumber(ctx context.Context, sobId uuid.UUID, fiscalYear, periodNumber int) (Period, error)
	PeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (Period, error)
	PeriodsByIds(ctx context.Context, periodIds []uuid.UUID) ([]Period, error)
}
