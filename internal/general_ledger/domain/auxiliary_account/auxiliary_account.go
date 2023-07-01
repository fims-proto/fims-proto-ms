package auxiliary_account

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account_category"
)

type AuxiliaryAccount struct {
	id          uuid.UUID
	category    *auxiliary_account_category.AuxiliaryAccountCategory
	key         string
	title       string
	description string
}

func New(
	id uuid.UUID,
	category *auxiliary_account_category.AuxiliaryAccountCategory,
	key string,
	title string,
	description string,
) (*AuxiliaryAccount, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil auxiliary account id")
	}

	if category == nil {
		return nil, errors.New("nil auxiliary account category")
	}

	if key == "" {
		return nil, errors.New("nil auxiliary account key")
	}

	if title == "" {
		return nil, errors.New("nil auxiliary account title")
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

func (a AuxiliaryAccount) Category() *auxiliary_account_category.AuxiliaryAccountCategory {
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
