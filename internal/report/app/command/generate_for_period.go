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

type GenerateForPeriodHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewGenerateForPeriodHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) GenerateForPeriodHandler {
	if repo == nil {
		panic("nil repository")
	}
	if generalLedgerService == nil {
		panic("nil general ledger service")
	}
	return GenerateForPeriodHandler{repo: repo, generalLedgerService: generalLedgerService}
}

func (h GenerateForPeriodHandler) Handle(ctx context.Context, cmd GenerateForPeriodCmd) error {
	templates, err := h.repo.ReadTemplatesBySobId(ctx, cmd.SobId)
	if err != nil {
		return fmt.Errorf("failed to read report templates: %w", err)
	}
	for _, tmpl := range templates {
		if err = h.generateForTemplate(ctx, tmpl, cmd.PeriodId); err != nil {
			return err
		}
	}
	return nil
}

func (h GenerateForPeriodHandler) generateForTemplate(ctx context.Context, tmpl *report.Report, periodId uuid.UUID) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		existing, err := h.repo.ReadInstanceBySobClassAndPeriod(txCtx, tmpl.SobId(), tmpl.Class(), periodId)
		if err != nil {
			return fmt.Errorf("failed to check existing report instance: %w", err)
		}
		if existing != nil {
			return nil
		}

		newReport, err := tmpl.Instantiate(uuid.New(), periodId, "")
		if err != nil {
			return err
		}
		ev := evaluator.New(newReport, h.generalLedgerService)
		if err = ev.Regenerate(txCtx); err != nil {
			return fmt.Errorf("failed to generate report: %w", err)
		}
		return h.repo.CreateReports(txCtx, []*report.Report{ev.Report()})
	})
}
