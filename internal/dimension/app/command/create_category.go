package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/dimension/domain"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain/category"

	"github.com/google/uuid"
)

type CreateCategoryCmd struct {
	CategoryId uuid.UUID
	SobId      uuid.UUID
	Name       string
}

type CreateCategoryHandler struct {
	repo domain.Repository
}

func NewCreateCategoryHandler(repo domain.Repository) CreateCategoryHandler {
	if repo == nil {
		panic("nil repo")
	}

	return CreateCategoryHandler{repo: repo}
}

func (h CreateCategoryHandler) Handle(ctx context.Context, cmd CreateCategoryCmd) error {
	c, err := category.New(cmd.CategoryId, cmd.SobId, cmd.Name)
	if err != nil {
		return fmt.Errorf("failed to create dimension category: %w", err)
	}

	return h.repo.CreateCategory(ctx, c)
}
