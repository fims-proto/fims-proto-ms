package ledger

import (
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/balance_direction"

	"github.com/shopspring/decimal"
)

func (l *Ledger) UpdateEndingBalance(debit, credit decimal.Decimal) {
	l.periodDebit = l.periodDebit.Add(debit)
	l.periodCredit = l.periodCredit.Add(credit)

	// 1. ending balance = (opening debit - opening credit) + (period debit - period credit)
	// 2. if ending balance positive, save it as ending debit, otherwise ending credit
	endingBalance := (l.openingDebitBalance.Sub(l.openingCreditBalance)).Add(l.periodDebit.Sub(l.periodCredit))
	if endingBalance.IsPositive() {
		l.endingDebitBalance = endingBalance
		l.endingCreditBalance = decimal.Zero
	} else {
		l.endingCreditBalance = endingBalance.Neg()
		l.endingDebitBalance = decimal.Zero
	}
}

func (l *Ledger) UpdateOpeningBalance(balance decimal.Decimal) {
	debit := decimal.Zero
	credit := decimal.Zero

	if l.account.BalanceDirection() == balance_direction.Debit {
		debit = balance
	}
	if l.account.BalanceDirection() == balance_direction.Credit {
		credit = balance
	}

	l.openingDebitBalance = debit
	l.openingCreditBalance = credit

	// update ending accordingly
	l.UpdateEndingBalance(decimal.Zero, decimal.Zero)
}
