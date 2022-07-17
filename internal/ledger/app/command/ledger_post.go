package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
)

type PostLedgersAmountCmd struct {
	Debit  decimal.Decimal
	Credit decimal.Decimal
}

type PostLedgersCmd struct {
	Accounts map[uuid.UUID]PostLedgersAmountCmd
	PeriodId uuid.UUID
}

type PostLedgersHandler struct {
	repo           domain.Repository
	accountService AccountService
}

func NewPostLedgersHandler(repo domain.Repository, accountService AccountService) PostLedgersHandler {
	if repo == nil {
		panic("nil ledger repository")
	}
	if accountService == nil {
		panic("nil account service")
	}
	return PostLedgersHandler{
		repo:           repo,
		accountService: accountService,
	}
}

func (h PostLedgersHandler) Handle(ctx context.Context, cmd PostLedgersCmd) error {
	log.Info(ctx, "handle post ledgers, cmd: %+v", cmd)

	// read all accounts with superiors
	var accountIds []uuid.UUID
	for accountId := range cmd.Accounts {
		accountIds = append(accountIds, accountId)
	}

	accounts, err := h.accountService.ReadAccountsWithSuperiorsByIds(ctx, accountIds)
	if err != nil {
		return errors.Wrap(err, "failed to read accounts")
	}

	// prepare account ids, and map: accountId -> amount command
	var allAccountIds []uuid.UUID
	account2AmountMap := make(map[uuid.UUID]struct {
		account accountQuery.Account
		cmd     PostLedgersAmountCmd
	})
	for _, account := range accounts {
		amountCmd, ok := cmd.Accounts[account.Id]
		if !ok {
			return errors.Errorf("failed to find account %s in command", account.Id)
		}
		accountAmount := struct {
			account accountQuery.Account
			cmd     PostLedgersAmountCmd
		}{
			account: account,
			cmd:     amountCmd,
		}
		allAccountIds = append(allAccountIds, account.Id)
		account2AmountMap[account.Id] = accountAmount

		internalAccount := account.SuperiorAccount
		for internalAccount != nil {
			allAccountIds = append(allAccountIds, internalAccount.Id)
			account2AmountMap[internalAccount.Id] = accountAmount

			internalAccount = internalAccount.SuperiorAccount
		}
	}

	return h.repo.UpdateLedgersByPeriodAndAccounts(
		ctx,
		cmd.PeriodId,
		allAccountIds,
		func(ledgers []*domain.Ledger) ([]*domain.Ledger, error) {
			for _, ledger := range ledgers {
				accountAmount, ok := account2AmountMap[ledger.AccountId()]
				if !ok {
					return nil, errors.Errorf("failed to find account %s in account-command map", ledger.AccountId())
				}

				ledger.UpdateBalance(accountAmount.cmd.Debit, accountAmount.cmd.Credit, accountAmount.account.BalanceDirection)
			}
			return ledgers, nil
		})
}
