package app

import "github/fims-proto/fims-proto-ms/internal/account/app/query"

type Queries struct {
	ReadAccounts query.ReadAccountsHandler
}

type Commands struct{}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(readModel query.AccountsReadModel) {
	a.Queries = Queries{
		ReadAccounts: query.NewReadAccountsHandler(readModel),
	}
	a.Commands = Commands{}
}
