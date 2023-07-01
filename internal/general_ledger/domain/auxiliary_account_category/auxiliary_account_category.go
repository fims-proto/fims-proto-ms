package auxiliary_account_category

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type AuxiliaryAccountCategory struct {
	id         uuid.UUID
	sobId      uuid.UUID
	key        string
	title      string
	isStandard bool
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	key string,
	title string,
	isStandard bool,
) (*AuxiliaryAccountCategory, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil auxiliary account category id")
	}

	if sobId == uuid.Nil {
		return nil, errors.New("nil sob id")
	}

	if key == "" {
		return nil, errors.New("nil auxiliary account category key")
	}

	if title == "" {
		return nil, errors.New("nil auxiliary account category title")
	}

	return &AuxiliaryAccountCategory{
		id:         id,
		sobId:      sobId,
		key:        key,
		title:      title,
		isStandard: isStandard,
	}, nil
}

func (a AuxiliaryAccountCategory) Id() uuid.UUID {
	return a.id
}

func (a AuxiliaryAccountCategory) SobId() uuid.UUID {
	return a.sobId
}

func (a AuxiliaryAccountCategory) Key() string {
	return a.key
}

func (a AuxiliaryAccountCategory) Title() string {
	return a.title
}

func (a AuxiliaryAccountCategory) IsStandard() bool {
	return a.isStandard
}
