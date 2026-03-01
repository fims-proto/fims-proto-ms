package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type GeneralLedgerReadModel interface {
	SearchAccounts(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Account], error)
	SearchAuxiliaryCategories(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[AuxiliaryCategory], error)
	SearchAuxiliaryAccounts(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[AuxiliaryAccount], error)
	SearchLedgers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Ledger], error)
	SearchAuxiliaryLedgers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[AuxiliaryLedger], error)
	SearchPeriods(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Period], error)
	SearchVouchers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Voucher], error)

	LedgersByPeriodRange(ctx context.Context, sobId, accountId uuid.UUID, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int) ([]Ledger, error)
	AllLedgersByPeriodRange(ctx context.Context, sobId uuid.UUID, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int) ([]Ledger, error)
	AuxiliariesByPeriodRange(ctx context.Context, sobId, accountId, auxiliaryCategoryId uuid.UUID, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int, pageRequest data.PageRequest) (data.Page[AuxiliaryLedger], error)
	LedgerEntriesByPeriodRange(ctx context.Context, sobId, accountId uuid.UUID, auxiliaryAccountId *uuid.UUID, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int, pageRequest data.PageRequest) (data.Page[LedgerEntry], error)

	CurrentPeriod(ctx context.Context, sobId uuid.UUID) (Period, error)
	FirstPeriod(ctx context.Context, sobId uuid.UUID) (Period, error)
	CheckPeriodContinuity(ctx context.Context, sobId uuid.UUID, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int) error

	VoucherById(ctx context.Context, voucherId uuid.UUID) (Voucher, error)
}
