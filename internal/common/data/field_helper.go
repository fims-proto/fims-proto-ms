package data

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	matchFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	matchAllCap   = regexp.MustCompile(`([a-z\d])([A-Z])`)
)

func toSnakeCase(field string) (string, error) {
	if match, _ := regexp.Match(`^[a-zA-Z]+$`, []byte(field)); !match {
		return "", errors.Errorf("invalid field name %s", field)
	}

	snake := matchFirstCap.ReplaceAllString(field, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake), nil
}
