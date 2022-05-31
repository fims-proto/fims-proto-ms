package data

import (
	"strings"

	"github.com/pkg/errors"
)

type Sort interface {
	Field() string
	Order() string
}

type sortRequest struct {
	field string
	order string
}

func newSort(field, order string) (Sort, error) {
	fieldSnakeCase, err := toSnakeCase(field)
	if err != nil {
		return nil, err
	}
	if order != "desc" && order != "asc" {
		return nil, errors.Errorf("invalid order %s", order)
	}

	return sortRequest{
		field: fieldSnakeCase,
		order: order,
	}, nil
}

func (s sortRequest) Field() string {
	return s.field
}

func (s sortRequest) Order() string {
	return s.order
}

func newSortsFromQuery(sort string) ([]Sort, error) {
	if sort == "" {
		return nil, nil
	}

	sortFields := make(map[string]string)
	sortSegments := strings.Split(sort, ",")
	for _, segment := range sortSegments {
		elements := strings.Split(strings.TrimSpace(segment), " ")
		if len(elements) == 1 {
			sortFields[elements[0]] = "asc"
		} else if len(elements) == 2 {
			sortFields[elements[0]] = elements[1]
		} else {
			return nil, errors.Errorf("invalid sort query parameter %s", sort)
		}
	}

	var sorts []Sort
	for field, order := range sortFields {
		sort, err := newSort(field, order)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create sorts request")
		}
		sorts = append(sorts, sort)
	}

	return sorts, nil
}
