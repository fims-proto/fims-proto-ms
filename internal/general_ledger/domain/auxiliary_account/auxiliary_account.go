package auxiliary_account

import (
	"errors"
	"unicode/utf8"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"

	"github.com/google/uuid"
)

type AuxiliaryAccount struct {
	id          uuid.UUID
	category    *auxiliary_category.AuxiliaryCategory
	key         string
	title       string
	description string
}

func New(
	id uuid.UUID,
	category *auxiliary_category.AuxiliaryCategory,
	key string,
	title string,
	description string,
) (*AuxiliaryAccount, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil auxiliary account id")
	}

	if category == nil {
		return nil, errors.New("nil auxiliary category")
	}

	if key == "" {
		return nil, errors.New("nil auxiliary account key")
	}

	if utf8.RuneCountInString(key) > 20 {
		return nil, errors.New("auxiliary account key too long")
	}

	if title == "" {
		return nil, errors.New("nil auxiliary account title")
	}

	if utf8.RuneCountInString(title) > 50 {
		return nil, errors.New("auxiliary account title too long")
	}

	if utf8.RuneCountInString(description) > 500 {
		return nil, errors.New("auxiliary account description too long")
	}

	return &AuxiliaryAccount{
		id:          id,
		category:    category,
		key:         key,
		title:       title,
		description: description,
	}, nil
}

func (a AuxiliaryAccount) Id() uuid.UUID {
	return a.id
}

func (a AuxiliaryAccount) Category() *auxiliary_category.AuxiliaryCategory {
	return a.category
}

func (a AuxiliaryAccount) Key() string {
	return a.key
}

func (a AuxiliaryAccount) Title() string {
	return a.title
}

func (a AuxiliaryAccount) Description() string {
	return a.description
}

type AuxiliaryPair struct {
	CategoryKey string
	AccountKey  string
}
