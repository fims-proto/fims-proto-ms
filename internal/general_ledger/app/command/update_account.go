package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
	sobQuery "github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github.com/google/uuid"
)

type UpdateAccountCmd struct {
	AccountId            uuid.UUID
	SobId                uuid.UUID
	Title                string
	LevelNumber          int
	BalanceDirection     string
	Group                int
	DimensionCategoryIds []uuid.UUID
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
	// Fetch SoB only if we're updating the account number (for code length validation)
	var sob *sobQuery.Sob
	if cmd.LevelNumber != 0 {
		s, err := h.sobService.ReadById(ctx, cmd.SobId)
		if err != nil {
			return fmt.Errorf("failed to read sob: %w", err)
		}
		sob = &s
	}

	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.update(txCtx, cmd, sob)
	})
}

func (h UpdateAccountHandler) update(ctx context.Context, cmd UpdateAccountCmd, sob *sobQuery.Sob) error {
	return h.repo.UpdateAccount(ctx, cmd.AccountId, func(a *account.Account) (*account.Account, error) {
		if cmd.Title != "" {
			if err := a.UpdateTitle(cmd.Title); err != nil {
				return nil, fmt.Errorf("failed to update title: %w", err)
			}
		}

		if cmd.LevelNumber != 0 {
			// Validate: levelNumber string length must not exceed code length for this level
			levelNumberStr := fmt.Sprintf("%d", cmd.LevelNumber)
			codeLength := sob.AccountsCodeLength[a.Level()-1]
			if len(levelNumberStr) > codeLength {
				return nil, fmt.Errorf("level number %d (length %d) exceeds code length %d for level %d",
					cmd.LevelNumber, len(levelNumberStr), codeLength, a.Level())
			}

			if err := a.UpdateNumber(cmd.LevelNumber); err != nil {
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

		// DimensionCategoryIds is always applied (nil means "clear all", empty slice also clears).
		// Callers should only set this field when they intend to update dimension bindings.
		if cmd.DimensionCategoryIds != nil {
			a.UpdateDimensionCategories(cmd.DimensionCategoryIds)
		}

		return a, nil
	})
}
