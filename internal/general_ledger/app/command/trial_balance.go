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
	var totalOpeningAmount,
		totalPeriodAmount,
		totalEndingAmount decimal.Decimal

	// sum
	for _, l := range ledgers {
		totalOpeningAmount = totalOpeningAmount.Add(l.OpeningAmount())
		totalPeriodAmount = totalPeriodAmount.Add(l.PeriodAmount())
		totalEndingAmount = totalEndingAmount.Add(l.EndingAmount())
	}

	// Trial balance: sum of all signed amounts should be zero (debits = credits)
	if !totalOpeningAmount.IsZero() {
		return commonErrors.NewSlugError("period-close-openingBalanceUnequal")
	}
	if !totalPeriodAmount.IsZero() {
		return commonErrors.NewSlugError("period-close-periodBalanceUnequal")
	}
	if !totalEndingAmount.IsZero() {
		return commonErrors.NewSlugError("period-close-endingBalanceUnequal")
	}

	return nil
}
