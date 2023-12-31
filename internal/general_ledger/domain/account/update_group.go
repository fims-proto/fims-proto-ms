package account

import "github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/class"

func (a *Account) UpdateGroup(group int) error {
	if err := class.Validate(a.class, class.Group(group)); err != nil {
		return err
	}

	a.group = class.Group(group)
	return nil
}
