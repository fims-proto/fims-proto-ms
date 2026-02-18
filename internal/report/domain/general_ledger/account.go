package general_ledger

import "github/fims-proto/fims-proto-ms/internal/report/domain/general_ledger/balance_direction"

type Account struct {
	balanceDirection balance_direction.BalanceDirection
}

func NewAccount(balanceDirection string) (*Account, error) {
	newBalanceDirection, err := balance_direction.FromString(balanceDirection)
	if err != nil {
		return nil, err
	}

	return &Account{
		balanceDirection: newBalanceDirection,
	}, nil
}

func (a Account) BalanceDirection() balance_direction.BalanceDirection {
	return a.balanceDirection
}
