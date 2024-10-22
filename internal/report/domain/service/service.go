package service

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/report/domain/general_ledger"
)

type GeneralLedgerService interface {
	ReadPeriodIdByFiscalYearAndNumber(ctx context.Context, sobId uuid.UUID, fiscalYear, number int) (uuid.UUID, error)
	ReadPeriodById(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) (*general_ledger.Period, error)
	ReadFirstPeriodOfTheYear(ctx context.Context, sobId uuid.UUID, fiscalYear int) (*general_ledger.Period, error)

	ReadAccountIdsByNumbers(ctx context.Context, sobId uuid.UUID, accountNumbers []string) (map[string]uuid.UUID, error)

	ReadLedgersByAccountAndPeriodsOrderByPeriod(
		ctx context.Context,
		sobId uuid.UUID,
		accountId uuid.UUID,
		periods []*general_ledger.Period,
	) ([]*general_ledger.Ledger, error)
}
