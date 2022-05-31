package data

import (
	"strings"

	"github.com/pkg/errors"
)

type Choose interface {
	Field() string
}

type chooseRequest struct {
	field string
}

func newChoose(field string) (Choose, error) {
	fieldSnakeCase, err := toSnakeCase(field)
	if err != nil {
		return nil, err
	}

	return chooseRequest{
		field: fieldSnakeCase,
	}, nil
}

func (c chooseRequest) Field() string {
	return c.field
}

func newChoosesFromQuery(choose string) ([]Choose, error) {
	if choose == "" {
		return nil, nil
	}
	var chooses []Choose
	chooseSegments := strings.Split(choose, ",")
	for _, segment := range chooseSegments {
		choose, err := newChoose(segment)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create chooses request")
		}
		chooses = append(chooses, choose)
	}
	return chooses, nil
}
