package account

import "github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_category"

func (a *Account) AssignAuxiliaryCategories(auxiliaryCategories []*auxiliary_category.AuxiliaryCategory) error {
	// check leaf node?
	// check superior account?

	a.auxiliaryCategories = auxiliaryCategories

	return nil
}
