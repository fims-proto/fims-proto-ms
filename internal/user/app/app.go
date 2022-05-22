package app

import (
	"github/fims-proto/fims-proto-ms/internal/user/app/command"
	"github/fims-proto/fims-proto-ms/internal/user/app/query"
	"github/fims-proto/fims-proto-ms/internal/user/domain"
)

type Queries struct {
	ReadUsers query.ReadUsersHandler
}

type Commands struct {
	UpdateUser command.UpdateUserHandler
	Migrate    command.MigrationHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(
	readModel query.UsersReadModel,
	repo domain.Repository,
) {
	a.Queries = Queries{
		ReadUsers: query.NewReadUsersHandler(readModel),
	}

	a.Commands = Commands{
		UpdateUser: command.NewUpdateUserHandler(repo),
		Migrate:    command.NewMigrationHandler(repo),
	}
}
