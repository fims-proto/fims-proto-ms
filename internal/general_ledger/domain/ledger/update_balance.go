package ledger

import (
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/balance_direction"
)

func (l *Ledger) UpdateBalance(debit, credit decimal.Decimal) {
	l.periodDebit = l.periodDebit.Add(debit)
	l.periodCredit = l.periodCredit.Add(credit)

	adding := l.periodDebit
	subtracting := l.periodCredit

	if l.account.BalanceDirection() == balance_direction.Credit {
		adding = l.periodCredit
		subtracting = l.periodDebit
	}

	l.endingBalance = l.openingBalance.Add(adding).Sub(subtracting)
}
