package account

import (
	"context"

	accountQuery "github/fims-proto/fims-proto-ms/internal/account/app/query"

	accountPort "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	accountInterface accountPort.AccountInterface
}

func NewIntraProcessAdapter(accountInterface accountPort.AccountInterface) IntraProcessAdapter {
	return IntraProcessAdapter{accountInterface: accountInterface}
}

func (a IntraProcessAdapter) ReadAccountsBySobId(ctx context.Context, sobId uuid.UUID) ([]accountQuery.Account, error) {
	return a.accountInterface.ReadAccountsBySobId(ctx, sobId)
}

func (a IntraProcessAdapter) ReadAccountsByIds(ctx context.Context, accountIds []uuid.UUID) (map[uuid.UUID]accountQuery.Account, error) {
	return a.accountInterface.ReadAccountsByIds(ctx, accountIds)
}

func (a IntraProcessAdapter) ReadAccountsWithSuperiorsByIds(ctx context.Context, accountIds []uuid.UUID) ([]accountQuery.Account, error) {
	return a.accountInterface.ReadAccountsWithSuperiorsByIds(ctx, accountIds)
}
