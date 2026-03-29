package app

import (
	"github/fims-proto/fims-proto-ms/internal/report/app/command"
	"github/fims-proto/fims-proto-ms/internal/report/app/query"
	"github/fims-proto/fims-proto-ms/internal/report/app/service"
	"github/fims-proto/fims-proto-ms/internal/report/domain"
	domainService "github/fims-proto/fims-proto-ms/internal/report/domain/service"
)

type Queries struct {
	PagingReports query.PagingReportsHandler
	ReportById    query.ReportByIdHandler
}

type Commands struct {
	Initialize command.InitializeHandler

	Generate   command.GenerateHandler
	Regenerate command.RegenerateHandler

	UpdateReport command.UpdateReportHandler

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
	generalLedgerService domainService.GeneralLedgerService,
	sobService service.SobService,
) {
	a.Queries = Queries{
		PagingReports: query.NewPagingReportsHandler(readModel),
		ReportById:    query.NewReportByIdHandler(readModel),
	}
	a.Commands = Commands{
		Initialize: command.NewInitializeHandler(repo, generalLedgerService, sobService),

		Generate:   command.NewGenerateHandler(repo, generalLedgerService),
		Regenerate: command.NewRegenerateHandler(repo, generalLedgerService),

		UpdateReport: command.NewUpdateReportHandler(repo, generalLedgerService, sobService),

		Migrate: command.NewMigrationHandler(repo),
	}
}
