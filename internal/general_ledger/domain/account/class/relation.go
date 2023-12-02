package class

type pair struct {
	Class  Class
	Groups []Group
}

// Classes must be sorted by Class increasingly
var Classes = []pair{
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
