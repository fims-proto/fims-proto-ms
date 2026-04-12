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
		return nil, commonErrors.NewInternalError(commonErrors.SlugDimOptionEmptyId)
	}

	if categoryId == uuid.Nil {
		return nil, commonErrors.NewInternalError(commonErrors.SlugDimOptionEmptyCategoryId)
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugDimOptionEmptyName)
	}

	if utf8.RuneCountInString(name) > maxNameRunes {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugDimOptionNameTooLong)
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
		return commonErrors.NewInvalidInputError(commonErrors.SlugDimOptionEmptyName)
	}

	if utf8.RuneCountInString(newName) > maxNameRunes {
		return commonErrors.NewInvalidInputError(commonErrors.SlugDimOptionNameTooLong)
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
