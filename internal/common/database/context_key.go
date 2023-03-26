package database

import (
	"context"

	"gorm.io/gorm"
)

const ctxDBKey string = "fims/db"

func GetContextDBKey() string {
	return ctxDBKey
}

func ReadDBFromContext(ctx context.Context) *gorm.DB {
	return ctx.Value(ctxDBKey).(*gorm.DB)
}

func ContextWithDB(parent context.Context, db *gorm.DB) context.Context {
	return context.WithValue(parent, ctxDBKey, db)
}
