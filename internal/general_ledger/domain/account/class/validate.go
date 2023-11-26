package class

import (
	"cmp"
	"slices"

	"github/fims-proto/fims-proto-ms/internal/common/errors"
)

type pair struct {
	class         Class
	allowedGroups []Group
}

// allowedClasses must be sorted by class increasingly
var allowedClasses = []pair{
	{
		Assets,
		[]Group{CurrentAsset, NonCurrentAsset},
	},
	{
		Liabilities,
		[]Group{CurrentLiability, NonCurrentLiability},
	},
	{
		Equities,
		[]Group{Equity},
	},
	{
		Costs,
		[]Group{Cost},
	},
	{
		ProfitsAndLosses,
		[]Group{OperatingIncome, OtherIncome, PeriodCost, OtherCost, OperatingCostAndTax, PriorYearIncomeAdjustment, IncomeTax},
	},
	{
		Commons,
		[]Group{Common},
	},
}

func Validate(c Class, g Group) error {
	i, found := slices.BinarySearchFunc(allowedClasses, pair{class: c}, func(a pair, b pair) int {
		return cmp.Compare(a.class, b.class)
	})

	if !found {
		return errors.ErrInvalidAccountClass(c.String())
	}

	if !slices.Contains(allowedClasses[i].allowedGroups, g) {
		return errors.ErrInvalidAccountGroup(c.String(), g.String())
	}

	return nil
}
