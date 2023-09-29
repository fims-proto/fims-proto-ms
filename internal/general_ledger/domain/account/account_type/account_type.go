package account_type

import (
	"fmt"
)

type AccountType struct {
	slug string
}

func (a AccountType) String() string {
	return a.slug
}

var (
	Unknown       = AccountType{""}                // invalid account type
	Assets        = AccountType{"assets"}          // assets 资产类
	Cost          = AccountType{"cost"}            // cost 成本类
	Liabilities   = AccountType{"liabilities"}     // liabilities 负债类
	Equity        = AccountType{"equity"}          // equity 所有者权益类
	ProfitAndLoss = AccountType{"profit_and_loss"} // profit_and_loss 损益类
	Common        = AccountType{"common"}          // common 共同类
)

func FromString(s string) (AccountType, error) {
	switch s {
	case Assets.slug:
		return Assets, nil
	case Cost.slug:
		return Cost, nil
	case Liabilities.slug:
		return Liabilities, nil
	case Equity.slug:
		return Equity, nil
	case ProfitAndLoss.slug:
		return ProfitAndLoss, nil
	case Common.slug:
		return Common, nil
	}

	return Unknown, fmt.Errorf("unknown account type: %s", s)
}
