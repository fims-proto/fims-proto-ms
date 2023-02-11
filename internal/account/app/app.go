package app

import (
	"github/fims-proto/fims-proto-ms/internal/account/app/command"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/app/service"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
)

type Queries struct {
	PagingAccounts        query.PagingAccountsHandler
	AccountsByIds         query.AccountsByIdsHandler
	AccountsByNumbers     query.AccountsByNumbersHandler
	OpenPeriod            query.OpenPeriodHandler
	PagingPeriods         query.PagingPeriodsHandler
	PeriodByTime          query.PeriodByTimeHandler
	PeriodsByIds          query.PeriodsByIdsHandler
	PagingLedgersByPeriod query.PagingLedgersByPeriodHandler
}

type Commands struct {
	InitialAccounts command.InitialAccountsHandler
	CreatePeriod    command.CreatePeriodHandler
	CreateLedgers   command.CreateLedgersHandler
	PostAccounts    command.PostAccountsHandler
	Migrate         command.MigrationHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(
	repo domain.Repository,
	readModel query.AccountReadModel,
	sobService service.SobService,
	numberingService service.NumberingService,
) {
	a.Queries = Queries{
		PagingAccounts:        query.NewPagingAccountsHandler(readModel),
		AccountsByNumbers:     query.NewAccountsByNumbersHandler(readModel),
		AccountsByIds:         query.NewAccountsByIdsHandler(readModel),
		OpenPeriod:            query.NewOpenPeriodHandler(readModel),
		PagingPeriods:         query.NewPagingPeriodsHandler(readModel),
		PeriodByTime:          query.NewPeriodByTimeHandler(readModel),
		PeriodsByIds:          query.NewPeriodsByIdsHandler(readModel),
		PagingLedgersByPeriod: query.NewPagingLedgersByPeriodHandler(readModel),
	}
	a.Commands = Commands{
		InitialAccounts: command.NewInitialAccountHandler(repo, sobService),
		CreatePeriod: command.NewCreatePeriodHandler(
			repo,
			numberingService,
			readModel,
		),
		CreateLedgers: command.NewCreateLedgersHandler(repo, readModel),
		PostAccounts:  command.NewPostAccountsHandler(repo, readModel),
		Migrate:       command.NewMigrationHandler(repo),
	}
}
