package ledger

import (
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/account/balance_direction"
)

func (l *Ledger) UpdateEndingBalance(debit, credit decimal.Decimal) {
	l.periodDebit = l.periodDebit.Add(debit)
	l.periodCredit = l.periodCredit.Add(credit)

	l.endingDebitBalance = l.openingDebitBalance.Add(l.periodDebit)
	l.endingCreditBalance = l.openingCreditBalance.Add(l.periodCredit)
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
