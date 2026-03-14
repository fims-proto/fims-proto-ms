package account

import "github.com/google/uuid"

func (a *Account) UpdateDimensionCategories(categoryIds []uuid.UUID) {
	a.dimensionCategoryIds = categoryIds
}
