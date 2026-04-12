package identifier_configuration

import (
	"strings"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
)

type PropertyMatcher struct {
	name  string
	value string
}

func NewPropertyMatcher(name, value string) (*PropertyMatcher, error) {
	if strings.TrimSpace(name) == "" {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugNumberingPropertyNameEmpty)
	}

	if strings.TrimSpace(value) == "" {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugNumberingPropertyValueEmpty)
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
