package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"
)

type CreateAuxiliaryCategoryCmd struct {
	SobId      uuid.UUID
	CategoryId uuid.UUID
	Key        string
	Title      string
	IsStandard bool
}

type CreateAuxiliaryCategoryHandler struct {
	repo domain.Repository
}

func NewCreateAuxiliaryCategoryHandler(repo domain.Repository) CreateAuxiliaryCategoryHandler {
	if repo == nil {
		panic("nil repo")
	}

	return CreateAuxiliaryCategoryHandler{repo: repo}
}

func (h CreateAuxiliaryCategoryHandler) Handle(ctx context.Context, cmd CreateAuxiliaryCategoryCmd) error {
	auxiliaryCategory, err := auxiliary_category.New(
		cmd.CategoryId,
		cmd.SobId,
		cmd.Key,
		cmd.Title,
		cmd.IsStandard,
	)
	if err != nil {
		return fmt.Errorf("failed to create auxiliary category: %w", err)
	}

	return h.repo.CreateAuxiliaryCategories(ctx, utils.AsSlice(auxiliaryCategory))
}
