package data

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type Operator int

const (
	OptEq         Operator = 1 << iota // equal
	OptBt                              // between
	OptLt                              // less than
	OptLte                             // less than equal
	OptGt                              // greater than
	OptGte                             // greater than equal
	OptIn                              // in
	OptStartsWith                      // starts with
)

func (o Operator) String() string {
	switch o {
	case OptEq:
		return "="
	case OptBt:
		return "BETWEEN"
	case OptLt:
		return "<"
	case OptLte:
		return "<="
	case OptGt:
		return ">"
	case OptGte:
		return ">="
	case OptIn:
		return "IN"
	case OptStartsWith:
		return "startsWith"
	default:
		return "unknown"
	}
}

func newOperator(o string) (Operator, error) {
	switch o {
	case "eq":
		return OptEq, nil
	case "bt":
		return OptBt, nil
	case "lt":
		return OptLt, nil
	case "lte":
		return OptLte, nil
	case "gt":
		return OptGt, nil
	case "gte":
		return OptGte, nil
	case "in":
		return OptIn, nil
	case "startsWith":
		return OptStartsWith, nil
	default:
		return 0, errors.Errorf("operator %s not supported", o)
	}
}

type Filter interface {
	Field() string
	Operator() Operator
	Values() []string
}

type filterImpl struct {
	field    string
	operator Operator
	values   []string
}

func newFilterImpl(field, operator string, values []string) (Filter, error) {
	snakeCase, err := toSnakeCase(field)
	if err != nil {
		return nil, err
	}

	o, err := newOperator(operator)
	if err != nil {
		return nil, err
	}

	switch o {
	case OptBt:
		if len(values) != 2 {
			return nil, errors.Errorf("invalid values for operator %s", o)
		}
	default:
		if len(values) == 0 {
			return nil, errors.Errorf("invalid values for operator %s", o)
		}
	}

	return filterImpl{
		field:    snakeCase,
		operator: o,
		values:   removeSingleQuote(values),
	}, nil
}

func removeSingleQuote(values []string) []string {
	singleQuoted := regexp.MustCompile(`^'(.+)'$`)
	for i := 0; i < len(values); i++ {
		if singleQuoted.MatchString(values[i]) {
			values[i] = values[i][1 : len(values[i])-1]
		}
	}
	return values
}

func (f filterImpl) Field() string {
	return f.field
}

func (f filterImpl) Operator() Operator {
	return f.operator
}

func (f filterImpl) Values() []string {
	return f.values
}

func newFiltersFromQuery(filter string) ([]Filter, error) {
	if filter == "" {
		return nil, nil
	}

	andSep := regexp.MustCompile(`\band\b|\bAND\b`)
	spaceSep := regexp.MustCompile(`\s+`)
	commaSep := regexp.MustCompile(`\s*,\s*`)

	conditions := andSep.Split(filter, -1)

	var result []Filter

	for _, raw := range conditions {
		condition := strings.TrimSpace(raw)
		components := spaceSep.Split(condition, -1)
		if len(components) < 3 {
			return nil, errors.Errorf("invalid filter %s", condition)
		}
		field := components[0]
		opt := components[1]
		valuesString := strings.Join(components[2:], " ")
		values := commaSep.Split(valuesString, -1)

		filterImpl, err := newFilterImpl(field, opt, values)
		if err != nil {
			return nil, err
		}
		result = append(result, filterImpl)
	}

	return result, nil
}
