package sortable

import (
	"strings"

	"github.com/pkg/errors"
)

func NewSortableFromQuery(sort string) (Sortable, error) {
	if sort == "" {
		return Unsorted(), nil
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
		sort, err := NewSort(field, order)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create sorts request")
		}
		sorts = append(sorts, sort)
	}

	return New(sorts), nil
}
