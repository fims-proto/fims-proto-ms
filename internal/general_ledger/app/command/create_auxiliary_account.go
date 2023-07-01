package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_ledger"
)

type CreateAuxiliaryAccountCmd struct {
	AccountId   uuid.UUID
	CategoryId  uuid.UUID
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
		auxiliaryAccount, err := h.createAccount(txCtx, cmd)
		if err != nil {
			return err
		}

		// create auxiliary ledger
		return h.createLedger(txCtx, auxiliaryAccount)
	})
}

func (h CreateAuxiliaryAccountHandler) createAccount(ctx context.Context, cmd CreateAuxiliaryAccountCmd) (*auxiliary_account.AuxiliaryAccount, error) {
	category, err := h.repo.ReadAuxiliaryAccountCategoryById(ctx, cmd.CategoryId)
	if err != nil {
		return nil, fmt.Errorf("failed to get auxiliary account category: %w", err)
	}

	auxiliaryAccount, err := auxiliary_account.New(
		cmd.AccountId,
		category,
		cmd.Key,
		cmd.Title,
		cmd.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create auxiliary account: %w", err)
	}

	if err = h.repo.CreateAuxiliaryAccounts(ctx, utils.AsSlice(auxiliaryAccount)); err != nil {
		return nil, err
	}

	return auxiliaryAccount, nil
}

func (h CreateAuxiliaryAccountHandler) createLedger(ctx context.Context, auxiliaryAccount *auxiliary_account.AuxiliaryAccount) error {
	p, err := h.repo.ReadCurrentPeriod(ctx, auxiliaryAccount.Category().SobId())
	if err != nil {
		return fmt.Errorf("failed to read current period: %w", err)
	}

	auxiliaryLedger, err := auxiliary_ledger.New(
		uuid.New(),
		p.Id(),
		auxiliaryAccount,
		decimal.Zero,
		decimal.Zero,
		decimal.Zero,
		decimal.Zero,
	)
	if err != nil {
		return fmt.Errorf("failed to create auxiliary ledger: %w", err)
	}

	return h.repo.CreateAuxiliaryLedgers(ctx, utils.AsSlice(auxiliaryLedger))
}
