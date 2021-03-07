package domain

import (
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"strings"

	"github.com/pkg/errors"
)

type Account struct {
	number         string
	title          string
	superiorNumber string
	accountType    commonAccount.Type
}

func NewAccount(number string, title string, superiorNumber string, accType commonAccount.Type) (*Account, error) {
	if number == "" {
		return nil, errors.New("empty account number")
	}
	if title == "" {
		return nil, errors.New("empty account title")
	}
	if superiorNumber != "" && !strings.HasPrefix(number, superiorNumber) {
		return nil, errors.New("invalid superior account number")
	}

	return &Account{
		number:         number,
		title:          title,
		superiorNumber: superiorNumber,
		accountType:    accType,
	}, nil
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

func (acc Account) Type() commonAccount.Type {
	return acc.accountType
}
