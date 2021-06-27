package domain

import (
	commonaccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Ledger struct {
	sob            string
	number         string
	title          string
	superiorNumber string
	accountType    commonaccount.Type
	debit          decimal.Decimal
	credit         decimal.Decimal
	balance        decimal.Decimal
}

func NewLedger(sob, number, title, superiorNumber string, accountType commonaccount.Type) (*Ledger, error) {
	if sob == "" {
		return nil, errors.New("empty sob")
	}
	if number == "" {
		return nil, errors.New("empty ledger number")
	}
	if title == "" {
		return nil, errors.New("empty ledger title")
	}
	if superiorNumber != "" && !strings.HasPrefix(number, superiorNumber) {
		return nil, errors.New("invalid superior ledger number")
	}

	accType, err := commonaccount.NewAccountType(accountType)
	if err != nil {
		return nil, errors.Wrap(err, "invalid account type")
	}

	return &Ledger{
		sob:            sob,
		number:         number,
		title:          title,
		superiorNumber: superiorNumber,
		accountType:    accType,
		debit:          decimal.RequireFromString("0"),
		credit:         decimal.RequireFromString("0"),
		balance:        decimal.RequireFromString("0"),
	}, nil
}

func (l Ledger) Sob() string {
	return l.sob
}

func (l Ledger) Number() string {
	return l.number
}

func (l Ledger) Title() string {
	return l.title
}

func (l Ledger) SuperiorNumber() string {
	return l.superiorNumber
}

func (l Ledger) AccountType() commonaccount.Type {
	return l.accountType
}

func (l Ledger) Debit() decimal.Decimal {
	return l.debit
}

func (l Ledger) Credit() decimal.Decimal {
	return l.credit
}

func (l Ledger) Balance() decimal.Decimal {
	return l.balance
}
