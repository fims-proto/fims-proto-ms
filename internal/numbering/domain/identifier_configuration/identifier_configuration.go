package identifier_configuration

import (
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

type IdentifierConfiguration struct {
	id                   uuid.UUID
	targetBusinessObject string
	propertyMatchers     []PropertyMatcher
	counter              int
	prefix               string
	suffix               string
}

func New(id uuid.UUID, targetBusinessObject string, propertyMatchers []PropertyMatcher, counter int, prefix, suffix string) (*IdentifierConfiguration, error) {
	if id == uuid.Nil {
		return nil, commonErrors.NewInternalError(commonErrors.SlugNumberingIdEmpty)
	}

	if targetBusinessObject == "" {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugNumberingTargetObjectEmpty)
	}

	if len(propertyMatchers) == 0 {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugNumberingPropertyMatchersEmpty)
	}

	if counter < 0 {
		return nil, fmt.Errorf("unexpected configuration counter %d", counter)
	}

	return &IdentifierConfiguration{
		id:                   id,
		targetBusinessObject: targetBusinessObject,
		propertyMatchers:     propertyMatchers,
		counter:              counter,
		prefix:               prefix,
		suffix:               suffix,
	}, nil
}

func (c *IdentifierConfiguration) Id() uuid.UUID {
	return c.id
}

func (c *IdentifierConfiguration) TargetBusinessObject() string {
	return c.targetBusinessObject
}

func (c *IdentifierConfiguration) PropertyMatchers() []PropertyMatcher {
	return c.propertyMatchers
}

func (c *IdentifierConfiguration) Counter() int {
	return c.counter
}

func (c *IdentifierConfiguration) Prefix() string {
	return c.prefix
}

func (c *IdentifierConfiguration) Suffix() string {
	return c.suffix
}

func (c *IdentifierConfiguration) Stringify() string {
	return fmt.Sprintf("%s%d%s", c.prefix, c.counter, c.suffix)
}
