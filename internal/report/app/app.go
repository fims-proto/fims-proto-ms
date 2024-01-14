package app

import (
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/report/domain"
)

type Queries struct{}

type Commands struct{}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(
	repo domain.Repository,
	readModel query.GeneralLedgerReadModel,
) {
	a.Queries = Queries{}
	a.Commands = Commands{}
}
