package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/account/app"
)

type AccountInterface struct {
	app app.Application
}

func NewHandler(app app.Application) AccountInterface {
	return AccountInterface{app: app}
}

func (h AccountInterface) ValidateExistence(ctx context.Context, accNumbers []string) error {
	return h.app.Queries.ValidateAccounts.HandleValidateExistence(ctx, accNumbers)
}
