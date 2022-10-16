package sortable

import (
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/common/data/field"
)

type Sort interface {
	Field() field.Field
	Order() string
}

type sortImpl struct {
	field field.Field
	order string
}

// new

func NewSort(fieldName, order string) (Sort, error) {
	f, err := field.New(fieldName)
	if err != nil {
		return nil, err
	}
	if order != "desc" && order != "asc" {
		return nil, errors.Errorf("invalid order %s", order)
	}

	return sortImpl{
		field: f,
		order: order,
	}, nil
}

// impl

func (s sortImpl) Field() field.Field {
	return s.field
}

func (s sortImpl) Order() string {
	return s.order
}
