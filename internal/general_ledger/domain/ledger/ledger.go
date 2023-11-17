package ledger

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account"
)

type Ledger struct {
	id                   uuid.UUID
	sobId                uuid.UUID
	periodId             uuid.UUID
	accountId            uuid.UUID
	account              *account.Account
	openingDebitBalance  decimal.Decimal
	openingCreditBalance decimal.Decimal
	periodDebit          decimal.Decimal
	periodCredit         decimal.Decimal
	endingDebitBalance   decimal.Decimal
	endingCreditBalance  decimal.Decimal
}

func New(
	id uuid.UUID,
	sobId uuid.UUID,
	periodId uuid.UUID,
	accountId uuid.UUID,
	account *account.Account,
	openingDebitBalance decimal.Decimal,
	openingCreditBalance decimal.Decimal,
	periodDebit decimal.Decimal,
	periodCredit decimal.Decimal,
	endingDebitBalance decimal.Decimal,
	endingCreditBalance decimal.Decimal,
) (*Ledger, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil ledger id")
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

	if account == nil {
		return nil, errors.New("nil account")
	}

	return &Ledger{
		id:                   id,
		sobId:                sobId,
		accountId:            accountId,
		periodId:             periodId,
		openingDebitBalance:  openingDebitBalance,
		openingCreditBalance: openingCreditBalance,
		periodDebit:          periodDebit,
		periodCredit:         periodCredit,
		endingDebitBalance:   endingDebitBalance,
		endingCreditBalance:  endingCreditBalance,
		account:              account,
	}, nil
}

func (l *Ledger) Id() uuid.UUID {
	return l.id
}

func (l *Ledger) SobId() uuid.UUID {
	return l.sobId
}

func (l *Ledger) PeriodId() uuid.UUID {
	return l.periodId
}

func (l *Ledger) AccountId() uuid.UUID {
	return l.accountId
}

func (l *Ledger) Account() *account.Account {
	return l.account
}

func (l *Ledger) OpeningDebitBalance() decimal.Decimal {
	return l.openingDebitBalance
}

func (l *Ledger) OpeningCreditBalance() decimal.Decimal {
	return l.openingCreditBalance
}

func (l *Ledger) PeriodDebit() decimal.Decimal {
	return l.periodDebit
}

func (l *Ledger) PeriodCredit() decimal.Decimal {
	return l.periodCredit
}

func (l *Ledger) EndingDebitBalance() decimal.Decimal {
	return l.endingDebitBalance
}

func (l *Ledger) EndingCreditBalance() decimal.Decimal {
	return l.endingCreditBalance
}
