package category

import (
	"strings"
	"unicode/utf8"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

const maxNameRunes = 50

type DimensionCategory struct {
	id    uuid.UUID
	sobId uuid.UUID
	name  string
}

func New(id, sobId uuid.UUID, name string) (*DimensionCategory, error) {
	if id == uuid.Nil {
		return nil, commonErrors.NewSlugError("dimension-category-emptyId")
	}

	if sobId == uuid.Nil {
		return nil, commonErrors.NewSlugError("dimension-category-emptySobId")
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, commonErrors.NewSlugError("dimension-category-emptyName")
	}

	if utf8.RuneCountInString(name) > maxNameRunes {
		return nil, commonErrors.NewSlugError("dimension-category-nameTooLong")
	}

	return &DimensionCategory{
		id:    id,
		sobId: sobId,
		name:  name,
	}, nil
}

func (c *DimensionCategory) Rename(newName string) error {
	newName = strings.TrimSpace(newName)

	if newName == "" {
		return commonErrors.NewSlugError("dimension-category-emptyName")
	}

	if utf8.RuneCountInString(newName) > maxNameRunes {
		return commonErrors.NewSlugError("dimension-category-nameTooLong")
	}

	c.name = newName

	return nil
}

func (c *DimensionCategory) Id() uuid.UUID {
	return c.id
}

func (c *DimensionCategory) SobId() uuid.UUID {
	return c.sobId
}

func (c *DimensionCategory) Name() string {
	return c.name
}
