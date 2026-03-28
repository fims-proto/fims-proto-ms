package command

import (
	"context"
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"

	"github.com/google/uuid"
)

type DeleteAccountCmd struct {
	AccountId uuid.UUID
	SobId     uuid.UUID
}

type DeleteAccountHandler struct {
	repo domain.Repository
}

func NewDeleteAccountHandler(repo domain.Repository) DeleteAccountHandler {
	if repo == nil {
		panic("nil repo")
	}

	return DeleteAccountHandler{repo: repo}
}

func (h DeleteAccountHandler) Handle(ctx context.Context, cmd DeleteAccountCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		acct, err := h.repo.ReadAccountById(txCtx, cmd.AccountId)
		if err != nil {
			return fmt.Errorf("failed to read account: %w", err)
		}

		hasChildren, err := h.repo.ExistsChildAccountsByAccountId(txCtx, cmd.AccountId)
		if err != nil {
			return fmt.Errorf("failed to check child accounts: %w", err)
		}
		if hasChildren {
			return commonErrors.NewSlugError("account-delete-hasChildren")
		}

		usedByJournalLine, err := h.repo.ExistsJournalLinesByAccountId(txCtx, cmd.AccountId)
		if err != nil {
			return fmt.Errorf("failed to check journal line usage: %w", err)
		}
		if usedByJournalLine {
			return commonErrors.NewSlugError("account-delete-usedByJournalLine")
		}

		hasOpeningBalance, err := h.repo.ExistsLedgerWithOpeningBalanceByAccountId(txCtx, cmd.AccountId)
		if err != nil {
			return fmt.Errorf("failed to check ledger opening balance: %w", err)
		}
		if hasOpeningBalance {
			return commonErrors.NewSlugError("account-delete-hasOpeningBalance")
		}

		if err := h.repo.DeleteLedgersByAccountId(txCtx, cmd.AccountId); err != nil {
			return fmt.Errorf("failed to delete ledgers for account: %w", err)
		}

		if err := h.repo.DeleteAccount(txCtx, cmd.AccountId); err != nil {
			return fmt.Errorf("failed to delete account: %w", err)
		}

		// Restore parent's isLeaf flag if it has no remaining children
		superiorAccountId := acct.SuperiorAccountId()
		if superiorAccountId == uuid.Nil {
			return nil
		}

		stillHasChildren, err := h.repo.ExistsChildAccountsByAccountId(txCtx, superiorAccountId)
		if err != nil {
			return fmt.Errorf("failed to check remaining child accounts: %w", err)
		}
		if stillHasChildren {
			return nil
		}

		return h.repo.UpdateAccount(txCtx, superiorAccountId, func(a *account.Account) (*account.Account, error) {
			if err := a.UpdateLeaf(true); err != nil {
				return nil, fmt.Errorf("failed to restore superior account leaf indicator: %w", err)
			}
			return a, nil
		})
	})
}
