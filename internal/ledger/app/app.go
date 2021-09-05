package app

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
)

type Queries struct {
	ReadLedgers query.ReadLedgersHandler
}

type Commands struct {
	UpdateLedgerBalance command.UpdateLedgerBalanceHandler
	LoadLedgers         command.LedgerDataloadHandler
	Migrate             command.MigrationHanlder
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(repo domain.Repository, readModel query.LedgersReadModel, accountService command.AccountService, voucherService command.VoucherService) {
	a.Queries = Queries{
		ReadLedgers: query.NewReadLedgersHandler(readModel),
	}
	a.Commands = Commands{
		UpdateLedgerBalance: command.NewUpdateLedgerBalanceHandler(repo, accountService, voucherService),
		LoadLedgers:         command.NewLedgerDataloadHandler(repo),
		Migrate:             command.NewMigrationHanlder(repo),
	}
}
