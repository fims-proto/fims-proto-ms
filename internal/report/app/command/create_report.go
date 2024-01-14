package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"

	"github.com/google/uuid"
)

type CreateReportCmd struct {
	ReportId      uuid.UUID
	SobId         uuid.UUID
	PeriodId      uuid.UUID
	Version       int
	RefTemplateId uuid.UUID
}

type CreateReportHandler struct {
	repo domain.Repository
}

func NewcreateReportHandler(repo domain.Repository) CreateReportHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CreateReportHandler{repo: repo}
}

func (h CreateReportHandler) Handle(ctx context.Context, cmd CreateReportCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.createReport(txCtx, cmd)
	})
}

func (h CreateReportHandler) createReport(ctx context.Context, cmd CreateReportCmd) error {
	innerTemplate, err := h.repo.DeepCopyTemplate(ctx, cmd.RefTemplateId)
	if err != nil {
		return err
	}
	newReport, err := report.New(
		cmd.ReportId,
		cmd.PeriodId,
		cmd.Version,
		cmd.RefTemplateId,
		innerTemplate)
	if err != nil {
		return err
	}
	return h.repo.CreateReport(ctx, newReport)
}
