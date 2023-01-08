package app

import (
	"github/fims-proto/fims-proto-ms/internal/numbering/app/command"
	"github/fims-proto/fims-proto-ms/internal/numbering/app/query"
	"github/fims-proto/fims-proto-ms/internal/numbering/domain"
)

type Queries struct {
	IdentifierById query.IdentifierByIdHandler
}

type Commands struct {
	CreateIdentifierConfiguration command.CreateIdentifierConfigurationHandler
	GenerateNextIdentifier        command.GenerateNextIdentifierHandler
	Migrate                       command.MigrationHandler
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
	readModel query.NumberingReadModel,
) {
	a.Queries = Queries{
		IdentifierById: query.NewIdentifierByIdHandler(readModel),
	}
	a.Commands = Commands{
		CreateIdentifierConfiguration: command.NewCreateIdentifierConfigurationHandler(repo),
		GenerateNextIdentifier:        command.NewGenerateNextIdentifierHandler(repo, readModel),
		Migrate:                       command.NewMigrationHandler(repo),
	}
}
