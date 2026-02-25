package ledger

import (
	"github.com/shopspring/decimal"
)

func (l *Ledger) UpdateBalance(amount decimal.Decimal) {
	l.periodAmount = l.periodAmount.Add(amount)
	l.endingAmount = l.openingAmount.Add(l.periodAmount)

	// Update performance fields (periodDebit, periodCredit)
	if amount.IsPositive() {
		l.periodDebit = l.periodDebit.Add(amount)
	} else {
		l.periodCredit = l.periodCredit.Add(amount.Abs())
	}
}

func (l *Ledger) UpdateOpeningBalance(balance decimal.Decimal) {
	l.openingAmount = balance

	// Update ending accordingly
	l.endingAmount = l.openingAmount.Add(l.periodAmount)
}
