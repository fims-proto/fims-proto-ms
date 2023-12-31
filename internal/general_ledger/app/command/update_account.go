package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"
)

type UpdateAccountCmd struct {
	AccountId        uuid.UUID
	SobId            uuid.UUID
	Title            string
	LevelNumber      int
	BalanceDirection string
	Group            int
	CategoryKeys     []string
}

type UpdateAccountHandler struct {
	repo       domain.Repository
	sobService service.SobService
}

func NewUpdateAccountHandler(repo domain.Repository, sobService service.SobService) UpdateAccountHandler {
	if repo == nil {
		panic("nil repo")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	return UpdateAccountHandler{
		repo:       repo,
		sobService: sobService,
	}
}

func (h UpdateAccountHandler) Handle(ctx context.Context, cmd UpdateAccountCmd) error {
	return h.repo.UpdateAccount(ctx, cmd.AccountId, func(a *account.Account) (*account.Account, error) {
		if cmd.Title != "" {
			if err := a.UpdateTitle(cmd.Title); err != nil {
				return nil, fmt.Errorf("failed to update title: %w", err)
			}
		}

		if cmd.LevelNumber != 0 {
			sob, err := h.sobService.ReadById(ctx, cmd.SobId)
			if err != nil {
				return nil, fmt.Errorf("failed to read sob: %w", err)
			}

			if err = a.UpdateNumber(cmd.LevelNumber, sob.AccountsCodeLength); err != nil {
				return nil, fmt.Errorf("failed to update account number: %w", err)
			}
		}

		if cmd.BalanceDirection != "" {
			if err := a.UpdateBalanceDirection(cmd.BalanceDirection); err != nil {
				return nil, fmt.Errorf("failed to update balance direction: %w", err)
			}
		}

		if cmd.Group != 0 {
			if err := a.UpdateGroup(cmd.Group); err != nil {
				return nil, fmt.Errorf("failed to update group: %w", err)
			}
		}

		if cmd.CategoryKeys != nil {
			var categories []*auxiliary_category.AuxiliaryCategory
			for _, key := range cmd.CategoryKeys {
				auxiliaryCategory, err := h.repo.ReadAuxiliaryCategoryByKey(ctx, key)
				if err != nil {
					return nil, fmt.Errorf("failed to read auxiliary category: %w", err)
				}
				categories = append(categories, auxiliaryCategory)
			}

			if err := a.AssignAuxiliaryCategories(categories); err != nil {
				return nil, fmt.Errorf("failed to update auxiliary assignment: %w", err)
			}
		}

		return a, nil
	})
}
