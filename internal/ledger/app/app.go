package app

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
)

type Queries struct {
	ReadLedgers query.ReadLedgerHandler
}

type Commands struct {
	CreatePeriod           command.CreatePeriodHandler
	ClosePeriod            command.ClosePeriodHandler
	CreateLedgersForPeriod command.CreatePeriodLedgersHandler
	PostLedgers            command.PostLedgersHandler
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
	selfService service.SelfService,
	accountService service.AccountService,
	numberingService service.NumberingService,
) {
	a.Queries = Queries{
		ReadLedgers: query.NewReadLedgerHandler(readModel, accountService),
	}
	a.Commands = Commands{
		CreatePeriod:           command.NewCreatePeriodHandler(repo, readModel, selfService, numberingService),
		ClosePeriod:            command.NewClosePeriodHandler(repo),
		CreateLedgersForPeriod: command.NewCreatePeriodLedgersHandler(repo, readModel, accountService),
		PostLedgers:            command.NewPostLedgersHandler(repo, accountService),
		Migrate:                command.NewMigrationHandler(repo),
	}
}
