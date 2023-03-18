package filterable

import (
	"github/fims-proto/fims-proto-ms/internal/common/data/field"

	"github.com/pkg/errors"
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

func NewFilter[T any](fieldName string, operator Operator, values ...T) (Filter, error) {
	f, err := field.New(fieldName)
	if err != nil {
		return nil, err
	}

	o := operator

	switch o {
	case OptBtw:
		if len(values) != 2 {
			return nil, errors.Errorf("invalid values for operator %s", o)
		}
	default:
		if len(values) == 0 {
			return nil, errors.Errorf("invalid values for operator %s", o)
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

// impl for Filter

func (f filterImpl) Field() field.Field {
	return f.field
}

func (f filterImpl) Operator() Operator {
	return f.operator
}

func (f filterImpl) Values() []any {
	return f.values
}

// impl for Filterable
func (f filterImpl) Children() []Filterable {
	return nil
}

func (f filterImpl) IsFiltered() bool {
	return true
}

func (f filterImpl) FilterableType() FilterableType {
	return TypeATOM
}

// misc

type Operator int

const (
	OptBtw Operator = 1 << iota // between
	OptCtn                      // contain
	OptEq                       // equal
	OptGt                       // greater than
	OptGte                      // greater than equal
	OptIn                       // in
	OptLt                       // less than
	OptLte                      // less than equal
	OptStw                      // starts with
)

func (o Operator) String() string {
	switch o {
	case OptEq:
		return "="
	case OptBtw:
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
	case OptStw:
		return "startsWith"
	case OptCtn:
		return "contain"
	default:
		return "unknown"
	}
}
