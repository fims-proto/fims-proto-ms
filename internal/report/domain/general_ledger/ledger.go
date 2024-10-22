package general_ledger

import (
	"github.com/shopspring/decimal"
)

type Ledger struct {
	account              *Account
	period               *Period
	openingDebitBalance  decimal.Decimal
	openingCreditBalance decimal.Decimal
	periodDebit          decimal.Decimal
	periodCredit         decimal.Decimal
	endingDebitBalance   decimal.Decimal
	endingCreditBalance  decimal.Decimal
}

func NewLedger(
	account *Account,
	period *Period,
	openingDebitBalance decimal.Decimal,
	openingCreditBalance decimal.Decimal,
	periodDebit decimal.Decimal,
	periodCredit decimal.Decimal,
	endingDebitBalance decimal.Decimal,
	endingCreditBalance decimal.Decimal,
) *Ledger {
	return &Ledger{
		account:              account,
		period:               period,
		openingDebitBalance:  openingDebitBalance,
		openingCreditBalance: openingCreditBalance,
		periodDebit:          periodDebit,
		periodCredit:         periodCredit,
		endingDebitBalance:   endingDebitBalance,
		endingCreditBalance:  endingCreditBalance,
	}
}

func (l Ledger) Account() *Account {
	return l.account
}

func (l Ledger) Period() *Period {
	return l.period
}

func (l Ledger) OpeningDebitBalance() decimal.Decimal {
	return l.openingDebitBalance
}

func (l Ledger) OpeningCreditBalance() decimal.Decimal {
	return l.openingCreditBalance
}

func (l Ledger) PeriodDebit() decimal.Decimal {
	return l.periodDebit
}

func (l Ledger) PeriodCredit() decimal.Decimal {
	return l.periodCredit
}

func (l Ledger) EndingDebitBalance() decimal.Decimal {
	return l.endingDebitBalance
}

func (l Ledger) EndingCreditBalance() decimal.Decimal {
	return l.endingCreditBalance
}
