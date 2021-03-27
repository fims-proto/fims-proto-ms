package domain

import (
	commonaccount "github/fims-proto/fims-proto-ms/internal/common/account"

	"github.com/shopspring/decimal"
)

func (l *Ledger) UpdateBalance(debit decimal.Decimal, credit decimal.Decimal) error {
	l.debit.Add(debit)
	l.credit.Add(credit)

	switch l.AccountType() {
	case commonaccount.Assets, commonaccount.Cost:
		l.balance = l.balance.Add(debit).Sub(credit)
	case commonaccount.Liabilities, commonaccount.Equity:
		l.balance = l.balance.Add(credit).Sub(debit)
	case commonaccount.ProfitAndLoss:
		// no balance for this type of account
	case commonaccount.Common:
		// TODO not sure how to handle
		panic("common account type not supported yet")
	}

	return nil
}
