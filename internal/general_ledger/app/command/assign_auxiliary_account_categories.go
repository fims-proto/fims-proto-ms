package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"
)

type AssignAuxiliaryCategoryCmd struct {
	AccountId    uuid.UUID
	CategoryKeys []string
}

type AssignAuxiliaryCategoryHandler struct {
	repo domain.Repository
}

func NewAssignAuxiliaryCategoryHandler(repo domain.Repository) AssignAuxiliaryCategoryHandler {
	if repo == nil {
		panic("nil repo")
	}

	return AssignAuxiliaryCategoryHandler{repo: repo}
}

func (h AssignAuxiliaryCategoryHandler) Handle(ctx context.Context, cmd AssignAuxiliaryCategoryCmd) error {
	var categories []*auxiliary_category.AuxiliaryCategory
	for _, key := range cmd.CategoryKeys {
		auxiliaryCategory, err := h.repo.ReadAuxiliaryCategoryByKey(ctx, key)
		if err != nil {
			return fmt.Errorf("failed to read auxiliary category: %w", err)
		}
		categories = append(categories, auxiliaryCategory)
	}

	return h.repo.UpdateAccount(ctx, cmd.AccountId, func(a *account.Account) (*account.Account, error) {
		a.AssignAuxiliaryCategories(categories)
		return a, nil
	})
}
