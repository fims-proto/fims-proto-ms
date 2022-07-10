package intraprocess

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"

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

func (i LedgerInterface) ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (query.Period, error) {
	return i.app.Queries.ReadLedgers.HandleReadPeriodByTime(ctx, sobId, timePoint)
}

func (i LedgerInterface) ReadPeriodsByIds(ctx context.Context, periodIds []uuid.UUID) (map[uuid.UUID]query.Period, error) {
	return i.app.Queries.ReadLedgers.HandleReadPeriodsByIds(ctx, periodIds)
}
