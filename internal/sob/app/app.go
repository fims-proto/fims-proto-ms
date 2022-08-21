package app

import (
	"github/fims-proto/fims-proto-ms/internal/sob/app/command"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"github/fims-proto/fims-proto-ms/internal/sob/app/service"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"
)

type Queries struct {
	PagingSobs query.PagingSobsHandler
	SobById    query.SobByIdHandler
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
	repo domain.Repository,
	readModel query.SobReadModel,
	accountService service.AccountService,
) {
	a.Queries = Queries{
		PagingSobs: query.NewPagingSobsHandler(readModel),
		SobById:    query.NewSobByIdHandler(readModel),
	}
	a.Commands = Commands{
		CreateSob: command.NewCreateSobHandler(repo, accountService),
		UpdateSob: command.NewUpdateSobHandler(repo),
		Migrate:   command.NewMigrationHandler(repo),
	}
}
