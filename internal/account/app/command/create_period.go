package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/account/domain/ledger"

	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"

	"github/fims-proto/fims-proto-ms/internal/account/app/service"

	"github/fims-proto/fims-proto-ms/internal/account/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// create first period in the SoB
// create period, checking number, using ending time as opening time

type CreatePeriodCmd struct {
	SobId            uuid.UUID
	PeriodId         uuid.UUID
	PreviousPeriodId uuid.UUID
	FinancialYear    int
	Number           int
	OpeningTime      time.Time
	EndingTime       time.Time
}

type CreatePeriodHandler struct {
	repo             domain.Repository
	numberingService service.NumberingService
	readModel        query.AccountReadModel
}

func NewCreatePeriodHandler(
	repo domain.Repository,
	numberingService service.NumberingService,
	readModel query.AccountReadModel,
) CreatePeriodHandler {
	if repo == nil {
		panic("nil ledger repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	if readModel == nil {
		panic("nil read model")
	}

	return CreatePeriodHandler{
		repo:             repo,
		numberingService: numberingService,
		readModel:        readModel,
	}
}

func (h CreatePeriodHandler) Handle(ctx context.Context, cmd CreatePeriodCmd) error {
	// use previous period ending time as new opening time if previous period provided
	// otherwise using given opening time
	openingTime := cmd.OpeningTime
	if cmd.PreviousPeriodId != uuid.Nil {
		previousPeriod, err := h.readModel.PeriodById(ctx, cmd.PreviousPeriodId)
		if err != nil {
			return errors.Wrap(err, "failed to read previous period")
		}
		if previousPeriod.SobId != cmd.SobId {
			return errors.Wrap(err, "sob id not equals to the one from previous period")
		}
		if !previousPeriod.IsClosed {
			return errors.Wrap(err, "previous period not closed")
		}
		openingTime = previousPeriod.EndingTime
	}

	p, err := period.New(cmd.PeriodId, cmd.SobId, cmd.PreviousPeriodId, cmd.FinancialYear, cmd.Number, openingTime, cmd.EndingTime, false)
	if err != nil {
		return errors.Wrap(err, "failed to create period domain model")
	}

	return h.repo.CreatePeriod(ctx, p, func() error {
		// create accounts for this period
		if err = h.createLedgersForPeriod(ctx, *p); err != nil {
			return errors.Wrap(err, "failed to create accounts for period")
		}

		// create numbering configuration for journal entries in this period
		if err = h.numberingService.InitializeIdentifierConfigurationForJournal(ctx, cmd.PeriodId); err != nil {
			return errors.Wrap(err, "failed to create numbering configuration for period")
		}

		return nil
	})
}

func (h CreatePeriodHandler) createLedgersForPeriod(ctx context.Context, period period.Period) error {
	// read all account configurations
	configurations, err := h.readModel.AllAccounts(ctx, period.SobId())
	if err != nil {
		return errors.Wrap(err, "failed to create accounts for period")
	}

	// read all accounts in previous period if applicable
	accountsInPreviousPeriod := make(map[uuid.UUID]query.Ledger) // key: Id, value: account
	if period.PreviousPeriodId() != uuid.Nil {
		accounts, err := h.readModel.LedgersInPeriod(ctx, period.SobId(), period.PreviousPeriodId())
		if err != nil {
			return errors.Wrap(err, "failed to read accounts in previous period")
		}

		for _, previousAccount := range accounts {
			accountsInPreviousPeriod[previousAccount.AccountId] = previousAccount
		}
	}

	// create accounts based on configuration
	var ledgers []*ledger.Ledger
	for _, configuration := range configurations {
		// move previous ending balance to opening balance
		openingBalance := decimal.Zero
		previousAccount, ok := accountsInPreviousPeriod[configuration.Id]
		if ok {
			openingBalance = previousAccount.EndingBalance
		}

		accountConfiguration, err := account.New(configuration.Id, configuration.SobId, configuration.SuperiorAccountId, configuration.Title, configuration.AccountNumber, configuration.NumberHierarchy, configuration.Level, configuration.AccountType, configuration.BalanceDirection)
		if err != nil {
			// should not happen
			return errors.Wrap(err, "should not happen, failed to create account configuration")
		}

		domainLedger, err := ledger.New(
			uuid.New(),
			configuration.SobId,
			configuration.Id,
			period.Id(),
			openingBalance,
			decimal.Zero,
			decimal.Zero,
			decimal.Zero,
			*accountConfiguration,
		)
		if err != nil {
			return errors.Wrap(err, "should not happen, failed to create account")
		}

		ledgers = append(ledgers, domainLedger)
	}

	return h.repo.CreateLedgers(ctx, ledgers)
}
