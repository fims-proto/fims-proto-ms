package datasource

import (
	"context"
	"database/sql"

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

func HasTransactionInContext(ctx context.Context) bool {
	db, ok := ctx.Value(key).(*gorm.DB)
	if !ok || db == nil {
		return false
	}

	// Check if the DB is actually in a transaction by inspecting the connection pool
	// If it's a transaction, the ConnPool will be a *sql.Tx
	if _, ok = db.Statement.ConnPool.(*sql.Tx); ok {
		return true
	}

	return false
}

func WrapInNewContext(parent context.Context, db *gorm.DB) context.Context {
	return context.WithValue(parent, key, db)
}
