package intraprocess

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/app"
	"github/fims-proto/fims-proto-ms/internal/report/app/command"

	"github.com/google/uuid"
)

type ReportInterface struct {
	app *app.Application
}

func NewReportInterface(app *app.Application) ReportInterface {
	return ReportInterface{app: app}
}

func (i ReportInterface) Initialize(ctx context.Context, cmd command.InitializeCmd) error {
	return i.app.Commands.Initialize.Handle(ctx, cmd)
}

func (i ReportInterface) GenerateForPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) error {
	return i.app.Commands.GenerateForPeriod.Handle(ctx, command.GenerateForPeriodCmd{
		SobId:    sobId,
		PeriodId: periodId,
	})
}
