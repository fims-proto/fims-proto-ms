package command

import (
	"context"
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain"

	"github.com/google/uuid"
)

type DeleteOptionCmd struct {
	OptionId uuid.UUID
}

type DeleteOptionHandler struct {
	repo domain.Repository
}

func NewDeleteOptionHandler(repo domain.Repository) DeleteOptionHandler {
	if repo == nil {
		panic("nil repo")
	}

	return DeleteOptionHandler{repo: repo}
}

func (h DeleteOptionHandler) Handle(ctx context.Context, cmd DeleteOptionCmd) error {
	used, err := h.repo.ExistsOptionUsedByJournalLine(ctx, cmd.OptionId)
	if err != nil {
		return fmt.Errorf("failed to check option usage: %w", err)
	}

	if used {
		return commonErrors.NewSlugError("dimension-deleteOption-isUsed")
	}

	return h.repo.DeleteOption(ctx, cmd.OptionId)
}
