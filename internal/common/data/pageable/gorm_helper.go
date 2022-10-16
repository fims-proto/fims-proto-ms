package pageable

import (
	"gorm.io/gorm"
)

func Paging(p Pageable) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.IsPaged() {
			db = db.Offset(p.Offset()).Limit(p.PageSize())
		}

		return db
	}
}
