package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/evaluator"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"
)

type RecalculateHandler struct {
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewRecalculateHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) RecalculateHandler {
	if repo == nil {
		panic("nil repository")
	}
	if generalLedgerService == nil {
		panic("nil general ledger service")
	}
	return RecalculateHandler{repo: repo, generalLedgerService: generalLedgerService}
}

func (h RecalculateHandler) Handle(ctx context.Context, cmd RecalculateReportCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.repo.UpdateReport(txCtx, cmd.ReportId, func(r *report.Report) (*report.Report, error) {
			if r.Template() {
				return nil, fmt.Errorf("report template cannot be recalculated")
			}
			ev := evaluator.New(r, h.generalLedgerService)
			if err := ev.Regenerate(txCtx); err != nil {
				return nil, fmt.Errorf("failed to recalculate report: %w", err)
			}
			return ev.Report(), nil
		})
	})
}
