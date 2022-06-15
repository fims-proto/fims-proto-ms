package domain

import (
	"github.com/google/uuid"
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
		return nil, newDomainErr(errLineItemEmptyId)
	}
	if accountId == uuid.Nil {
		return nil, newDomainErr(errLineItemEmptyAccountId)
	}
	if summary == "" {
		return nil, newDomainErr(errLineItemEmptySummary)
	}
	if debit.IsZero() && credit.IsZero() {
		return nil, newDomainErr(errLineItemEmptyDebitCredit)
	}

	if !debit.IsZero() && !credit.IsZero() {
		return nil, newDomainErr(errLineItemDebitCreditCoExist)
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
