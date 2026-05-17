package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/generator"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"
	"github/fims-proto/fims-proto-ms/internal/report/domain/validator"

	"github.com/google/uuid"
)

type GenerateReportCmd struct {
	TemplateId       uuid.UUID
	ReportId         uuid.UUID
	SobId            uuid.UUID
	Title            string
	AmountTypes      []string
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

func (h GenerateHandler) Handle(ctx context.Context, cmd GenerateReportCmd) (uuid.UUID, error) {
	var actualId uuid.UUID
	err := h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		id, err := h.handle(txCtx, cmd)
		if err != nil {
			return err
		}
		actualId = id
		return nil
	})
	return actualId, err
}

func (h GenerateHandler) handle(ctx context.Context, cmd GenerateReportCmd) (uuid.UUID, error) {
	reportTemplate, err := h.repo.ReadReportById(ctx, cmd.TemplateId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to read report: %w", err)
	}

	periodId, err := h.generalLedgerService.ReadPeriodIdByFiscalYearAndNumber(
		ctx,
		cmd.SobId,
		cmd.PeriodFiscalYear,
		cmd.PeriodNumber,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to read period: %w", err)
	}

	existing, err := h.repo.ReadInstanceBySobClassAndPeriod(ctx, cmd.SobId, reportTemplate.Class(), periodId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to check existing report instance: %w", err)
	}

	if existing != nil {
		v := validator.NewValidatorFactory(existing.Class())
		updateErr := h.repo.UpdateReport(ctx, existing.Id(), func(r *report.Report) (*report.Report, error) {
			g := generator.NewGenerator(r, h.generalLedgerService, v)
			if err = g.Regenerate(ctx); err != nil {
				return nil, fmt.Errorf("failed to regenerate report: %w", err)
			}
			return g.Report(), nil
		})
		return existing.Id(), updateErr
	}

	v := validator.NewValidatorFactory(reportTemplate.Class())
	reportGenerator := generator.NewGenerator(reportTemplate, h.generalLedgerService, v)
	newReport, err := reportGenerator.Generate(ctx, cmd.ReportId, periodId, cmd.Title, cmd.AmountTypes)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate report: %w", err)
	}

	if err = h.repo.CreateReports(ctx, []*report.Report{newReport}); err != nil {
		return uuid.Nil, fmt.Errorf("failed to save report: %w", err)
	}

	return cmd.ReportId, nil
}
