package command

import (
	"context"
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain"

	"github.com/google/uuid"
)

type DeleteCategoryCmd struct {
	CategoryId uuid.UUID
}

type DeleteCategoryHandler struct {
	repo domain.Repository
}

func NewDeleteCategoryHandler(repo domain.Repository) DeleteCategoryHandler {
	if repo == nil {
		panic("nil repo")
	}

	return DeleteCategoryHandler{repo: repo}
}

func (h DeleteCategoryHandler) Handle(ctx context.Context, cmd DeleteCategoryCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		used, err := h.repo.ExistsCategoryUsedByJournalLine(txCtx, cmd.CategoryId)
		if err != nil {
			return fmt.Errorf("failed to check category usage: %w", err)
		}

		if used {
			return commonErrors.NewSlugError("dimension-deleteCategory-hasUsedOptions")
		}

		// Cascade-delete options (safe because none are used).
		if err := h.repo.DeleteOptionsByCategoryId(txCtx, cmd.CategoryId); err != nil {
			return fmt.Errorf("failed to delete options for category: %w", err)
		}

		return h.repo.DeleteCategory(txCtx, cmd.CategoryId)
	})
}
