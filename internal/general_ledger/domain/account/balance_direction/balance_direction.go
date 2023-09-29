package balance_direction

import (
	"fmt"
)

type BalanceDirection struct {
	slug string
}

func (b BalanceDirection) String() string {
	return b.slug
}

var (
	Unknown    = BalanceDirection{""}
	Debit      = BalanceDirection{"debit"}
	Credit     = BalanceDirection{"credit"}
	NotDefined = BalanceDirection{"not_defined"}
)

func FromString(s string) (BalanceDirection, error) {
	switch s {
	case Debit.slug:
		return Debit, nil
	case Credit.slug:
		return Credit, nil
	case NotDefined.slug:
		return NotDefined, nil
	}

	return Unknown, fmt.Errorf("unknown balance direction: %s", s)
}
