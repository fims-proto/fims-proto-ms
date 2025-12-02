package auxiliary_category

import (
	"errors"
	"unicode/utf8"

	"github.com/google/uuid"
)

type AuxiliaryCategory struct {
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
) (*AuxiliaryCategory, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil auxiliary category id")
	}

	if sobId == uuid.Nil {
		return nil, errors.New("nil sob id")
	}

	if key == "" {
		return nil, errors.New("nil auxiliary category key")
	}

	if utf8.RuneCountInString(key) > 20 {
		return nil, errors.New("auxiliary category key too long")
	}

	if title == "" {
		return nil, errors.New("nil auxiliary category title")
	}

	if utf8.RuneCountInString(title) > 50 {
		return nil, errors.New("auxiliary category title too long")
	}

	return &AuxiliaryCategory{
		id:         id,
		sobId:      sobId,
		key:        key,
		title:      title,
		isStandard: isStandard,
	}, nil
}

func (a AuxiliaryCategory) Id() uuid.UUID {
	return a.id
}

func (a AuxiliaryCategory) SobId() uuid.UUID {
	return a.sobId
}

func (a AuxiliaryCategory) Key() string {
	return a.key
}

func (a AuxiliaryCategory) Title() string {
	return a.title
}

func (a AuxiliaryCategory) IsStandard() bool {
	return a.isStandard
}
