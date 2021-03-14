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

func (s IntraprocessAdapter) ValidateExistence(ctx context.Context, accNumbers []string) error {
	return s.accountInterface.ValidateExistence(ctx, accNumbers)
}
