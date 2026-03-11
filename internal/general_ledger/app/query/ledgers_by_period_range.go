package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type LedgersByPeriodRangeHandler struct {
	readModel            GeneralLedgerReadModel
	periodRangeValidator periodRangeValidator
}

func NewLedgersByPeriodRangeHandler(readModel GeneralLedgerReadModel) LedgersByPeriodRangeHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return LedgersByPeriodRangeHandler{
		readModel:            readModel,
		periodRangeValidator: newPeriodRangeValidator(readModel),
	}
}

func (h LedgersByPeriodRangeHandler) Handle(ctx context.Context, sobId uuid.UUID, fromPeriod, toPeriod string) ([]Ledger, error) {
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber, err := h.periodRangeValidator.validate(ctx, sobId, fromPeriod, toPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid period range: %w", err)
	}

	ledgers, err := h.readModel.AllLedgersByPeriodRange(ctx, sobId, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ledgers: %w", err)
	}

	return aggregateLedgersByAccount(ledgers), nil
}

// aggregateLedgersByAccount merges per-period ledger rows into one Ledger per account.
// Input must be sorted by account_number asc, then period asc (guaranteed by AllLedgersByPeriodRange).
//
//   - openingAmount: first period's openingAmount
//   - periodDebit/periodCredit/periodAmount: sum across all periods
//   - endingAmount: last period's endingAmount
func aggregateLedgersByAccount(ledgers []Ledger) []Ledger {
	result := make([]Ledger, 0, len(ledgers))

	var current *Ledger
	for _, l := range ledgers {
		if current == nil || current.AccountId != l.AccountId {
			if current != nil {
				result = append(result, *current)
			}
			current = new(l)
		} else {
			current.PeriodAmount = current.PeriodAmount.Add(l.PeriodAmount)
			current.PeriodDebit = current.PeriodDebit.Add(l.PeriodDebit)
			current.PeriodCredit = current.PeriodCredit.Add(l.PeriodCredit)
			current.EndingAmount = l.EndingAmount
		}
	}
	if current != nil {
		result = append(result, *current)
	}

	return result
}
