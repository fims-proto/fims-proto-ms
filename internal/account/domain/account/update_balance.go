package account

import (
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/account/domain/balance_direction"
)

func (a *Account) UpdateBalance(debit, credit decimal.Decimal) {
	a.periodDebit = a.periodDebit.Add(debit)
	a.periodCredit = a.periodCredit.Add(credit)

	adding := a.periodDebit
	subtracting := a.periodCredit

	if a.configuration.BalanceDirection() == balance_direction.Credit {
		adding = a.periodCredit
		subtracting = a.periodDebit
	}

	a.endingBalance = a.openingBalance.Add(adding).Sub(subtracting)
}
