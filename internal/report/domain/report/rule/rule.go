package rule

import "fmt"

type Rule struct {
	slug string
}

func (r Rule) String() string {
	return r.slug
}

var (
	Unknown       = Rule{""}
	Balance       = Rule{"B"}
	DebitBalance  = Rule{"DB"}
	CreditBalance = Rule{"CB"}

	ProfitAndLoss = Rule{"PL"}
	DebitAmount   = Rule{"DA"}
	CreditAmount  = Rule{"CA"}
)

var stringToRule = map[string]Rule{
	"B":  Balance,
	"DB": DebitBalance,
	"CB": CreditBalance,
	"PL": ProfitAndLoss,
	"DA": DebitAmount,
	"CA": CreditAmount,
}

func FromString(s string) (Rule, error) {
	rule, ok := stringToRule[s]
	if ok {
		return rule, nil
	}
	return Unknown, fmt.Errorf("unknown rule: %s", s)
}
