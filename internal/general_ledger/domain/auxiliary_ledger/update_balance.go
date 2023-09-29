package auxiliary_ledger

import (
	"github.com/shopspring/decimal"
)

func (l *AuxiliaryLedger) UpdateBalance(debit, credit decimal.Decimal) {
	l.periodDebit = l.periodDebit.Add(debit)
	l.periodCredit = l.periodCredit.Add(credit)

	adding := l.periodDebit
	subtracting := l.periodCredit

	l.endingBalance = l.openingBalance.Add(adding).Sub(subtracting)
}
