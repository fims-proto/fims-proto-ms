package sortable

import (
	"strings"

	"github/fims-proto/fims-proto-ms/internal/common/datav3/field"
	"gorm.io/gorm"
)

func Sorting(s Sortable, resolveEntity func(entity string) (string, error)) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if s.IsSorted() {
			var orderStr []string
			for _, sort := range s.Sorts() {
				fieldName, err := field.ToColumn(sort.Field(), resolveEntity)
				if err != nil {
					_ = db.AddError(err)
					return db
				}
				orderStr = append(orderStr, strings.Join([]string{fieldName, sort.Order()}, " "))
			}
			db = db.Order(strings.Join(orderStr, ","))
		}

		return db
	}
}
