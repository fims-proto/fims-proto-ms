package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type GeneralLedgerReadModel interface {
	SearchAccounts(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Account], error)
	SearchLedgers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Ledger], error)
	SearchPeriods(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Period], error)
	SearchVouchers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Voucher], error)

	AllAccounts(ctx context.Context, sobId uuid.UUID) ([]Account, error)
	AccountsByIds(ctx context.Context, accountIds []uuid.UUID) ([]Account, error)
	AccountsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]Account, error)
	SuperiorAccounts(ctx context.Context, accountId uuid.UUID) ([]Account, error)

	LedgersInPeriod(ctx context.Context, sobId, periodId uuid.UUID) ([]Ledger, error)
	PagingLedgersByPeriod(ctx context.Context, sobId, periodId uuid.UUID, pageRequest data.PageRequest) (data.Page[Ledger], error)
	FirstLevelLedgersInPeriod(ctx context.Context, sobId, periodId uuid.UUID) ([]Ledger, error)
	ExistsProfitAndLossLedgersHavingBalanceInPeriod(ctx context.Context, sobId, periodId uuid.UUID) (bool, error)

	CurrentPeriod(ctx context.Context, sobId uuid.UUID) (Period, error)
	PeriodById(ctx context.Context, periodId uuid.UUID) (Period, error)
	PeriodsByIds(ctx context.Context, periodIds []uuid.UUID) ([]Period, error)
	PeriodByFiscalYearAndNumber(ctx context.Context, sobId uuid.UUID, fiscalYear, periodNumber int) (Period, error)

	VoucherById(ctx context.Context, voucherId uuid.UUID) (Voucher, error)
	ExistsVouchersNotPostedInPeriod(ctx context.Context, sobId, periodId uuid.UUID) (bool, error)
}
