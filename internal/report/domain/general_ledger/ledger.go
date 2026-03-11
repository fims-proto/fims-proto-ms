package general_ledger

import (
	"github.com/shopspring/decimal"
)

type Ledger struct {
	account       *Account
	period        *Period
	openingAmount decimal.Decimal
	periodAmount  decimal.Decimal
	periodDebit   decimal.Decimal // positive amount, only for query performance
	periodCredit  decimal.Decimal // positive amount, only for query performance
	endingAmount  decimal.Decimal
}

func NewLedger(
	account *Account,
	period *Period,
	openingAmount decimal.Decimal,
	periodAmount decimal.Decimal,
	periodDebit decimal.Decimal,
	periodCredit decimal.Decimal,
	endingAmount decimal.Decimal,
) *Ledger {
	return &Ledger{
		account:       account,
		period:        period,
		openingAmount: openingAmount,
		periodAmount:  periodAmount,
		periodDebit:   periodDebit,
		periodCredit:  periodCredit,
		endingAmount:  endingAmount,
	}
}

func (l Ledger) Account() *Account {
	return l.account
}

func (l Ledger) Period() *Period {
	return l.period
}

func (l Ledger) OpeningAmount() decimal.Decimal {
	return l.openingAmount
}

func (l Ledger) PeriodAmount() decimal.Decimal {
	return l.periodAmount
}

func (l Ledger) PeriodDebit() decimal.Decimal {
	return l.periodDebit
}

func (l Ledger) PeriodCredit() decimal.Decimal {
	return l.periodCredit
}

func (l Ledger) EndingAmount() decimal.Decimal {
	return l.endingAmount
}
