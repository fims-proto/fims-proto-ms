package field

import (
	"fmt"
	"regexp"
	"strings"
)

type Field interface {
	Entity() string
	Name() string
	Equals(f Field) bool
}

type fieldImpl struct {
	entity string
	name   string
}

func (f fieldImpl) Entity() string {
	return f.entity
}

func (f fieldImpl) Name() string {
	return f.name
}

func (f fieldImpl) Equals(another Field) bool {
	return f.entity == another.Entity() && f.name == another.Name()
}

func New(fieldName string) (Field, error) {
	// can only contain camel cased field name with . as separator
	if match, _ := regexp.Match(`^[a-zA-Z.]+$`, []byte(fieldName)); !match {
		return nil, fmt.Errorf("invalid field name %s", fieldName)
	}

	parts := strings.Split(fieldName, ".")

	switch len(parts) {
	case 1:
		return fieldImpl{
			entity: "",
			name:   parts[0],
		}, nil
	case 2:
		return fieldImpl{
			entity: parts[0],
			name:   parts[1],
		}, nil
	default:
		return nil, fmt.Errorf("invalid field name %s", fieldName)
	}
}
