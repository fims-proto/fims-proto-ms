package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/report/app/service"
	"github/fims-proto/fims-proto-ms/internal/report/domain"

	"github.com/google/uuid"
)

type InitializeCmd struct {
	SobId uuid.UUID
}

type InitializeHandler struct {
	repo       domain.Repository
	sobService service.SobService
}

func NewInitializeHandler(repo domain.Repository, sobService service.SobService) InitializeHandler {
	if repo == nil {
		panic("nil repo")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	return InitializeHandler{
		repo:       repo,
		sobService: sobService,
	}
}

func (h InitializeHandler) Handle(ctx context.Context, cmd InitializeCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		_, err := h.sobService.ReadById(txCtx, cmd.SobId)
		if err != nil {
			return fmt.Errorf("failed to read sob: %w", err)
		}

		// TODO initialize default template from json?

		return nil
	})
}
