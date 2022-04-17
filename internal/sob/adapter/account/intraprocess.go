package account

import (
	"context"

	accountPort "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	accountInterface accountPort.AccountInterface
}

func NewIntraProcessAdapter(accountInterface accountPort.AccountInterface) IntraProcessAdapter {
	return IntraProcessAdapter{accountInterface: accountInterface}
}

func (i IntraProcessAdapter) InitializeAccounts(ctx context.Context, sobId uuid.UUID) error {
	return i.accountInterface.InitializeAccounts(ctx, sobId)
}
