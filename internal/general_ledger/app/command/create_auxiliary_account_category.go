package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account_category"
)

type CreateAuxiliaryAccountCategoryCmd struct {
	SobId      uuid.UUID
	CategoryId uuid.UUID
	Key        string
	Title      string
	IsStandard bool
}

type CreateAuxiliaryAccountCategoryHandler struct {
	repo domain.Repository
}

func NewCreateAuxiliaryAccountCategoryHandler(repo domain.Repository) CreateAuxiliaryAccountCategoryHandler {
	if repo == nil {
		panic("nil repo")
	}

	return CreateAuxiliaryAccountCategoryHandler{repo: repo}
}

func (h CreateAuxiliaryAccountCategoryHandler) Handle(ctx context.Context, cmd CreateAuxiliaryAccountCategoryCmd) error {
	auxiliaryAccountCategory, err := auxiliary_account_category.New(
		cmd.CategoryId,
		cmd.SobId,
		cmd.Key,
		cmd.Title,
		cmd.IsStandard,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create auxiliary account category")
	}

	return h.repo.CreateAuxiliaryAccountCategories(ctx, utils.AsSlice(auxiliaryAccountCategory))
}
