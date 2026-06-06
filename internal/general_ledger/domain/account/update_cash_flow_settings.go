package account

import "github.com/google/uuid"

func (a *Account) UpdateCashEquivalent(isCashEquivalent bool) {
	a.isCashEquivalent = isCashEquivalent
}

func (a *Account) UpdateDefaultCashFlowItems(debitItemId, creditItemId *uuid.UUID) {
	a.defaultCashFlowItemIdForDebit = debitItemId
	a.defaultCashFlowItemIdForCredit = creditItemId
}
