package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"

	"github.com/google/uuid"
)

type CreateAuxiliaryAccountCmd struct {
	AccountId   uuid.UUID
	SobId       uuid.UUID
	CategoryKey string
	Key         string
	Title       string
	Description string
}

type CreateAuxiliaryAccountHandler struct {
	repo domain.Repository
}

func NewCreateAuxiliaryAccountHandler(repo domain.Repository) CreateAuxiliaryAccountHandler {
	if repo == nil {
		panic("nil repo")
	}

	return CreateAuxiliaryAccountHandler{repo: repo}
}

func (h CreateAuxiliaryAccountHandler) Handle(ctx context.Context, cmd CreateAuxiliaryAccountCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		// create auxiliary account
		return h.createAccount(txCtx, cmd)
	})
}

func (h CreateAuxiliaryAccountHandler) createAccount(ctx context.Context, cmd CreateAuxiliaryAccountCmd) error {
	category, err := h.repo.ReadAuxiliaryCategoryByKey(ctx, cmd.SobId, cmd.CategoryKey)
	if err != nil {
		return fmt.Errorf("failed to get auxiliary category: %w", err)
	}

	auxiliaryAccount, err := auxiliary_account.New(
		cmd.AccountId,
		category,
		cmd.Key,
		cmd.Title,
		cmd.Description,
	)
	if err != nil {
		return fmt.Errorf("failed to create auxiliary account: %w", err)
	}

	return h.repo.CreateAuxiliaryAccounts(ctx, utils.AsSlice(auxiliaryAccount))
}
