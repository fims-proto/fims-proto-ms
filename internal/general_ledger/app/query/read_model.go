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
	SearchJournals(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Journal], error)

	LedgersByPeriodRange(ctx context.Context, sobId, accountId uuid.UUID, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int) ([]Ledger, error)
	AllLedgersByPeriodRange(ctx context.Context, sobId uuid.UUID, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int) ([]Ledger, error)
	LedgerEntriesByPeriodRange(ctx context.Context, sobId, accountId uuid.UUID, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int, pageRequest data.PageRequest) (data.Page[LedgerEntry], error)

	CurrentPeriod(ctx context.Context, sobId uuid.UUID) (Period, error)
	FirstPeriod(ctx context.Context, sobId uuid.UUID) (Period, error)
	CheckPeriodContinuity(ctx context.Context, sobId uuid.UUID, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber int) error

	JournalById(ctx context.Context, journalId uuid.UUID) (Journal, error)
}
