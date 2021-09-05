package domain

import (
	commonaccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Account struct {
	id             uuid.UUID
	sob            string
	number         string
	title          string
	superiorNumber string
	accountType    commonaccount.Type
}

func NewAccount(id uuid.UUID, sob, number, title, superiorNumber, accType string) (*Account, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil uuid")
	}
	if sob == "" {
		return nil, errors.New("empty sob")
	}
	if number == "" {
		return nil, errors.New("empty account number")
	}
	if title == "" {
		return nil, errors.New("empty account title")
	}
	if superiorNumber != "" && !strings.HasPrefix(number, superiorNumber) {
		return nil, errors.New("invalid superior account number")
	}
	at, err := commonaccount.NewAccountTypeFromString(accType)
	if err != nil {
		return nil, errors.Wrap(err, "invalid account type")
	}

	return &Account{
		id:             id,
		sob:            sob,
		number:         number,
		title:          title,
		superiorNumber: superiorNumber,
		accountType:    at,
	}, nil
}

func (acc Account) Id() uuid.UUID {
	return acc.id
}

func (acc Account) Sob() string {
	return acc.sob
}

func (acc Account) Number() string {
	return acc.number
}

func (acc Account) Title() string {
	return acc.title
}

func (acc Account) SuperiorNumber() string {
	return acc.superiorNumber
}

func (acc Account) Type() commonaccount.Type {
	return acc.accountType
}
