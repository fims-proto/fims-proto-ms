package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type LineItem struct {
	id        uuid.UUID
	accountId uuid.UUID
	summary   string
	debit     decimal.Decimal
	credit    decimal.Decimal
}

func NewLineItem(id, accountId uuid.UUID, summary string, debit, credit decimal.Decimal) (*LineItem, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil id")
	}
	if accountId == uuid.Nil {
		return nil, errors.New("nil account id")
	}
	if summary == "" {
		return nil, errors.New("empty summary")
	}
	if debit.IsZero() && credit.IsZero() {
		return nil, errors.New("credit and debit cannot both be zero")
	}

	if !debit.IsZero() && !credit.IsZero() {
		return nil, errors.New("credit and debit cannot both be non zero")
	}

	return &LineItem{
		id:        id,
		summary:   summary,
		accountId: accountId,
		debit:     debit,
		credit:    credit,
	}, nil
}

func (l LineItem) Id() uuid.UUID {
	return l.id
}

func (l LineItem) Summary() string {
	return l.summary
}

func (l LineItem) AccountId() uuid.UUID {
	return l.accountId
}

func (l LineItem) Debit() decimal.Decimal {
	return l.debit
}

func (l LineItem) Credit() decimal.Decimal {
	return l.credit
}
