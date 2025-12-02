package account

import (
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/balance_direction"

	"github.com/google/uuid"
)

func (a *Account) UpdateBalanceDirection(direction string) error {
	balanceDirection, err := balance_direction.FromString(direction)
	if err != nil {
		return err
	}

	if a.superiorAccountId != uuid.Nil && a.balanceDirection != balanceDirection {
		return fmt.Errorf("balance direction cannot be update when superior account exists")
	}

	a.balanceDirection = balanceDirection

	return nil
}
