package filterable

import (
	"errors"
	"fmt"
	"strings"

	"github/fims-proto/fims-proto-ms/internal/common/data/schema"

	"github/fims-proto/fims-proto-ms/internal/common/data/field"

	"gorm.io/gorm"
)

func helpValueToString(values any) string {
	strVal, ok := values.(string)
	if ok {
		return fmt.Sprintf(`'%s'`, strVal)
	}
	return fmt.Sprintf("%v", values)
}

func assembleSQL(f Filterable, targetSchema schema.Schema) (string, error) {
	fType := f.FilterableType()
	switch fType {
	case TypeAND:
		{
			var whereStr []string
			whereStr = append(whereStr, "(")
			for _, child := range f.Children() {
				substr, err := assembleSQL(child, targetSchema)
				if err != nil {
					return substr, err
				}
				whereStr = append(whereStr, substr, "AND")
			}
			whereStr[len(whereStr)-1] = ")"
			return strings.Join(whereStr, " "), nil
		}
	case TypeOR:
		{
			var whereStr []string
			whereStr = append(whereStr, "(")
			for _, child := range f.Children() {
				substr, err := assembleSQL(child, targetSchema)
				if err != nil {
					return substr, err
				}
				whereStr = append(whereStr, substr, "OR")
			}
			whereStr[len(whereStr)-1] = ")"
			return strings.Join(whereStr, " "), nil
		}
	case TypeNOT:
		{
			var whereStr []string
			whereStr = append(whereStr, "NOT", "(")
			child := f.Children()[0]
			substr, err := assembleSQL(child, targetSchema)
			if err != nil {
				return substr, err
			}
			whereStr = append(whereStr, substr, ")")
			return strings.Join(whereStr, " "), nil
		}
	case TypeATOM:
		{
			// try to assert type into filterImpl
			var filter Filter
			ok := false
			filter, ok = f.(filterImpl)
			if !ok {
				return "failed", errors.New("type assertion failed for filterImpl")
			}
			fieldName, err := field.ToColumn(filter.Field(), targetSchema)
			if err != nil {
				return "failed", err
			}

			var whereStr []string
			whereStr = append(whereStr, "(", fieldName)
			switch filter.Operator() {
			case OptBtw:
				{
					whereStr = append(
						whereStr,
						filter.Operator().String(), helpValueToString(filter.Values()[0]),
						"AND",
						helpValueToString(filter.Values()[1]),
					)
				}
			case OptIn:
				{
					whereStr = append(whereStr, filter.Operator().String())
					var strValList []string
					for _, val := range filter.Values() {
						strValList = append(strValList, helpValueToString(val))
					}
					whereStr = append(whereStr, "(", strings.Join(strValList, ","), ")")
				}
			case OptEq, OptLt, OptLte, OptGt, OptGte:
				{
					whereStr = append(whereStr, filter.Operator().String(), helpValueToString(filter.Values()[0]))
				}
			case OptStw:
				{
					whereStr = append(whereStr, "LIKE")
					whereStr = append(whereStr, fmt.Sprintf("'%s%%'", filter.Values()[0]))
				}
			case OptCtn:
				{
					whereStr = append(whereStr, "LIKE")
					whereStr = append(whereStr, fmt.Sprintf(`'%%%s%%'`, filter.Values()[0]))
				}
			default:
				{
					return "failed", fmt.Errorf("currently not supported Operator Type %s", filter.Operator().String())
				}
			}
			whereStr = append(whereStr, ")")
			return strings.Join(whereStr, " "), nil
		}
	default:
		{
			return "failed", errors.New("unknow filterNode type")
		}
	}
}

func Filtering(f Filterable, targetSchema schema.Schema) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.IsFiltered() {
			strFilter, err := assembleSQL(f, targetSchema)
			if err != nil {
				_ = db.AddError(err)
				return db
			}
			db = db.Where(strFilter)
		}

		return db
	}
}
