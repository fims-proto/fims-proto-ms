package account

import (
	"context"
	accountport "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
)

type IntraprocessAdapter struct {
	accountInterface accountport.AccountInterface
}

func NewIntraprocessAdapter(accountInterface accountport.AccountInterface) IntraprocessAdapter {
	return IntraprocessAdapter{accountInterface: accountInterface}
}

func (s IntraprocessAdapter) ReadSuperiorNumbers(ctx context.Context, accountNumber string) ([]string, error) {
	return s.accountInterface.ReadSuperiorNumbers(ctx, accountNumber)
}
