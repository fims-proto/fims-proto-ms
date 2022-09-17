package field

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	matchAllCap   = regexp.MustCompile(`([a-z\d])([A-Z])`)
)

func ToColumn(f Field, resolveEntity func(entity string) (string, error)) (string, error) {
	snake := matchFirstCap.ReplaceAllString(f.Name(), "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	fieldName := strings.ToLower(snake)

	if f.Entity() == "" {
		return fieldName, nil
	}

	entity, err := resolveEntity(f.Entity())
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.%s", entity, fieldName), nil
}
