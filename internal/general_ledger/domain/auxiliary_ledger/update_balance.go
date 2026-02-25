package auxiliary_ledger

import (
	"github.com/shopspring/decimal"
)

func (l *AuxiliaryLedger) UpdateBalance(amount decimal.Decimal) {
	l.periodAmount = l.periodAmount.Add(amount)
	l.endingAmount = l.openingAmount.Add(l.periodAmount)

	// Update performance fields (periodDebit, periodCredit)
	if amount.IsPositive() {
		l.periodDebit = l.periodDebit.Add(amount)
	} else {
		l.periodCredit = l.periodCredit.Add(amount.Abs())
	}
}
