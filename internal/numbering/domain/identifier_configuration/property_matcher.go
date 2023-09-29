package identifier_configuration

import (
	"errors"
	"strings"
)

type PropertyMatcher struct {
	name  string
	value string
}

func NewPropertyMatcher(name, value string) (*PropertyMatcher, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("property name cannot be empty")
	}

	if strings.TrimSpace(value) == "" {
		return nil, errors.New("property value cannot be empty")
	}

	return &PropertyMatcher{
		name:  strings.TrimSpace(name),
		value: strings.TrimSpace(value),
	}, nil
}

func (m PropertyMatcher) Name() string {
	return m.name
}

func (m PropertyMatcher) Value() string {
	return m.value
}
