package balance_direction

import "github.com/pkg/errors"

type BalanceDirection struct {
	slug string
}

func (b BalanceDirection) String() string {
	return b.slug
}

var (
	Unknown = BalanceDirection{""}
	Debit   = BalanceDirection{"debit"}
	Credit  = BalanceDirection{"credit"}
)

func FromString(s string) (BalanceDirection, error) {
	switch s {
	case Debit.slug:
		return Debit, nil
	case Credit.slug:
		return Credit, nil
	}

	return Unknown, errors.Errorf("unknown balance direction: %s", s)
}
