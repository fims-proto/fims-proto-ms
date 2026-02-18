package formula_rule

import "fmt"

type FormulaRule struct {
	slug string
}

func (r FormulaRule) String() string {
	return r.slug
}

var (
	unknown     = FormulaRule{""}
	Net         = FormulaRule{"net"}
	Debit       = FormulaRule{"debit"}
	Credit      = FormulaRule{"credit"}
	Transaction = FormulaRule{"transaction"}
)

var stringToFormulaRule = map[string]FormulaRule{
	"net":         Net,
	"debit":       Debit,
	"credit":      Credit,
	"transaction": Transaction,
}

func FromString(s string) (FormulaRule, error) {
	rule, ok := stringToFormulaRule[s]
	if ok {
		return rule, nil
	}
	return unknown, fmt.Errorf("unknown formula rule: %s", s)
}
