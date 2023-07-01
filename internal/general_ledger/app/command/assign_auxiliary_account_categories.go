package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account_category"
)

type AssignAuxiliaryAccountCategoryCmd struct {
	AccountId   uuid.UUID
	CategoryIds []uuid.UUID
}

type AssignAuxiliaryAccountCategoryHandler struct {
	repo domain.Repository
}

func NewAssignAuxiliaryAccountCategoryHandler(repo domain.Repository) AssignAuxiliaryAccountCategoryHandler {
	if repo == nil {
		panic("nil repo")
	}

	return AssignAuxiliaryAccountCategoryHandler{repo: repo}
}

func (h AssignAuxiliaryAccountCategoryHandler) Handle(ctx context.Context, cmd AssignAuxiliaryAccountCategoryCmd) error {
	var categories []*auxiliary_account_category.AuxiliaryAccountCategory
	for _, categoryId := range cmd.CategoryIds {
		auxiliaryAccountCategory, err := h.repo.ReadAuxiliaryAccountCategoryById(ctx, categoryId)
		if err != nil {
			return fmt.Errorf("failed to read auxiliary account category: %w", err)
		}
		categories = append(categories, auxiliaryAccountCategory)
	}

	return h.repo.UpdateAccount(ctx, cmd.AccountId, func(a *account.Account) (*account.Account, error) {
		a.AssignAuxiliaryAccountCategories(categories)
		return a, nil
	})
}
