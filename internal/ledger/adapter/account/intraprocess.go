package account

import (
	"context"

	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"

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

func (a IntraProcessAdapter) ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]accountQuery.Account, error) {
	return a.accountInterface.ReadAccountsByIds(ctx, accountIds)
}

func (a IntraProcessAdapter) ReadAllAccountsBySobId(ctx context.Context, sobId uuid.UUID) ([]accountQuery.Account, error) {
	return a.accountInterface.ReadAllAccountsBySobId(ctx, sobId)
}
