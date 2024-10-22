package intraprocess

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/app"
	"github/fims-proto/fims-proto-ms/internal/report/app/command"
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
