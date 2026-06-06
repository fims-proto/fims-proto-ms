package command

import (
	"context"
	"errors"
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// initializeAllLedgers creates ledgers for current period.
// If it's first period in the SoB, the opening balances are zero.
// Initializing the opening balance for the first period is handled by InitializeLedgersBalanceHandler
func initializeAllLedgers(ctx context.Context, repo domain.Repository, sobId uuid.UUID) error {
	// read current period
	currentPeriod, err := repo.ReadCurrentPeriod(ctx, sobId)
	if err != nil {
		return fmt.Errorf("failed to read current period: %w", err)
	}

	// read previous period
	previousPeriod, err := repo.ReadPreviousPeriod(ctx, currentPeriod.Id())
	if err != nil && !errors.Is(err, commonErrors.ErrRecordNotFound()) {
		return fmt.Errorf("initialize ledger failed: %w", err)
	}

	return initializeLedgers(ctx, repo, sobId, currentPeriod, previousPeriod)
}

func initializeLedgers(ctx context.Context, repo domain.Repository, sobId uuid.UUID, currentPeriod *period.Period, previousPeriod *period.Period) error {
	// read ledgers in previous period
	ledgersInPreviousPeriod := make(map[uuid.UUID]ledger.Ledger) // key: account id, value: ledger

	if previousPeriod != nil {
		// normal ledgers
		ledgers, err := repo.ReadLedgersByPeriod(ctx, previousPeriod.Id())
		if err != nil {
			return fmt.Errorf("failed to read ledgers in previous period: %w", err)
		}
		for _, previousLedger := range ledgers {
			ledgersInPreviousPeriod[previousLedger.AccountId()] = *previousLedger
		}
	}

	// create ledgers based on accounts
	accounts, err := repo.ReadAllAccounts(ctx, sobId)
	if err != nil {
		return fmt.Errorf("failed to read accounts: %w", err)
	}

	var ledgers []*ledger.Ledger
	for _, account := range accounts {
		// move previous ending balance to current balance
		openingAmount := decimal.Zero
		endingAmount := decimal.Zero

		previousLedger, ok := ledgersInPreviousPeriod[account.Id()]
		if ok {
			openingAmount = previousLedger.EndingAmount()
			endingAmount = previousLedger.EndingAmount()
		}

		ledgerBO, err := ledger.New(
			uuid.New(),
			account.SobId(),
			currentPeriod.Id(),
			account.Id(),
			account,
			openingAmount, // openingAmount
			decimal.Zero,  // periodAmount
			decimal.Zero,  // periodDebit
			decimal.Zero,  // periodCredit
			endingAmount,  // endingAmount
		)
		if err != nil {
			return fmt.Errorf("should not happen, failed to create ledger: %w", err)
		}

		ledgers = append(ledgers, ledgerBO)
	}

	return repo.CreateLedgers(ctx, ledgers)
}
