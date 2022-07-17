package domain

import (
	"github.com/shopspring/decimal"
	commonAccount "github/fims-proto/fims-proto-ms/internal/common/account"
)

func (l *Ledger) UpdateBalance(debit, credit decimal.Decimal, balanceDirection commonAccount.Direction) {
	l.debit = debit
	l.credit = credit

	addingAmount := l.debit
	subAmount := l.credit

	if balanceDirection == commonAccount.Credit {
		addingAmount = l.credit
		subAmount = l.debit
	}

	l.endingBalance = l.openingBalance.Add(addingAmount).Sub(subAmount)
}
