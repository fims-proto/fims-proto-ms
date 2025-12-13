package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"

	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/class"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateAccountCmd struct {
	AccountId             uuid.UUID
	SobId                 uuid.UUID
	Title                 string
	LevelNumber           int
	SuperiorAccountNumber string
	BalanceDirection      string
	Class                 int
	Group                 int
	CategoryKeys          []string
}

type CreateAccountHandler struct {
	repo       domain.Repository
	sobService service.SobService
}

func NewCreateAccountHandler(repo domain.Repository, sobService service.SobService) CreateAccountHandler {
	if repo == nil {
		panic("nil repo")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	return CreateAccountHandler{
		repo:       repo,
		sobService: sobService,
	}
}

func (h CreateAccountHandler) Handle(ctx context.Context, cmd CreateAccountCmd) error {
	accountClass := class.Class(cmd.Class)
	accountGroup := class.Group(cmd.Group)
	superiorAccountId := uuid.Nil
	level := 1
	numberHierarchy := []int{cmd.LevelNumber}

	sob, err := h.sobService.ReadById(ctx, cmd.SobId)
	if err != nil {
		return fmt.Errorf("failed to read sob: %w", err)
	}

	if err := class.Validate(accountClass, accountGroup); err != nil {
		return fmt.Errorf("invalid class or group: %w", err)
	}

	if cmd.SuperiorAccountNumber != "" {
		superiorAccount, err := h.repo.ReadAccountByNumber(ctx, cmd.SobId, cmd.SuperiorAccountNumber)
		if err != nil {
			return err
		}

		if superiorAccount.Class() != accountClass {
			return fmt.Errorf("class %s does not match superior account's class %s", accountClass, superiorAccount.Class())
		}

		if superiorAccount.Group() != accountGroup {
			return fmt.Errorf("group %s does not match superior account's group %s", accountGroup, superiorAccount.Group())
		}

		if superiorAccount.Level()+1 > len(sob.AccountsCodeLength) {
			return fmt.Errorf("level %d exceeds limit %d", superiorAccount.Level()+1, len(sob.AccountsCodeLength))
		}

		superiorAccountId = superiorAccount.Id()
		level = superiorAccount.Level() + 1
		numberHierarchy = append(superiorAccount.NumberHierarchy(), cmd.LevelNumber)
	}

	categories, err := h.repo.ReadAuxiliaryCategoriesByKeys(ctx, cmd.SobId, cmd.CategoryKeys)
	if err != nil {
		return fmt.Errorf("failed to read auxiliary categories: %w", err)
	}

	newAccount, err := account.New(
		cmd.AccountId,
		cmd.SobId,
		superiorAccountId,
		cmd.Title,
		numberHierarchy,
		sob.AccountsCodeLength,
		level,
		true,
		cmd.Class,
		cmd.Group,
		cmd.BalanceDirection,
		categories,
	)
	if err != nil {
		return fmt.Errorf("failed to create new account: %w", err)
	}

	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		// create account
		if err := h.createAccount(txCtx, superiorAccountId, newAccount); err != nil {
			return err
		}

		// create ledger for current period
		return h.createLedger(txCtx, newAccount)
	})
}

func (h CreateAccountHandler) createAccount(ctx context.Context, superiorAccountId uuid.UUID, newAccount *account.Account) error {
	if superiorAccountId != uuid.Nil {
		// update superior account to not be a leaf
		if err := h.repo.UpdateAccount(ctx, superiorAccountId, func(a *account.Account) (*account.Account, error) {
			if err := a.UpdateLeaf(false); err != nil {
				return nil, fmt.Errorf("failed to update superior account leaf indicator: %w", err)
			}
			return a, nil
		}); err != nil {
			return fmt.Errorf("failed to update superior account: %w", err)
		}
	}

	// save new account
	if err := h.repo.CreateAccount(ctx, newAccount); err != nil {
		return fmt.Errorf("failed to create new account: %w", err)
	}
	return nil
}

func (h CreateAccountHandler) createLedger(ctx context.Context, acct *account.Account) error {
	p, err := h.repo.ReadCurrentPeriod(ctx, acct.SobId())
	if err != nil {
		return fmt.Errorf("failed to read current period: %w", err)
	}

	ledgerEntity, err := ledger.New(
		uuid.New(),
		acct.SobId(),
		p.Id(),
		acct.Id(),
		acct,
		decimal.Zero,
		decimal.Zero,
		decimal.Zero,
		decimal.Zero,
		decimal.Zero,
		decimal.Zero,
	)
	if err != nil {
		return fmt.Errorf("failed to create ledger: %w", err)
	}

	return h.repo.CreateLedgers(ctx, utils.AsSlice(ledgerEntity))
}
