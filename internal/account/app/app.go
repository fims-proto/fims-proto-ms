package app

import (
	"github/fims-proto/fims-proto-ms/internal/account/app/command"
	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/account/domain"
)

type Queries struct {
	ReadAccounts query.ReadAccountsHandler
}

type Commands struct {
	LoadAccounts command.AccountDataLoadHandler
	Migrate      command.MigrationHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(
	readModel query.AccountsReadModel,
	repo domain.Repository,
	sobService command.SobService,
) {
	a.Queries = Queries{
		ReadAccounts: query.NewReadAccountsHandler(readModel),
	}
	a.Commands = Commands{
		LoadAccounts: command.NewAccountDataLoadHandler(repo, sobService),
		Migrate:      command.NewMigrationHandler(repo),
	}
}
