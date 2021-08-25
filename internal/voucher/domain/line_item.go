package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type LineItem struct {
	id            uuid.UUID
	summary       string
	accountNumber string
	debit         decimal.Decimal
	credit        decimal.Decimal
}

func NewLineItem(id uuid.UUID, summary string, accountNumber string, debit string, credit string) (*LineItem, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil id")
	}
	if summary == "" {
		return nil, errors.New("empty summary")
	}
	if accountNumber == "" {
		return nil, errors.New("empty account number")
	}
	if debit == "" && credit == "" {
		return nil, errors.New("empty debit and credit amount")
	}
	if debit != "" && credit != "" {
		return nil, errors.New("both debit and credit amount provided")
	}
	debitDecimal, err := decimal.NewFromString(debit)
	if debit != "" && err != nil {
		return nil, errors.New("invalid debit amount")
	}
	creditDecimal, err := decimal.NewFromString(credit)
	if credit != "" && err != nil {
		return nil, errors.New("invalid credit amount")
	}

	return &LineItem{
		id:            id,
		summary:       summary,
		accountNumber: accountNumber,
		debit:         debitDecimal,
		credit:        creditDecimal,
	}, nil
}

func (l LineItem) Id() uuid.UUID {
	return l.id
}

func (l LineItem) Summary() string {
	return l.summary
}

func (l LineItem) AccountNumber() string {
	return l.accountNumber
}

func (l LineItem) Debit() decimal.Decimal {
	return l.debit
}

func (l LineItem) Credit() decimal.Decimal {
	return l.credit
}
