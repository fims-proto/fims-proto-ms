package auxiliary_ledger

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
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

func (a *AuxiliaryLedger) Id() uuid.UUID {
	return a.id
}

func (a *AuxiliaryLedger) PeriodId() uuid.UUID {
	return a.periodId
}

func (a *AuxiliaryLedger) AuxiliaryAccount() *auxiliary_account.AuxiliaryAccount {
	return a.auxiliaryAccount
}

func (a *AuxiliaryLedger) OpeningBalance() decimal.Decimal {
	return a.openingBalance
}

func (a *AuxiliaryLedger) EndingBalance() decimal.Decimal {
	return a.endingBalance
}

func (a *AuxiliaryLedger) PeriodDebit() decimal.Decimal {
	return a.periodDebit
}

func (a *AuxiliaryLedger) PeriodCredit() decimal.Decimal {
	return a.periodCredit
}
