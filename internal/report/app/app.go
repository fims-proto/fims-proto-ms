package app

import (
	"github/fims-proto/fims-proto-ms/internal/report/app/command"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"
	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"
)

type Queries struct {
	PagingReports query.PagingReportsHandler
	ReportById    query.ReportByIdHandler
}

type Commands struct {
	Initialize command.InitializeHandler

	Generate   command.GenerateHandler
	Regenerate command.RegenerateHandler

	UpdateItem command.UpdateItemHandler
	AddItem    command.AddItemHandler
	DeleteItem command.DeleteItemHandler

	Migrate command.MigrationHandler
}

type Application struct {
	Queries  Queries
	Commands Commands
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(
	repo domain.Repository,
	readModel query.ReportReadModel,
	generalLedgerService service.GeneralLedgerService,
) {
	a.Queries = Queries{
		PagingReports: query.NewPagingReportsHandler(readModel),
		ReportById:    query.NewReportByIdHandler(readModel),
	}
	a.Commands = Commands{
		Initialize: command.NewInitializeHandler(repo, generalLedgerService),

		Generate:   command.NewGenerateHandler(repo, generalLedgerService),
		Regenerate: command.NewRegenerateHandler(repo, generalLedgerService),

		UpdateItem: command.NewUpdateItemHandler(repo, generalLedgerService),
		AddItem:    command.NewAddItemHandler(repo, generalLedgerService),
		DeleteItem: command.NewDeleteItemHandler(repo),

		Migrate: command.NewMigrationHandler(repo),
	}
}
