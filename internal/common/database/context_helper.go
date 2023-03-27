package database

import (
	"context"

	"gorm.io/gorm"
)

type ctxDBKey struct{}

func ReadDBFromContext(ctx context.Context) *gorm.DB {
	return ctx.Value(ctxDBKey{}).(*gorm.DB)
}

func NewContextWithDB(parent context.Context, db *gorm.DB) context.Context {
	return context.WithValue(parent, ctxDBKey{}, db)
}
