package filterable

import (
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/field"
)

type Filter interface {
	Field() field.Field
	Operator() Operator
	Values() []any
}

type filterImpl struct {
	field    field.Field
	operator Operator
	values   []any
}

// new

func NewFilter[T any](fieldName, operator string, values ...T) (Filter, error) {
	f, err := field.New(fieldName)
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
			return nil, fmt.Errorf("invalid values for operator %s", o)
		}
	default:
		if len(values) == 0 {
			return nil, fmt.Errorf("invalid values for operator %s", o)
		}
	}

	sliceAny := make([]any, len(values))
	for i, v := range values {
		sliceAny[i] = v
	}

	return filterImpl{
		field:    f,
		operator: o,
		values:   sliceAny,
	}, nil
}

// impl

func (f filterImpl) Field() field.Field {
	return f.field
}

func (f filterImpl) Operator() Operator {
	return f.operator
}

func (f filterImpl) Values() []any {
	return f.values
}

// misc

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
		return 0, fmt.Errorf("operator %s not supported", o)
	}
}
