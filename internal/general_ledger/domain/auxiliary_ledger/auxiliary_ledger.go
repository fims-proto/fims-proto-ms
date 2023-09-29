package auxiliary_ledger

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
)

type AuxiliaryLedger struct {
	id               uuid.UUID
	periodId         uuid.UUID
	auxiliaryAccount *auxiliary_account.AuxiliaryAccount
	openingBalance   decimal.Decimal
	endingBalance    decimal.Decimal
	periodDebit      decimal.Decimal
	periodCredit     decimal.Decimal
}

func New(
	id uuid.UUID,
	periodId uuid.UUID,
	auxiliaryAccount *auxiliary_account.AuxiliaryAccount,
	openingBalance decimal.Decimal,
	endingBalance decimal.Decimal,
	periodDebit decimal.Decimal,
	periodCredit decimal.Decimal,
) (*AuxiliaryLedger, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil auxiliary ledger id")
	}

	if periodId == uuid.Nil {
		return nil, errors.New("nil period id")
	}

	if auxiliaryAccount == nil {
		return nil, errors.New("nil auxiliary account")
	}

	return &AuxiliaryLedger{
		id:               id,
		periodId:         periodId,
		auxiliaryAccount: auxiliaryAccount,
		openingBalance:   openingBalance,
		endingBalance:    endingBalance,
		periodDebit:      periodDebit,
		periodCredit:     periodCredit,
	}, nil
}

func (l *AuxiliaryLedger) Id() uuid.UUID {
	return l.id
}

func (l *AuxiliaryLedger) PeriodId() uuid.UUID {
	return l.periodId
}

func (l *AuxiliaryLedger) AuxiliaryAccount() *auxiliary_account.AuxiliaryAccount {
	return l.auxiliaryAccount
}

func (l *AuxiliaryLedger) OpeningBalance() decimal.Decimal {
	return l.openingBalance
}

func (l *AuxiliaryLedger) EndingBalance() decimal.Decimal {
	return l.endingBalance
}

func (l *AuxiliaryLedger) PeriodDebit() decimal.Decimal {
	return l.periodDebit
}

func (l *AuxiliaryLedger) PeriodCredit() decimal.Decimal {
	return l.periodCredit
}
