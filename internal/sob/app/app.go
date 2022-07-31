package app

import (
	"github/fims-proto/fims-proto-ms/internal/sob/app/command"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"github/fims-proto/fims-proto-ms/internal/sob/app/service"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"
)

type Queries struct {
	ReadSobs query.ReadSobsHandler
}

type Commands struct {
	CreateSob command.CreateSobHandler
	UpdateSob command.UpdateSobHandler
	Migrate   command.MigrationHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(
	readModel query.SobsReadModel,
	repo domain.Repository,
	accountService service.AccountService,
) {
	a.Queries = Queries{
		ReadSobs: query.NewReadSobsHandler(readModel),
	}
	a.Commands = Commands{
		CreateSob: command.NewCreateSobHandler(repo, accountService),
		UpdateSob: command.NewUpdateSobHandler(repo),
		Migrate:   command.NewMigrationHandler(repo),
	}
}
