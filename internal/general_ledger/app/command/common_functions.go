package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"

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
	var accountNumbers []string
	for _, item := range commands {
		accountNumbers = append(accountNumbers, item.AccountNumber)
	}

	// validate account numbers
	accounts, err := repo.ReadAccountsByNumbers(ctx, sobId, accountNumbers)
	if err != nil {
		return nil, fmt.Errorf("failed to read accounts: %w", err)
	}

	accountsMap := utils.SliceToMap(
		accounts,
		func(a *account.Account) string { return a.AccountNumber() },
		func(a *account.Account) *account.Account { return a },
	)

	// prepare journal lines
	var journalLines []*journal.JournalLine
	for _, item := range commands {
		itemId := item.Id
		if itemId == uuid.Nil {
			itemId = uuid.New()
		}

		a := accountsMap[item.AccountNumber]

		// Validate dimension options for this journal line against the account's required categories.
		if err = dimensionService.ValidateOptions(ctx, a.DimensionCategoryIds(), item.DimensionOptionIds); err != nil {
			return nil, err
		}

		journalLine, err := journal.NewJournalLine(
			itemId,
			a,
			item.Text,
			item.Amount,
			item.DimensionOptionIds,
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
		return nil, commonErrors.ErrPeriodClosed()
	}

	return p, nil
}
