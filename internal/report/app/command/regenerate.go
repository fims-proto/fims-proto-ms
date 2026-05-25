package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/evaluator"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"
)

type RegenerateHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewRegenerateHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) RegenerateHandler {
	if repo == nil {
		panic("nil repository")
	}
	if generalLedgerService == nil {
		panic("nil general ledger service")
	}
	return RegenerateHandler{repo: repo, generalLedgerService: generalLedgerService}
}

func (h RegenerateHandler) Handle(ctx context.Context, cmd RegenerateReportCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.repo.UpdateReport(txCtx, cmd.ReportId, func(existing *report.Report) (*report.Report, error) {
			if existing.Template() {
				return nil, fmt.Errorf("report template cannot be regenerated")
			}
			tmpl, err := h.repo.ReadTemplateBySobIdAndClass(txCtx, existing.SobId(), existing.Class())
			if err != nil {
				return nil, fmt.Errorf("failed to read report template: %w", err)
			}
			if tmpl == nil {
				return nil, fmt.Errorf("report template %s not found", existing.Class())
			}
			rebuilt, err := tmpl.Instantiate(existing.Id(), existing.PeriodId(), existing.Title())
			if err != nil {
				return nil, err
			}
			ev := evaluator.New(rebuilt, h.generalLedgerService)
			if err = ev.Regenerate(txCtx); err != nil {
				return nil, fmt.Errorf("failed to regenerate report: %w", err)
			}
			return ev.Report(), nil
		})
	})
}
