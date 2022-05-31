package data

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func EnrichDb(pageable Pageable, db *gorm.DB) *gorm.DB {
	// page
	db = db.Offset(pageable.Offset()).Limit(pageable.Size())

	// sort
	if pageable.Sorts() != nil {
		var orderStr []string
		for _, sort := range pageable.Sorts() {
			orderStr = append(orderStr, strings.Join([]string{sort.Field(), sort.Order()}, " "))
		}
		db = db.Order(strings.Join(orderStr, ","))
	}

	// choose
	if pageable.Chooses() != nil {
		var selectStr []string
		for _, choose := range pageable.Chooses() {
			selectStr = append(selectStr, choose.Field())
		}
		db = db.Select(strings.Join(selectStr, ","))
	}

	// filter
	if pageable.Filters() != nil {
		var whereStr []string
		var args []any
		for _, filter := range pageable.Filters() {
			// field
			whereStr = append(whereStr, filter.Field())

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
