package intraprocess

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/ledger/app"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"
)

type LedgerInterface struct {
	app *app.Application
}

func NewLedgerInterface(app *app.Application) LedgerInterface {
	return LedgerInterface{app: app}
}

func (i LedgerInterface) PostVoucher(ctx context.Context, req UpdateLedgerBalanceRequest) error {
	return i.app.Commands.UpdateLedgerBalance.Handle(ctx, req.mapToCommand())
}

func (i LedgerInterface) LoadLedgers(ctx context.Context, reqs []LoadLedgersRequest) error {
	var cmds []command.LedgerDataloadCmd
	for _, req := range reqs {
		cmds = append(cmds, req.mapToCommand())
	}
	return i.app.Commands.LoadLedgers.Handle(ctx, cmds)
}
