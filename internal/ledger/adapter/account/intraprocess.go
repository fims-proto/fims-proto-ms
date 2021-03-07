package account

import (
	"context"
	accountport "github/fims-proto/fims-proto-ms/internal/account/port/private/intraprocess"
)

type IntraprocessAdapter struct {
	accInterface accountport.AccountInterface
}

func NewIntraprocessAdapter(accInterface accountport.AccountInterface) IntraprocessAdapter {
	return IntraprocessAdapter{accInterface: accInterface}
}

func (s IntraprocessAdapter) ReadSuperiorNumbers(ctx context.Context, accountNumber string) ([]string, error) {
	return s.accInterface.ReadSuperiorNumbers(ctx, accountNumber)
}
