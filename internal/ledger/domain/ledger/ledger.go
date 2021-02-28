package ledger

import (
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"
	"strings"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Ledger struct {
	number         string
	title          string
	superiorNumber string
	accountType    commonAccount.Type
	debit          decimal.Decimal
	credit         decimal.Decimal
	balance        decimal.Decimal
}

func NewLedger(number string, title string, superiorNumber string, accountType commonAccount.Type) (*Ledger, error) {
	if number == "" {
		return nil, errors.New("empty ledger number")
	}
	if title == "" {
		return nil, errors.New("empty ledger title")
	}
	if superiorNumber != "" && !strings.HasPrefix(number, superiorNumber) {
		return nil, errors.New("invalid superior ledger number")
	}

	accType, err := commonAccount.NewAccountType(accountType)
	if err != nil {
		return nil, errors.Wrap(err, "invalid account type")
	}

	return &Ledger{
		number:         number,
		title:          title,
		superiorNumber: superiorNumber,
		accountType:    accType,
		debit:          decimal.RequireFromString("0"),
		credit:         decimal.RequireFromString("0"),
		balance:        decimal.RequireFromString("0"),
	}, nil
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

func (l Ledger) AccountType() commonAccount.Type {
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
