package intraprocess

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app"
)

type GeneralLedgerInterface struct {
	app *app.Application
}

func NewGeneralLedgerInterface(app *app.Application) GeneralLedgerInterface {
	return GeneralLedgerInterface{app: app}
}

func (i GeneralLedgerInterface) Initialize(ctx context.Context, cmd command.InitializeCmd) error {
	return i.app.Commands.Initialize.Handle(ctx, cmd)
}
