package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/generator"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"
)

type GenerateReportCmd struct {
	TemplateId       uuid.UUID
	ReportId         uuid.UUID
	SobId            uuid.UUID
	PeriodFiscalYear int
	PeriodNumber     int
}

type GenerateHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewGenerateHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) GenerateHandler {
	if repo == nil {
		panic("nil repository")
	}

	if generalLedgerService == nil {
		panic("nil general ledger service")
	}
	return GenerateHandler{
		repo:                 repo,
		generalLedgerService: generalLedgerService,
	}
}

func (h GenerateHandler) Handle(ctx context.Context, cmd GenerateReportCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.handle(txCtx, cmd)
	})
}

func (h GenerateHandler) handle(ctx context.Context, cmd GenerateReportCmd) error {
	// read template first
	reportTemplate, err := h.repo.ReadReportById(ctx, cmd.TemplateId)
	if err != nil {
		return fmt.Errorf("failed to read report: %w", err)
	}

	periodId, err := h.generalLedgerService.ReadPeriodIdByFiscalYearAndNumber(
		ctx,
		cmd.SobId,
		cmd.PeriodFiscalYear,
		cmd.PeriodNumber,
	)
	if err != nil {
		return fmt.Errorf("failed to read period: %w", err)
	}

	reportGenerator := generator.NewGenerator(reportTemplate, h.generalLedgerService)
	newReport, err := reportGenerator.Generate(ctx, cmd.ReportId, periodId)
	if err != nil {
		return fmt.Errorf("failed to generate report: %w", err)
	}

	// save
	err = h.repo.CreateReports(ctx, []*report.Report{newReport})
	if err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	return nil
}
