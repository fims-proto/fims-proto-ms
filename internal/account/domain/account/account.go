package account

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/account/domain/account_configuration"
)

type Account struct {
	sobId          uuid.UUID
	accountId      uuid.UUID
	periodId       uuid.UUID
	openingBalance decimal.Decimal
	endingBalance  decimal.Decimal
	periodDebit    decimal.Decimal
	periodCredit   decimal.Decimal
	configuration  account_configuration.AccountConfiguration
}

func New(sobId, accountId, periodId uuid.UUID, openingBalance, endingBalance, periodDebit, periodCredit decimal.Decimal, configuration account_configuration.AccountConfiguration) (*Account, error) {
	if sobId == uuid.Nil {
		return nil, errors.New("nil sob id")
	}

	if accountId == uuid.Nil {
		return nil, errors.New("nil account id")
	}

	if periodId == uuid.Nil {
		return nil, errors.New("nil period id")
	}

	return &Account{
		sobId:          sobId,
		accountId:      accountId,
		periodId:       periodId,
		openingBalance: openingBalance,
		endingBalance:  endingBalance,
		periodDebit:    periodDebit,
		periodCredit:   periodCredit,
		configuration:  configuration,
	}, nil
}

func (a *Account) SobId() uuid.UUID {
	return a.sobId
}

func (a *Account) AccountId() uuid.UUID {
	return a.accountId
}

func (a *Account) PeriodId() uuid.UUID {
	return a.periodId
}

func (a *Account) OpeningBalance() decimal.Decimal {
	return a.openingBalance
}

func (a *Account) EndingBalance() decimal.Decimal {
	return a.endingBalance
}

func (a *Account) PeriodDebit() decimal.Decimal {
	return a.periodDebit
}

func (a *Account) PeriodCredit() decimal.Decimal {
	return a.periodCredit
}

func (a *Account) Configuration() account_configuration.AccountConfiguration {
	return a.configuration
}
