package datasource

import (
	"context"

	"gorm.io/gorm"
)

type ctxDBKey struct{}

func GetIfAbsentInContext(ctx context.Context, getter func() *gorm.DB) *gorm.DB {
	db := ctx.Value(ctxDBKey{}).(*gorm.DB)
	if db != nil {
		return db
	}
	return getter()
}

func WrapInNewContext(parent context.Context, db *gorm.DB) context.Context {
	return context.WithValue(parent, ctxDBKey{}, db)
}
