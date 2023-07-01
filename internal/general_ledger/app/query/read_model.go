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

	PagingLedgersByPeriod(ctx context.Context, sobId, periodId uuid.UUID, pageRequest data.PageRequest) (data.Page[Ledger], error)

	CurrentPeriod(ctx context.Context, sobId uuid.UUID) (Period, error)

	VoucherById(ctx context.Context, voucherId uuid.UUID) (Voucher, error)
}
