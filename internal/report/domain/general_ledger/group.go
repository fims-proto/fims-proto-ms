package general_ledger

import "fmt"

type Group int

const (
	CurrentAsset              Group = 101 // 流动资产
	NonCurrentAsset           Group = 102 // 非流动资产
	CurrentLiability          Group = 201 // 流动负债
	NonCurrentLiability       Group = 202 // 非流动负债
	Equity                    Group = 301 // 所有者权益
	Cost                      Group = 401 // 成本
	OperatingIncome           Group = 501 // 营业收入
	OtherIncome               Group = 502 // 其他收益
	PeriodCost                Group = 503 // 期间费用
	OtherCost                 Group = 504 // 其他损失
	OperatingCostAndTax       Group = 505 // 营业成本及税金
	PriorYearIncomeAdjustment Group = 506 // 以前年度损益调整
	IncomeTax                 Group = 507 // 所得税
	Common                    Group = 701 // 共同
)

var groupToString = map[Group]string{
	CurrentAsset:              "current asset",
	NonCurrentAsset:           "non current asset",
	CurrentLiability:          "current liability",
	NonCurrentLiability:       "non current liability",
	Equity:                    "equity",
	Cost:                      "cost",
	OperatingIncome:           "operating income",
	OtherIncome:               "other income",
	PeriodCost:                "period cost",
	OtherCost:                 "other cost",
	OperatingCostAndTax:       "operating cost and tax",
	PriorYearIncomeAdjustment: "prior year income adjustment",
	IncomeTax:                 "income tax",
	Common:                    "common",
}

func (g Group) String() string {
	s, ok := groupToString[g]
	if ok {
		return s
	}
	return fmt.Sprintf("unknown group [%d]", g)
}
