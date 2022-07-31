package command

import (
	"context"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account_configuration"
	"github/fims-proto/fims-proto-ms/internal/account/domain/period"
	"time"

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
	repo                              domain.Repository
	numberingService                  service.NumberingService
	periodByIdReadModel               query.PeriodByIdReadModel
	allAccountConfigurationsReadModel query.AllAccountConfigurationsReadModel
	accountsInPeriodReadModel         query.AccountsInPeriodReadModel
}

func NewCreatePeriodHandler(
	repo domain.Repository,
	numberingService service.NumberingService,
	periodByIdReadModel query.PeriodByIdReadModel,
	allAccountConfigurationsReadModel query.AllAccountConfigurationsReadModel,
	accountsInPeriodReadModel query.AccountsInPeriodReadModel,
) CreatePeriodHandler {
	if repo == nil {
		panic("nil ledger repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	if periodByIdReadModel == nil {
		panic("nil period read model")
	}

	if allAccountConfigurationsReadModel == nil {
		panic("nil period read model")
	}

	if accountsInPeriodReadModel == nil {
		panic("nil period read model")
	}

	return CreatePeriodHandler{
		repo:                              repo,
		numberingService:                  numberingService,
		periodByIdReadModel:               periodByIdReadModel,
		allAccountConfigurationsReadModel: allAccountConfigurationsReadModel,
		accountsInPeriodReadModel:         accountsInPeriodReadModel,
	}
}

func (h CreatePeriodHandler) Handle(ctx context.Context, cmd CreatePeriodCmd) error {
	// use previous period ending time as new opening time if previous period provided
	// otherwise using given opening time
	openingTime := cmd.OpeningTime
	if cmd.PreviousPeriodId != uuid.Nil {
		previousPeriod, err := h.periodByIdReadModel.PeriodById(ctx, cmd.PreviousPeriodId)
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
		if err = h.createAccountsForPeriod(ctx, *p); err != nil {
			return errors.Wrap(err, "failed to create accounts for period")
		}

		// create numbering configuration for vouchers in this period
		if err = h.numberingService.InitializeIdentifierConfigurationForVoucher(ctx, cmd.PeriodId); err != nil {
			return errors.Wrap(err, "failed to create numbering configuration for period")
		}

		return nil
	})
}

func (h CreatePeriodHandler) createAccountsForPeriod(ctx context.Context, period period.Period) error {
	// read all account configurations
	configurations, err := h.allAccountConfigurationsReadModel.AllAccountConfigurations(ctx, period.SobId())
	if err != nil {
		return errors.Wrap(err, "failed to create accounts for period")
	}

	// read all accounts in previous period if applicable
	accountsInPreviousPeriod := make(map[uuid.UUID]query.Account) // key: AccountId, value: account
	if period.PreviousPeriodId() != uuid.Nil {
		accounts, err := h.accountsInPeriodReadModel.AccountsInPeriod(ctx, period.SobId(), period.PreviousPeriodId())
		if err != nil {
			return errors.Wrap(err, "failed to read accounts in previous period")
		}

		for _, previousAccount := range accounts {
			accountsInPreviousPeriod[previousAccount.AccountId] = previousAccount
		}
	}

	// create accounts based on configuration
	var accounts []*account.Account
	for _, configuration := range configurations {
		// move previous ending balance to opening balance
		openingBalance := decimal.Zero
		previousAccount, ok := accountsInPreviousPeriod[configuration.AccountId]
		if ok {
			openingBalance = previousAccount.EndingBalance
		}

		accountConfiguration, err := account_configuration.New(
			configuration.SobId,
			configuration.AccountId,
			configuration.SuperiorAccountId,
			configuration.Title,
			configuration.AccountNumber,
			configuration.NumberHierarchy,
			configuration.Level,
			configuration.AccountType,
			configuration.BalanceDirection,
		)
		if err != nil {
			// should not happen
			return errors.Wrap(err, "should not happen, failed to create account configuration")
		}

		domainAccount, err := account.New(
			configuration.SobId,
			configuration.AccountId,
			period.PeriodId(),
			openingBalance,
			decimal.Zero,
			decimal.Zero,
			decimal.Zero,
			*accountConfiguration,
		)
		if err != nil {
			return errors.Wrap(err, "should not happen, failed to create account")
		}

		accounts = append(accounts, domainAccount)
	}

	return h.repo.CreateAccounts(ctx, accounts)
}
