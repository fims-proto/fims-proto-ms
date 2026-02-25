package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LedgerSummary struct {
	AccountId     uuid.UUID
	PeriodId      uuid.UUID
	OpeningAmount decimal.Decimal
	PeriodAmount  decimal.Decimal
	PeriodDebit   decimal.Decimal
	PeriodCredit  decimal.Decimal
	EndingAmount  decimal.Decimal
}

type LedgerSummaryHandler struct {
	readModel            GeneralLedgerReadModel
	periodRangeValidator periodRangeValidator
}

func NewLedgerSummaryHandler(readModel GeneralLedgerReadModel) LedgerSummaryHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return LedgerSummaryHandler{
		readModel:            readModel,
		periodRangeValidator: newPeriodRangeValidator(readModel),
	}
}

func (h LedgerSummaryHandler) Handle(ctx context.Context, sobId, accountId uuid.UUID, fromPeriod, toPeriod string) (LedgerSummary, error) {
	// Validate period continuity
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber, err := h.periodRangeValidator.validate(ctx, sobId, fromPeriod, toPeriod)
	if err != nil {
		return LedgerSummary{}, fmt.Errorf("invalid period range: %w", err)
	}

	// Query ledgers with SQL-level period range filtering
	ledgers, err := h.readModel.LedgersByPeriodRange(
		ctx,
		sobId,
		accountId,
		fromFiscalYear,
		fromPeriodNumber,
		toFiscalYear,
		toPeriodNumber,
	)
	if err != nil {
		return LedgerSummary{}, fmt.Errorf("failed to search ledgers: %w", err)
	}

	// Aggregate across all periods
	var sumDebit, sumCredit decimal.Decimal
	firstLedger := ledgers[0]
	lastLedger := ledgers[len(ledgers)-1]

	for _, l := range ledgers {
		sumDebit = sumDebit.Add(l.PeriodDebit)
		sumCredit = sumCredit.Add(l.PeriodCredit)
	}

	return LedgerSummary{
		AccountId:     accountId,
		PeriodId:      lastLedger.PeriodId,
		OpeningAmount: firstLedger.OpeningAmount,
		PeriodAmount:  firstLedger.PeriodAmount,
		PeriodDebit:   sumDebit,
		PeriodCredit:  sumCredit,
		EndingAmount:  firstLedger.EndingAmount,
	}, nil
}
