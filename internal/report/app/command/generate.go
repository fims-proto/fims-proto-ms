package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/evaluator"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"

	"github.com/google/uuid"
)

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
	return GenerateHandler{repo: repo, generalLedgerService: generalLedgerService}
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
	periodId, err := h.generalLedgerService.ReadPeriodIdByFiscalYearAndNumber(ctx, cmd.SobId, cmd.FiscalYear, cmd.PeriodNumber)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to read period: %w", err)
	}

	existing, err := h.repo.ReadInstanceBySobClassAndPeriod(ctx, cmd.SobId, cmd.Class, periodId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to check existing report instance: %w", err)
	}
	if existing != nil {
		return existing.Id(), nil
	}

	tmpl, err := h.repo.ReadTemplateBySobIdAndClass(ctx, cmd.SobId, cmd.Class)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to read report template: %w", err)
	}
	if tmpl == nil {
		return uuid.Nil, fmt.Errorf("report template %s not found", cmd.Class)
	}

	newReport, err := tmpl.Instantiate(uuid.New(), periodId, "")
	if err != nil {
		return uuid.Nil, err
	}
	ev := evaluator.New(newReport, h.generalLedgerService)
	if err = ev.Regenerate(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate report: %w", err)
	}
	if err = h.repo.CreateReports(ctx, []*report.Report{ev.Report()}); err != nil {
		return uuid.Nil, fmt.Errorf("failed to save report: %w", err)
	}
	return ev.Report().Id(), nil
}
