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

func (i LedgerInterface) InitializeFirstPeriod(ctx context.Context, sobId uuid.UUID, financialYear, number int) error {
	startDateOfMonth := time.Date(financialYear, time.Month(number), 1, 0, 0, 0, 0, time.UTC)
	cmd := command.CreatePeriodCmd{
		PreviousPeriodId: uuid.Nil,
		SobId:            sobId,
		FinancialYear:    financialYear,
		Number:           number,
		OpeningTime:      startDateOfMonth,
	}

	_, err := i.app.Commands.CreatePeriod.Handle(ctx, cmd)
	return err
}

func (i LedgerInterface) InitializeLedgersForPeriod(ctx context.Context, periodId uuid.UUID) error {
	return i.app.Commands.CreateLedgersForPeriod.Handle(ctx, command.CreatePeriodLedgersCmd{PeriodId: periodId})
}

func (i LedgerInterface) PostLedgers(ctx context.Context, cmd command.PostLedgersCmd) error {
	return i.app.Commands.PostLedgers.Handle(ctx, cmd)
}

func (i LedgerInterface) ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (query.Period, error) {
	return i.app.Queries.ReadLedgers.HandleReadPeriodByTime(ctx, sobId, timePoint)
}

func (i LedgerInterface) ReadPeriodsByIds(ctx context.Context, periodIds []uuid.UUID) (map[uuid.UUID]query.Period, error) {
	return i.app.Queries.ReadLedgers.HandleReadPeriodsByIds(ctx, periodIds)
}
