package amount_type

import "fmt"

type AmountType struct {
	slug string
}

func (c AmountType) String() string {
	return c.slug
}

var (
	unknown             = AmountType{""}
	YearOpeningBalance  = AmountType{"year_opening_balance"}
	PeriodEndingBalance = AmountType{"period_ending_balance"}
	YearToDateAmount    = AmountType{"year_to_date_amount"}
	LastYearAmount      = AmountType{"last_year_amount"}
	PeriodAmount        = AmountType{"period_amount"}
)

var stringToAmountType = map[string]AmountType{
	"year_opening_balance":  YearOpeningBalance,
	"period_ending_balance": PeriodEndingBalance,
	"year_to_date_amount":   YearToDateAmount,
	"last_year_amount":      LastYearAmount,
	"period_amount":         PeriodAmount,
}

func FromString(s string) (AmountType, error) {
	class, ok := stringToAmountType[s]
	if ok {
		return class, nil
	}
	return unknown, fmt.Errorf("unknown amount type %s", s)
}
