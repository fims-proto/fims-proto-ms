package account

import "github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account_category"

func (a *Account) AssignAuxiliaryAccountCategories(auxiliaryAccountCategories []*auxiliary_account_category.AuxiliaryAccountCategory) {
	// check leaf node?
	// check superior account?

	a.auxiliaryAccountCategories = auxiliaryAccountCategories
}
