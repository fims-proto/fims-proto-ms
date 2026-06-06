package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/dimension/domain"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain/category"

	"github.com/google/uuid"
)

type UpdateCategoryCmd struct {
	CategoryId uuid.UUID
	NewName    string
}

type UpdateCategoryHandler struct {
	repo domain.Repository
}

func NewUpdateCategoryHandler(repo domain.Repository) UpdateCategoryHandler {
	if repo == nil {
		panic("nil repo")
	}

	return UpdateCategoryHandler{repo: repo}
}

func (h UpdateCategoryHandler) Handle(ctx context.Context, cmd UpdateCategoryCmd) error {
	return h.repo.UpdateCategory(ctx, cmd.CategoryId, func(c *category.DimensionCategory) (*category.DimensionCategory, error) {
		if err := c.Rename(cmd.NewName); err != nil {
			return nil, fmt.Errorf("failed to rename dimension category: %w", err)
		}

		return c, nil
	})
}
