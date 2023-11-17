package ledger

import (
	"github.com/shopspring/decimal"
)

func (l *Ledger) UpdateBalance(debit, credit decimal.Decimal) {
	l.periodDebit = l.periodDebit.Add(debit)
	l.periodCredit = l.periodCredit.Add(credit)

	l.endingDebitBalance = l.openingDebitBalance.Add(l.periodDebit)
	l.endingCreditBalance = l.openingCreditBalance.Add(l.periodCredit)
}
