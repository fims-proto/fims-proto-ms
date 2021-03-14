package app

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
)

type Queries struct{}

type Commands struct {
	UpdateLedgerBalance command.UpdateLedgerBalanceHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(repo domain.Repository, accountService command.AccountService, voucherService command.VoucherService) {
	a.Queries = Queries{}
	a.Commands = Commands{
		UpdateLedgerBalance: command.NewUpdateLedgerBalanceHandler(repo, accountService, voucherService),
	}
}
