package command

import (
	"context"
	"errors"
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_ledger"
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

	if err = initializeLedgers(ctx, repo, sobId, currentPeriod, previousPeriod); err != nil {
		return nil
	}

	return initializeAuxiliaryLedgers(ctx, repo, sobId, currentPeriod, previousPeriod)
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
		openingDebitBalance := decimal.Zero
		openingCreditBalance := decimal.Zero
		endingDebitBalance := decimal.Zero
		endingCreditBalance := decimal.Zero

		previousLedger, ok := ledgersInPreviousPeriod[account.Id()]
		if ok {
			openingDebitBalance = previousLedger.EndingDebitBalance()
			openingCreditBalance = previousLedger.EndingCreditBalance()
			endingDebitBalance = previousLedger.EndingDebitBalance()
			endingCreditBalance = previousLedger.EndingCreditBalance()
		}

		ledgerBO, err := ledger.New(
			uuid.New(),
			account.SobId(),
			currentPeriod.Id(),
			account.Id(),
			account,
			openingDebitBalance,
			openingCreditBalance,
			decimal.Zero,
			decimal.Zero,
			endingDebitBalance,
			endingCreditBalance,
		)
		if err != nil {
			return fmt.Errorf("should not happen, failed to create ledger: %w", err)
		}

		ledgers = append(ledgers, ledgerBO)
	}

	return repo.CreateLedgers(ctx, ledgers)
}

func initializeAuxiliaryLedgers(ctx context.Context, repo domain.Repository, sobId uuid.UUID, currentPeriod *period.Period, previousPeriod *period.Period) error {
	// Auxiliary ledgers are account-scoped (sobId + periodId + accountId + categoryId + auxiliaryAccountId)
	// We only copy auxiliary ledgers that existed in the previous period (account+category+auxiliary combinations that were used)

	if previousPeriod == nil {
		// No previous period, no auxiliary ledgers to copy
		return nil
	}

	// Read auxiliary ledgers from previous period
	previousAuxiliaryLedgers, err := repo.ReadAuxiliaryLedgersByPeriod(ctx, previousPeriod.Id())
	if err != nil {
		return fmt.Errorf("failed to read auxiliary ledgers in previous period: %w", err)
	}

	// Create new auxiliary ledgers for current period, copying balances from previous period
	var auxiliaryLedgers []*auxiliary_ledger.AuxiliaryLedger
	for _, previousLedger := range previousAuxiliaryLedgers {
		// Copy ending balance to opening balance
		auxiliaryLedger, err := auxiliary_ledger.New(
			uuid.New(),
			sobId,
			currentPeriod.Id(),
			previousLedger.AccountId(),
			previousLedger.AuxiliaryCategoryId(),
			previousLedger.AuxiliaryAccount(),
			previousLedger.EndingDebitBalance(),  // opening = previous ending
			previousLedger.EndingCreditBalance(), // opening = previous ending
			decimal.Zero,                         // period debit = 0
			decimal.Zero,                         // period credit = 0
			previousLedger.EndingDebitBalance(),  // ending = opening (no transactions yet)
			previousLedger.EndingCreditBalance(), // ending = opening (no transactions yet)
		)
		if err != nil {
			return fmt.Errorf("failed to create auxiliary ledger: %w", err)
		}

		auxiliaryLedgers = append(auxiliaryLedgers, auxiliaryLedger)
	}

	if len(auxiliaryLedgers) > 0 {
		return repo.CreateAuxiliaryLedgers(ctx, auxiliaryLedgers)
	}

	return nil
}
