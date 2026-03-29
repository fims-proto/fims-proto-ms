package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github.com/google/uuid"
)

type LedgersByPeriodRangeHandler struct {
	readModel            GeneralLedgerReadModel
	sobService           service.SobService
	periodRangeValidator periodRangeValidator
}

func NewLedgersByPeriodRangeHandler(readModel GeneralLedgerReadModel, sobService service.SobService) LedgersByPeriodRangeHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	return LedgersByPeriodRangeHandler{
		readModel:            readModel,
		sobService:           sobService,
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

	aggregated := aggregateLedgersByAccount(ledgers)

	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return nil, fmt.Errorf("failed to read sob: %w", err)
	}

	return enrichLedgerAccountNumbers(sob.AccountsCodeLength, aggregated), nil
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
