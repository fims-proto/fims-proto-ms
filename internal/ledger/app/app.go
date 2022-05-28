package app

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
)

type Queries struct {
	ReadLedgers query.ReadLedgerHandler
}

type Commands struct {
	AppendLedgerLogs       command.AppendLedgerLogsHandler
	CreatePeriodLedgers    command.CreatePeriodLedgersHandler
	CalculateLedgerBalance command.CalculateLedgerBalanceHandler
	CreateAccountingPeriod command.CreateAccountingPeriodHandler
	CloseAccountingPeriod  command.CloseAccountingPeriodHandler
	Migrate                command.MigrationHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(
	readModel query.LedgerReadModel,
	repo domain.Repository,
	accountQueryService query.AccountService,
	accountService command.AccountService,
) {
	a.Queries = Queries{
		ReadLedgers: query.NewReadLedgerHandler(readModel, accountQueryService),
	}
	a.Commands = Commands{
		AppendLedgerLogs:       command.NewAppendLedgerLogsHandler(repo, accountService),
		CreatePeriodLedgers:    command.NewCreatePeriodLedgersHandler(repo, readModel, accountService),
		CalculateLedgerBalance: command.NewCalculateLedgerBalanceHandler(repo, readModel, accountService),
		CreateAccountingPeriod: command.NewCreateAccountingPeriodHandler(repo, readModel),
		CloseAccountingPeriod:  command.NewCloseAccountingPeriodHandler(repo),
		Migrate:                command.NewMigrationHandler(repo),
	}
}
