package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/generator"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"
	"github/fims-proto/fims-proto-ms/internal/report/domain/service"

	"github.com/google/uuid"
)

type RegenerateReportCmd struct {
	ReportId uuid.UUID
}

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

	return RegenerateHandler{
		repo:                 repo,
		generalLedgerService: generalLedgerService,
	}
}

func (h RegenerateHandler) Handle(ctx context.Context, cmd RegenerateReportCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.regenerate(txCtx, cmd)
	})
}

func (h RegenerateHandler) regenerate(ctx context.Context, cmd RegenerateReportCmd) error {
	return h.repo.UpdateReport(
		ctx,
		cmd.ReportId,
		func(r *report.Report) (*report.Report, error) {
			g := generator.NewGenerator(r, h.generalLedgerService)
			if err := g.Regenerate(ctx); err != nil {
				return nil, fmt.Errorf("failed to regenerate report: %w", err)
			}

			return g.Report(), nil
		},
	)
}
