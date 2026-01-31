package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/domain"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report"

	"github.com/google/uuid"
)

type DeleteItemCmd struct {
	ReportId  uuid.UUID
	SectionId uuid.UUID
	ItemId    uuid.UUID
}

type DeleteItemHandler struct {
	repo domain.Repository
}

func NewDeleteItemHandler(repo domain.Repository) DeleteItemHandler {
	if repo == nil {
		panic("nil repo")
	}

	return DeleteItemHandler{
		repo: repo,
	}
}

func (h DeleteItemHandler) Handle(ctx context.Context, cmd DeleteItemCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.deleteItem(txCtx, cmd)
	})
}

func (h DeleteItemHandler) deleteItem(ctx context.Context, cmd DeleteItemCmd) error {
	return h.repo.UpdateReport(ctx, cmd.ReportId, func(r *report.Report) (*report.Report, error) {
		// Delete item from the specified section
		if err := r.DeleteItemFromSection(cmd.SectionId, cmd.ItemId); err != nil {
			return nil, err
		}

		return r, nil
	})
}
