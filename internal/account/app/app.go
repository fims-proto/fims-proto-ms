package app

import (
	"github/fims-proto/fims-proto-ms/internal/account/app/command"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/app/service"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
)

type Queries struct {
	PagingAccountConfigurations    query.PagingAccountConfigurationsHandler
	AccountConfigurationsByIds     query.AccountConfigurationsByIdsHandler
	AccountConfigurationsByNumbers query.AccountConfigurationsByNumbersHandler
	PagingPeriods                  query.PagingPeriodsHandler
	PeriodByTime                   query.PeriodByTimeHandler
	PeriodsByIds                   query.PeriodsByIdsHandler
	PagingAccountsByPeriod         query.PagingAccountsByPeriodHandler
}

type Commands struct {
	InitialAccountConfigurations command.InitialAccountConfigurationHandler
	CreatePeriod                 command.CreatePeriodHandler
	PostAccounts                 command.PostAccountsHandler
	Migrate                      command.MigrationHandler
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
		PagingAccountConfigurations:    query.NewPagingAccountConfigurationsHandler(readModel),
		AccountConfigurationsByNumbers: query.NewAccountConfigurationsByNumbersHandler(readModel),
		AccountConfigurationsByIds:     query.NewAccountConfigurationsByIdsHandler(readModel),
		PagingPeriods:                  query.NewPagingPeriodsHandler(readModel),
		PeriodByTime:                   query.NewPeriodByTimeHandler(readModel),
		PeriodsByIds:                   query.NewPeriodsByIdsHandler(readModel),
		PagingAccountsByPeriod:         query.NewPagingAccountsByPeriodHandler(readModel),
	}
	a.Commands = Commands{
		InitialAccountConfigurations: command.NewInitialAccountConfigurationHandler(repo, sobService),
		CreatePeriod: command.NewCreatePeriodHandler(
			repo,
			numberingService,
			readModel,
		),
		PostAccounts: command.NewPostAccountsHandler(repo, readModel),
		Migrate:      command.NewMigrationHandler(repo),
	}
}
