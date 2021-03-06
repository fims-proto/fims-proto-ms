package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/app"
)

type AccountInterface struct {
	app app.Application
}

func NewAccountInterface(app app.Application) AccountInterface {
	return AccountInterface{app: app}
}

func (i AccountInterface) ValidateExistence(ctx context.Context, accNumbers []string) error {
	return i.app.Queries.ValidateAccounts.HandleValidateExistence(ctx, accNumbers)
}
