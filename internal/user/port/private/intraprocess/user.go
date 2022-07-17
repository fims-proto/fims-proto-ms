package intraprocess

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/user/app"
	"github/fims-proto/fims-proto-ms/internal/user/app/query"
)

type UserInterface struct {
	app *app.Application
}

func NewUserInterface(app *app.Application) UserInterface {
	return UserInterface{app: app}
}

func (i UserInterface) ReadUsersByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]query.User, error) {
	return i.app.Queries.ReadUsers.HandleReadByIds(ctx, ids)
}
