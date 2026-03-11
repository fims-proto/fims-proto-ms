package auxiliary_ledger

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AuxiliaryLedger struct {
	id                  uuid.UUID
	sobId               uuid.UUID
	periodId            uuid.UUID
	accountId           uuid.UUID
	auxiliaryCategoryId uuid.UUID
	auxiliaryAccountId  uuid.UUID
	openingAmount       decimal.Decimal
	periodAmount        decimal.Decimal
	periodDebit         decimal.Decimal // positive amount, only for query performance
	periodCredit        decimal.Decimal // positive amount, only for query performance
	endingAmount        decimal.Decimal
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	periodId uuid.UUID,
	accountId uuid.UUID,
	auxiliaryCategoryId uuid.UUID,
	auxiliaryAccountId uuid.UUID,
	openingAmount decimal.Decimal,
	periodAmount decimal.Decimal,
	periodDebit decimal.Decimal,
	periodCredit decimal.Decimal,
	endingAmount decimal.Decimal,
) (*AuxiliaryLedger, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil auxiliary ledger id")
	}

	if sobId == uuid.Nil {
		return nil, errors.New("nil sob id")
	}

	if periodId == uuid.Nil {
		return nil, errors.New("nil period id")
	}

	if accountId == uuid.Nil {
		return nil, errors.New("nil account id")
	}

	if auxiliaryCategoryId == uuid.Nil {
		return nil, errors.New("nil auxiliary category id")
	}

	if auxiliaryAccountId == uuid.Nil {
		return nil, errors.New("nil auxiliary account id")
	}

	return &AuxiliaryLedger{
		id:                  id,
		sobId:               sobId,
		periodId:            periodId,
		accountId:           accountId,
		auxiliaryCategoryId: auxiliaryCategoryId,
		auxiliaryAccountId:  auxiliaryAccountId,
		openingAmount:       openingAmount,
		periodAmount:        periodAmount,
		periodDebit:         periodDebit,
		periodCredit:        periodCredit,
		endingAmount:        endingAmount,
	}, nil
}

func (l *AuxiliaryLedger) Id() uuid.UUID {
	return l.id
}

func (l *AuxiliaryLedger) SobId() uuid.UUID {
	return l.sobId
}

func (l *AuxiliaryLedger) PeriodId() uuid.UUID {
	return l.periodId
}

func (l *AuxiliaryLedger) AccountId() uuid.UUID {
	return l.accountId
}

func (l *AuxiliaryLedger) AuxiliaryCategoryId() uuid.UUID {
	return l.auxiliaryCategoryId
}

func (l *AuxiliaryLedger) AuxiliaryAccountId() uuid.UUID {
	return l.auxiliaryAccountId
}

func (l *AuxiliaryLedger) OpeningAmount() decimal.Decimal {
	return l.openingAmount
}

func (l *AuxiliaryLedger) PeriodAmount() decimal.Decimal {
	return l.periodAmount
}

func (l *AuxiliaryLedger) PeriodDebit() decimal.Decimal {
	return l.periodDebit
}

func (l *AuxiliaryLedger) PeriodCredit() decimal.Decimal {
	return l.periodCredit
}

func (l *AuxiliaryLedger) EndingAmount() decimal.Decimal {
	return l.endingAmount
}
