package intraprocess

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app"
)

type GeneralLedgerInterface struct {
	app *app.Application
}

func NewGeneralLedgerInterface(app *app.Application) GeneralLedgerInterface {
	return GeneralLedgerInterface{app: app}
}

func (i GeneralLedgerInterface) InitializeAccounts(ctx context.Context, sobId uuid.UUID) error {
	return i.app.Commands.InitialAccounts.Handle(ctx, sobId)
}

func (i GeneralLedgerInterface) CreatePeriodByNumber(ctx context.Context, cmd command.CreatePeriodCmd) error {
	return i.app.Commands.CreatePeriod.Handle(ctx, cmd)
}

func (i GeneralLedgerInterface) CreateLedgers(ctx context.Context, cmd command.CreateLedgersCmd) error {
	return i.app.Commands.CreateLedgers.Handle(ctx, cmd)
}
