package account

import "github.com/pkg/errors"

// enum
type Type uint

const (
	Invalid       = Type(iota) // invalid accoun type
	Assets                     // assets 资产类
	Cost                       // cost 成本类
	Liabilities                // liabilities 负债类
	Equity                     // equity 所有者权益类
	ProfitAndLoss              // profit_and_loss 损益类
	Common                     // common 共同类
)

var availableTypes = map[Type]string{
	Assets:        "Assets",
	Cost:          "Cost",
	Liabilities:   "Liabilities",
	Equity:        "Equity",
	ProfitAndLoss: "ProfitAndLoss",
	Common:        "Common",
}

func NewAccountType(t Type) (Type, error) {
	for k := range availableTypes {
		if k == t {
			return k, nil
		}
	}

	return Invalid, errors.Errorf("invalid account Type: '%d'", t)
}

func NewAccountTypeFromString(s string) (Type, error) {
	for i, v := range availableTypes {
		if v == s {
			return i, nil
		}
	}

	return Invalid, errors.Errorf("invalid account name: '%s'", s)
}

func (t Type) String() string {
	for k, v := range availableTypes {
		if k == t {
			return v
		}
	}
	panic("account type string error, should not happen")
}
