package app

import (
	"github/fims-proto/fims-proto-ms/internal/sob/app/command"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"
)

type Queries struct {
	ReadSobs query.ReadSobsHandler
}

type Commands struct {
	CreateSob command.CreateSobHandler
	Migrate   command.MigrationHanlder
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
	readModel query.SobsReadModel,
) {
	a.Queries = Queries{
		ReadSobs: query.NewReadSobsHandler(readModel),
	}
	a.Commands = Commands{
		CreateSob: command.NewCreateSobHandler(repo),
		Migrate:   command.NewMigrationHanlder(repo),
	}
}
