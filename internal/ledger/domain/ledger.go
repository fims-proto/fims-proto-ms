package domain

import (
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Ledger struct {
	id             uuid.UUID
	periodId       uuid.UUID
	accountId      uuid.UUID
	openingBalance decimal.Decimal
	endingBalance  decimal.Decimal
	debit          decimal.Decimal
	credit         decimal.Decimal
}

func NewLedger(id, periodId, accountId uuid.UUID, openingBalance, endingBalance, debit, credit decimal.Decimal) (*Ledger, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil ledger id")
	}
	if periodId == uuid.Nil {
		return nil, errors.New("nil period id")
	}
	if accountId == uuid.Nil {
		return nil, errors.New("nil account id")
	}

	return &Ledger{
		id:             id,
		periodId:       periodId,
		accountId:      accountId,
		openingBalance: openingBalance,
		endingBalance:  endingBalance,
		debit:          debit,
		credit:         credit,
	}, nil
}

func (l Ledger) Id() uuid.UUID {
	return l.id
}

func (l Ledger) PeriodId() uuid.UUID {
	return l.periodId
}

func (l Ledger) AccountId() uuid.UUID {
	return l.accountId
}

func (l Ledger) OpeningBalance() decimal.Decimal {
	return l.openingBalance
}

func (l Ledger) EndingBalance() decimal.Decimal {
	return l.endingBalance
}

func (l Ledger) Debit() decimal.Decimal {
	return l.debit
}

func (l Ledger) Credit() decimal.Decimal {
	return l.credit
}

func (l *Ledger) UpdatePeriodAmount(debit, credit decimal.Decimal, balanceDirection commonAccount.Direction) {
	l.debit = debit
	l.credit = credit

	addingAmount := l.debit
	subAmount := l.credit

	if balanceDirection == commonAccount.Credit {
		addingAmount = l.credit
		subAmount = l.debit
	}

	l.endingBalance = l.openingBalance.Add(addingAmount).Sub(subAmount)
}
