package account

import "github.com/pkg/errors"

type Type uint

const (
	InvalidType   = Type(iota) // invalid account type
	Assets                     // assets 资产类
	Cost                       // cost 成本类
	Liabilities                // liabilities 负债类
	Equity                     // equity 所有者权益类
	ProfitAndLoss              // profit_and_loss 损益类
	Common                     // common 共同类
)

var availableTypes = map[Type]string{
	Assets:        "assets",
	Cost:          "cost",
	Liabilities:   "liabilities",
	Equity:        "equity",
	ProfitAndLoss: "profit_and_loss",
	Common:        "common",
}

func NewAccountType(s string) (Type, error) {
	for i, v := range availableTypes {
		if v == s {
			return i, nil
		}
	}

	return InvalidType, errors.Errorf("invalid account name: '%s'", s)
}

func (t Type) String() string {
	for k, v := range availableTypes {
		if k == t {
			return v
		}
	}
	panic("account type string error, should not happen")
}
