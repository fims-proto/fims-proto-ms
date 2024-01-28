package datasource

import (
	"context"

	"gorm.io/gorm"
)

type DataSource interface {
	GetConnection(ctx context.Context) *gorm.DB
	EnableTransaction(ctx context.Context, transactionalFn func(txCtx context.Context) error) error
}
