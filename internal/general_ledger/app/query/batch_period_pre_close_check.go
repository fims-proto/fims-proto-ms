package query

import (
	"context"
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

type BatchPreCloseCheckResult struct {
	UnpostedJournals PreCloseCheckUnpostedJournals
	TrialBalance     PreCloseCheckTrialBalance
}

type BatchPeriodPreCloseCheckHandler struct {
	readModel GeneralLedgerReadModel
}

func NewBatchPeriodPreCloseCheckHandler(readModel GeneralLedgerReadModel) BatchPeriodPreCloseCheckHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return BatchPeriodPreCloseCheckHandler{readModel: readModel}
}

func (h BatchPeriodPreCloseCheckHandler) Handle(ctx context.Context, sobId uuid.UUID, targetYear, targetMonth int) (BatchPreCloseCheckResult, error) {
	currentPeriod, err := h.readModel.CurrentPeriod(ctx, sobId)
	if err != nil {
		return BatchPreCloseCheckResult{}, commonErrors.NewInvalidInputError(commonErrors.SlugPeriodNotFound)
	}

	if targetYear < currentPeriod.FiscalYear || (targetYear == currentPeriod.FiscalYear && targetMonth < currentPeriod.PeriodNumber) {
		return BatchPreCloseCheckResult{}, commonErrors.NewInvalidInputError(commonErrors.SlugPeriodBatchCloseTargetInPast)
	}

	existingPeriods, err := h.readModel.PeriodsInRange(ctx, sobId, currentPeriod.FiscalYear, currentPeriod.PeriodNumber, targetYear, targetMonth)
	if err != nil {
		return BatchPreCloseCheckResult{}, fmt.Errorf("failed to list periods in range: %w", err)
	}

	aggregate := BatchPreCloseCheckResult{
		UnpostedJournals: PreCloseCheckUnpostedJournals{Status: CheckStatusPassed},
		TrialBalance:     PreCloseCheckTrialBalance{Status: CheckStatusPassed},
	}

	for _, p := range existingPeriods {
		unposted, err := checkUnpostedJournals(ctx, h.readModel, sobId, p.Id)
		if err != nil {
			return BatchPreCloseCheckResult{}, fmt.Errorf("failed to check unposted journals for period %d-%02d: %w", p.FiscalYear, p.PeriodNumber, err)
		}
		aggregate.UnpostedJournals.Count += unposted.Count
		aggregate.UnpostedJournals.Journals = append(aggregate.UnpostedJournals.Journals, unposted.Journals...)
		if unposted.Status == CheckStatusFailed {
			aggregate.UnpostedJournals.Status = CheckStatusFailed
		}

		// Only check trial balance if unposted journals passed for this period.
		if unposted.Status == CheckStatusPassed {
			tb, err := checkTrialBalance(ctx, h.readModel, sobId, p.Id)
			if err != nil {
				return BatchPreCloseCheckResult{}, fmt.Errorf("failed to check trial balance for period %d-%02d: %w", p.FiscalYear, p.PeriodNumber, err)
			}
			if tb.Status == CheckStatusFailed {
				aggregate.TrialBalance.Status = CheckStatusFailed
				aggregate.TrialBalance.OpeningAmount = aggregate.TrialBalance.OpeningAmount.Add(tb.OpeningAmount)
				aggregate.TrialBalance.PeriodAmount = aggregate.TrialBalance.PeriodAmount.Add(tb.PeriodAmount)
				aggregate.TrialBalance.EndingAmount = aggregate.TrialBalance.EndingAmount.Add(tb.EndingAmount)
			}
		}
	}

	// If unposted journals failed, mark trial balance as undetermined.
	if aggregate.UnpostedJournals.Status == CheckStatusFailed {
		aggregate.TrialBalance.Status = CheckStatusUndetermined
	}

	return aggregate, nil
}
