package option

import (
	"strings"
	"unicode/utf8"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

const maxNameRunes = 50

type DimensionOption struct {
	id         uuid.UUID
	categoryId uuid.UUID
	name       string
}

func New(id, categoryId uuid.UUID, name string) (*DimensionOption, error) {
	if id == uuid.Nil {
		return nil, commonErrors.NewSlugError("dimension-option-emptyId")
	}

	if categoryId == uuid.Nil {
		return nil, commonErrors.NewSlugError("dimension-option-emptyCategoryId")
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, commonErrors.NewSlugError("dimension-option-emptyName")
	}

	if utf8.RuneCountInString(name) > maxNameRunes {
		return nil, commonErrors.NewSlugError("dimension-option-nameTooLong")
	}

	return &DimensionOption{
		id:         id,
		categoryId: categoryId,
		name:       name,
	}, nil
}

func (o *DimensionOption) Rename(newName string) error {
	newName = strings.TrimSpace(newName)

	if newName == "" {
		return commonErrors.NewSlugError("dimension-option-emptyName")
	}

	if utf8.RuneCountInString(newName) > maxNameRunes {
		return commonErrors.NewSlugError("dimension-option-nameTooLong")
	}

	o.name = newName

	return nil
}

func (o *DimensionOption) Id() uuid.UUID {
	return o.id
}

func (o *DimensionOption) CategoryId() uuid.UUID {
	return o.categoryId
}

func (o *DimensionOption) Name() string {
	return o.name
}
