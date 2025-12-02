package command

import (
	"context"
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func trialBalance(ctx context.Context, repo domain.Repository, sobId, periodId uuid.UUID) error {
	ledgers, err := repo.ReadFirstLevelLedgersInPeriod(ctx, sobId, periodId)
	if err != nil {
		return fmt.Errorf("failed to read 1st level ledgers: %w", err)
	}
	var totalOpeningDebit, totalOpeningCredit,
		totalPeriodDebit, totalPeriodCredit,
		totalEndingDebit, totalEndingCredit decimal.Decimal

	// sum
	for _, l := range ledgers {
		totalOpeningDebit = totalOpeningDebit.Add(l.OpeningDebitBalance())
		totalEndingDebit = totalEndingDebit.Add(l.EndingDebitBalance())

		totalPeriodDebit = totalPeriodDebit.Add(l.PeriodDebit())
		totalPeriodCredit = totalPeriodCredit.Add(l.PeriodCredit())

		totalOpeningCredit = totalOpeningCredit.Add(l.OpeningCreditBalance())
		totalEndingCredit = totalEndingCredit.Add(l.EndingCreditBalance())
	}

	if !totalOpeningDebit.Equal(totalOpeningCredit) {
		return commonErrors.NewSlugError("period-close-openingBalanceUnequal")
	}
	if !totalPeriodDebit.Equal(totalPeriodCredit) {
		return commonErrors.NewSlugError("period-close-periodBalanceUnequal")
	}
	if !totalEndingDebit.Equal(totalEndingCredit) {
		return commonErrors.NewSlugError("period-close-endingBalanceUnequal")
	}

	return nil
}
