package app

import (
	"github/fims-proto/fims-proto-ms/internal/account/app/command"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/app/service"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
)

type Queries struct {
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
	sobService service.SobService,
	numberingService service.NumberingService,
	periodByIdReadModel query.PeriodByIdReadModel,
	allAccountConfigurationsReadModel query.AllAccountConfigurationsReadModel,
	accountsInPeriodReadModel query.AccountsInPeriodReadModel,
	superiorAccountsReadModel query.SuperiorAccountConfigurationsReadModel,
) {
	a.Queries = Queries{}
	a.Commands = Commands{
		InitialAccountConfigurations: command.NewInitialAccountConfigurationHandler(repo, sobService),
		CreatePeriod: command.NewCreatePeriodHandler(
			repo,
			numberingService,
			periodByIdReadModel,
			allAccountConfigurationsReadModel,
			accountsInPeriodReadModel,
		),
		PostAccounts: command.NewPostAccountsHandler(repo, superiorAccountsReadModel),
		Migrate:      command.NewMigrationHandler(repo),
	}
}
