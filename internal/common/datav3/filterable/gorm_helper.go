package filterable

import (
	"fmt"
	"strings"

	"github/fims-proto/fims-proto-ms/internal/common/datav3/field"
	"gorm.io/gorm"
)

func Filtering(f Filterable, resolveEntity func(entity string) (string, error)) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.IsFiltered() {
			var whereStr []string
			var args []any
			for _, filter := range f.Filters() {
				// field
				fieldName, err := field.ToColumn(filter.Field(), resolveEntity)
				if err != nil {
					_ = db.AddError(err)
					return db
				}
				whereStr = append(whereStr, fieldName)

				// operator and variable placeholder, and args
				switch filter.Operator() {
				case OptBt:
					whereStr = append(whereStr, filter.Operator().String())
					whereStr = append(whereStr, "? AND ?")
					args = append(args, filter.Values()[0], filter.Values()[1])
				case OptIn:
					whereStr = append(whereStr, filter.Operator().String())
					whereStr = append(whereStr, "?")
					args = append(args, filter.Values())
				case OptStartsWith:
					whereStr = append(whereStr, "LIKE")
					whereStr = append(whereStr, "?")
					args = append(args, fmt.Sprintf("%s%%", filter.Values()[0]))
				default:
					whereStr = append(whereStr, filter.Operator().String())
					whereStr = append(whereStr, "?")
					args = append(args, filter.Values()[0])
				}

				whereStr = append(whereStr, "AND")
			}
			whereStr = whereStr[0 : len(whereStr)-1] // remove last AND

			db = db.Where(strings.Join(whereStr, " "), args...)
		}

		return db
	}
}
