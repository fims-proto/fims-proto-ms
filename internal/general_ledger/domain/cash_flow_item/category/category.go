package category

import "fmt"

type Category struct {
	slug string
}

func (c Category) String() string {
	return c.slug
}

var (
	unknown   = Category{""}
	Operating = Category{"OPERATING"}
	Investing = Category{"INVESTING"}
	Financing = Category{"FINANCING"}
)

func FromString(s string) (Category, error) {
	switch s {
	case Operating.slug:
		return Operating, nil
	case Investing.slug:
		return Investing, nil
	case Financing.slug:
		return Financing, nil
	}
	return unknown, fmt.Errorf("unknown cash flow category: %s", s)
}
