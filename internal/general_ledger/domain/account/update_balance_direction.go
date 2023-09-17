package account

import (
	"fmt"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/balance_direction"
)

func (a *Account) UpdateBalanceDirection(direction string) error {
	balanceDirection, err := balance_direction.FromString(direction)
	if err != nil {
		return err
	}

	if a.superiorAccountId != uuid.Nil {
		return fmt.Errorf("balance direction cannot be update when superior account exists")
	}

	a.balanceDirection = balanceDirection

	return nil
}
