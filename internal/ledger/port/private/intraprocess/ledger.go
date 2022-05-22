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

func (i LedgerInterface) AppendLedgerLogs(ctx context.Context, logs []command.AppendLedgerLogCmd) error {
	return i.app.Commands.AppendLedgerLogs.Handle(ctx, logs)
}
