package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/ledger/app"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"
)

type LedgerInterface struct {
	app app.Application
}

func NewLedgerInterface(app app.Application) LedgerInterface {
	return LedgerInterface{app: app}
}

func (i LedgerInterface) PostVoucher(ctx context.Context, cmd UpdateLedgerBalanceCmd) error {
	return i.app.Commands.UpdateLedgerBalanceHandler.Handle(ctx, cmd.(command.UpdateLedgerBalanceCmd))
}
