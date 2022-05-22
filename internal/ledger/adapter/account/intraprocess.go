package account

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"

	accountPort "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type IntraProcessAdapter struct {
	accountInterface accountPort.AccountInterface
}

func NewIntraProcessAdapter(accountInterface accountPort.AccountInterface) IntraProcessAdapter {
	return IntraProcessAdapter{accountInterface: accountInterface}
}

func (a IntraProcessAdapter) ReadSuperiorAccountIds(ctx context.Context, accountId uuid.UUID) ([]uuid.UUID, error) {
	account, err := a.accountInterface.ReadAccountById(ctx, accountId)
	if err != nil {
		return nil, errors.Wrap(err, "read account by id failed")
	}

	currentAccount := &account
	var superiorAccountIds []uuid.UUID
	for currentAccount.SuperiorAccount != nil {
		superiorAccountIds = append(superiorAccountIds, currentAccount.SuperiorAccountId)
		currentAccount = currentAccount.SuperiorAccount
	}
	return superiorAccountIds, nil
}

func (a IntraProcessAdapter) ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]query.Account, error) {
	return a.accountInterface.ReadAccountsByIds(ctx, accountIds)
}

func (a IntraProcessAdapter) ReadAllAccountIdsBySobId(ctx context.Context, sobId uuid.UUID) ([]uuid.UUID, error) {
	accounts, err := a.accountInterface.ReadAllAccountIdsBySobId(ctx, sobId)
	if err != nil {
		return nil, errors.Wrap(err, "read accounts by sob failed")
	}

	var accountIds []uuid.UUID
	for _, account := range accounts {
		accountIds = append(accountIds, account.Id)
	}
	return accountIds, nil
}
