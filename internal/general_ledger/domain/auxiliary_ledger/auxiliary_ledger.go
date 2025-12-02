package auxiliary_ledger

import (
	"errors"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AuxiliaryLedger struct {
	id                   uuid.UUID
	periodId             uuid.UUID
	auxiliaryAccount     *auxiliary_account.AuxiliaryAccount
	openingDebitBalance  decimal.Decimal
	openingCreditBalance decimal.Decimal
	periodDebit          decimal.Decimal
	periodCredit         decimal.Decimal
	endingDebitBalance   decimal.Decimal
	endingCreditBalance  decimal.Decimal
}

func New(
	id uuid.UUID,
	periodId uuid.UUID,
	auxiliaryAccount *auxiliary_account.AuxiliaryAccount,
	openingDebitBalance decimal.Decimal,
	openingCreditBalance decimal.Decimal,
	periodDebit decimal.Decimal,
	periodCredit decimal.Decimal,
	endingDebitBalance decimal.Decimal,
	endingCreditBalance decimal.Decimal,
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
		id:                   id,
		periodId:             periodId,
		auxiliaryAccount:     auxiliaryAccount,
		openingDebitBalance:  openingDebitBalance,
		openingCreditBalance: openingCreditBalance,
		periodDebit:          periodDebit,
		periodCredit:         periodCredit,
		endingDebitBalance:   endingDebitBalance,
		endingCreditBalance:  endingCreditBalance,
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

func (l *AuxiliaryLedger) OpeningDebitBalance() decimal.Decimal {
	return l.openingDebitBalance
}

func (l *AuxiliaryLedger) OpeningCreditBalance() decimal.Decimal {
	return l.openingCreditBalance
}

func (l *AuxiliaryLedger) PeriodDebit() decimal.Decimal {
	return l.periodDebit
}

func (l *AuxiliaryLedger) PeriodCredit() decimal.Decimal {
	return l.periodCredit
}

func (l *AuxiliaryLedger) EndingDebitBalance() decimal.Decimal {
	return l.endingDebitBalance
}

func (l *AuxiliaryLedger) EndingCreditBalance() decimal.Decimal {
	return l.endingCreditBalance
}
