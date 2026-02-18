package class

import "fmt"

type Class struct {
	slug string
}

func (c Class) String() string {
	return c.slug
}

var (
	unknown         = Class{""}
	BalanceSheet    = Class{"balance_sheet"}
	IncomeStatement = Class{"income_statement"}
)

var stringToClass = map[string]Class{
	"balance_sheet":    BalanceSheet,
	"income_statement": IncomeStatement,
}

func FromString(s string) (Class, error) {
	class, ok := stringToClass[s]
	if ok {
		return class, nil
	}
	return unknown, fmt.Errorf("unknown report class %s", s)
}
