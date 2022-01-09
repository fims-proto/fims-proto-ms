package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/sob/app"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github.com/google/uuid"
)

type SobInterface struct {
	app *app.Application
}

func NewSobInterface(app *app.Application) SobInterface {
	if app == nil {
		panic("nil sob app")
	}
	return SobInterface{app: app}
}

func (i SobInterface) ReadById(ctx context.Context, sobId uuid.UUID) (query.Sob, error) {
	return i.app.Queries.ReadSobs.HandleReadById(ctx, sobId)
}
