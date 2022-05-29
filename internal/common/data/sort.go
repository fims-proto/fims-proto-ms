package data

import "github.com/pkg/errors"

type Sort interface {
	Field() string
	Order() string
}

type sortRequest struct {
	field string
	order string
}

func newSort(field, order string) (Sort, error) {
	if field == "" {
		return nil, errors.New("empty field")
	}
	if order != "desc" && order != "asc" {
		return nil, errors.Errorf("invalid order %s", order)
	}

	return sortRequest{
		field: field,
		order: order,
	}, nil
}

func (s sortRequest) Field() string {
	return s.field
}

func (s sortRequest) Order() string {
	return s.order
}
