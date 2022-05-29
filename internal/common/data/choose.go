package data

import (
	"regexp"
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
	if match, _ := regexp.Match(`^[a-zA-Z]+$`, []byte(field)); !match {
		return nil, errors.Errorf("invalid field name %s", field)
	}

	return chooseRequest{
		field: toSnakeCase(field),
	}, nil
}

func (c chooseRequest) Field() string {
	return c.field
}

var (
	matchFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	matchAllCap   = regexp.MustCompile(`([a-z\d])([A-Z])`)
)

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
