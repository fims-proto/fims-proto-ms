package datasource

import (
	"context"

	"gorm.io/gorm"
)

type ctxDBKey string

const key ctxDBKey = "dataSource"

func GetIfAbsentInContext(ctx context.Context, getter func() *gorm.DB) *gorm.DB {
	db, ok := ctx.Value(key).(*gorm.DB)
	if !ok || db == nil {
		return getter()
	}
	return db
}

func WrapInNewContext(parent context.Context, db *gorm.DB) context.Context {
	return context.WithValue(parent, key, db)
}
