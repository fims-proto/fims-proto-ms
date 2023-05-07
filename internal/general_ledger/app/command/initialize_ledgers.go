package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"
)

type InitializeLedgersCmd struct {
	SobId    uuid.UUID
	PeriodId uuid.UUID
}

func initializeLedgers(ctx context.Context, cmd InitializeLedgersCmd, repo domain.Repository, readModel query.GeneralLedgerReadModel) error {
	// read period
	period, err := readModel.PeriodById(ctx, cmd.PeriodId)
	if err != nil {
		return errors.Wrap(err, "failed to create ledgers for period")
	}

	// read all accounts
	accounts, err := readModel.AllAccounts(ctx, cmd.SobId)
	if err != nil {
		return errors.Wrap(err, "failed to create ledgers for period")
	}

	// read all ledgers in previous period if applicable
	ledgersInPreviousPeriod := make(map[uuid.UUID]query.Ledger) // key: Id, value: account

	previousPeriodTime := period.OpeningTime.AddDate(0, -1, 0) // one month before

	previousPeriod, _ := readModel.PeriodByFiscalYearAndNumber(ctx, cmd.SobId, previousPeriodTime.Year(), int(previousPeriodTime.Month()))
	if previousPeriod.Id != uuid.Nil {
		ledgers, err := readModel.LedgersInPeriod(ctx, cmd.SobId, previousPeriod.Id)
		if err != nil {
			return errors.Wrap(err, "failed to read ledgers in previous period")
		}

		for _, previousLedger := range ledgers {
			ledgersInPreviousPeriod[previousLedger.AccountId] = previousLedger
		}
	}

	// create ledgers based on accounts
	var ledgers []*ledger.Ledger
	for _, accountDTO := range accounts {
		// move previous ending balance to current balance
		openingBalance := decimal.Zero
		endingBalance := decimal.Zero

		previousLedger, ok := ledgersInPreviousPeriod[accountDTO.Id]
		if ok {
			openingBalance = previousLedger.EndingBalance
			endingBalance = previousLedger.EndingBalance
		}

		accountBO, err := account.New(
			accountDTO.Id,
			accountDTO.SobId,
			accountDTO.SuperiorAccountId,
			accountDTO.Title,
			accountDTO.AccountNumber,
			accountDTO.NumberHierarchy,
			accountDTO.Level,
			accountDTO.AccountType,
			accountDTO.BalanceDirection,
		)
		if err != nil {
			// should not happen
			return errors.Wrap(err, "should not happen, failed to create account")
		}

		domainLedger, err := ledger.New(
			uuid.New(),
			accountDTO.SobId,
			accountDTO.Id,
			cmd.PeriodId,
			openingBalance,
			endingBalance,
			decimal.Zero,
			decimal.Zero,
			*accountBO,
		)
		if err != nil {
			return errors.Wrap(err, "should not happen, failed to create account")
		}

		ledgers = append(ledgers, domainLedger)
	}

	return repo.CreateLedgers(ctx, ledgers)
}
