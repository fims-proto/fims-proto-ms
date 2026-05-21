package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"

	"github.com/google/uuid"
)

// prepareJournalLines prepares journal line domain objects, validates dimension options,
// and performs necessary checks.
func prepareJournalLines(
	ctx context.Context,
	repo domain.Repository,
	dimensionService service.DimensionService,
	sobId uuid.UUID,
	commands []JournalLineCmd,
) ([]*journal.JournalLine, error) {
	var rawAccountNumbers []string
	for _, item := range commands {
		rawAccountNumbers = append(rawAccountNumbers, item.RawAccountNumber)
	}

	// validate account numbers
	accounts, err := repo.ReadAccountsByRawNumbers(ctx, sobId, rawAccountNumbers)
	if err != nil {
		return nil, fmt.Errorf("failed to read accounts: %w", err)
	}

	accountsMap := utils.SliceToMap(
		accounts,
		func(a *account.Account) string { return a.RawAccountNumber() },
		func(a *account.Account) *account.Account { return a },
	)

	// determine cash flow classification requirements
	hasCashEquivalentLine := false
	hasNonCashEquivalentLine := false
	for i := range commands {
		a := accountsMap[rawAccountNumbers[i]]
		if a.IsCashEquivalent() {
			hasCashEquivalentLine = true
		} else {
			hasNonCashEquivalentLine = true
		}
	}

	// validate cash flow item IDs when mixed entry (cash + non-cash lines)
	if hasCashEquivalentLine && hasNonCashEquivalentLine {
		var cfItemIds []uuid.UUID
		for i, item := range commands {
			a := accountsMap[rawAccountNumbers[i]]
			if !a.IsCashEquivalent() {
				if item.CashFlowItemId == nil {
					return nil, commonErrors.NewInvalidInputError(commonErrors.SlugJournalLineMissingCashFlowItem)
				}
				cfItemIds = append(cfItemIds, *item.CashFlowItemId)
			}
		}

		// batch-validate all provided CF item IDs exist in this SoB
		existing, err := repo.ReadExistingCashFlowItemIds(ctx, sobId, cfItemIds)
		if err != nil {
			return nil, fmt.Errorf("failed to validate cash flow item ids: %w", err)
		}
		existingSet := utils.SliceToMap(existing,
			func(id uuid.UUID) uuid.UUID { return id },
			func(id uuid.UUID) struct{} { return struct{}{} },
		)
		for _, id := range cfItemIds {
			if _, ok := existingSet[id]; !ok {
				return nil, commonErrors.NewInvalidInputError(commonErrors.SlugJournalLineCashFlowItemNotFound, id)
			}
		}
	}

	// prepare journal lines
	var journalLines []*journal.JournalLine
	for i, item := range commands {
		itemId := item.Id
		if itemId == uuid.Nil {
			itemId = uuid.New()
		}

		a := accountsMap[rawAccountNumbers[i]]

		// Validate dimension options for this journal line against the account's required categories.
		if err = dimensionService.ValidateOptions(ctx, a.DimensionCategoryIds(), item.DimensionOptionIds); err != nil {
			return nil, err
		}

		// Only attach CF item ID when this is a mixed entry (cash + non-cash)
		var cashFlowItemId *uuid.UUID
		if hasCashEquivalentLine && hasNonCashEquivalentLine && !a.IsCashEquivalent() {
			cashFlowItemId = item.CashFlowItemId
		}

		journalLine, err := journal.NewJournalLine(
			itemId,
			a,
			item.Text,
			item.Amount,
			item.DimensionOptionIds,
			cashFlowItemId,
		)
		if err != nil {
			return nil, err
		}

		journalLines = append(journalLines, journalLine)
	}

	return journalLines, nil
}

// readPeriodIdAndCheck tries to get period id by given transaction date of a journal, and will also check if the period is closed.
// if no period exists for given transaction date, it creates one
func readPeriodIdAndCheck(
	ctx context.Context,
	repo domain.Repository,
	numberingService service.NumberingService,
	sobId uuid.UUID,
	transactionDate transaction_date.TransactionDate,
) (*period.Period, error) {
	fiscalYear := transactionDate.Year
	periodNumber := transactionDate.Month

	p, err := createPeriodIfNotExists(ctx, createPeriodCmd{
		SobId:      sobId,
		PeriodId:   uuid.Nil,
		FiscalYear: fiscalYear,
		Number:     periodNumber,
	}, repo, numberingService)
	if err != nil {
		return nil, err
	}

	if p.IsClosed() {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugPeriodClosed)
	}

	return p, nil
}
