package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Period struct {
	Id           uuid.UUID
	FiscalYear   int
	PeriodNumber int
}

type Ledger struct {
	AccountId        uuid.UUID
	BalanceDirection string
	Period           Period
	OpeningAmount    decimal.Decimal
	PeriodAmount     decimal.Decimal
	PeriodDebit      decimal.Decimal
	PeriodCredit     decimal.Decimal
	EndingAmount     decimal.Decimal
}

type GeneralLedgerService interface {
	ReadPeriodIdByFiscalYearAndNumber(ctx context.Context, sobId uuid.UUID, fiscalYear int, number int) (uuid.UUID, error)
	ReadPeriodById(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) (*Period, error)
	ReadFirstPeriodOfTheYear(ctx context.Context, sobId uuid.UUID, fiscalYear int) (*Period, error)
	ReadAccountIdsByRawNumbers(ctx context.Context, sobId uuid.UUID, rawAccountNumbers []string) (map[string]uuid.UUID, error)
	ReadCashFlowItemIdsByCodes(ctx context.Context, sobId uuid.UUID, codes []string) (map[string]uuid.UUID, error)

	ReadLedgersByAccountAndPeriodsOrderByPeriod(
		ctx context.Context,
		sobId uuid.UUID,
		accountId uuid.UUID,
		periods []*Period,
	) ([]Ledger, error)

	SumAbsJournalLineAmountsByCashFlowItemsAndPeriods(
		ctx context.Context,
		sobId uuid.UUID,
		cashFlowItemIds []uuid.UUID,
		periods []*Period,
	) (map[uuid.UUID]decimal.Decimal, error)
}
