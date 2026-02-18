package multitenant_datasource

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/datasource"

	"gorm.io/gorm"
)

type MultiTenantDataSource struct{}

func NewMultiTenantDataSource() *MultiTenantDataSource {
	return &MultiTenantDataSource{}
}

func (m *MultiTenantDataSource) GetConnection(_ context.Context) *gorm.DB {
	panic("not implemented")
	// 0. check if context already has connection
	// 1. read subdomain from context
	// 2. remote call to tenant manager module to acquire database information/credential
	// 3. open database connection
}

func (m *MultiTenantDataSource) EnableTransaction(ctx context.Context, transactionalFn func(txCtx context.Context) error) error {
	db := m.GetConnection(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		return transactionalFn(datasource.WrapInNewContext(ctx, tx))
	})
}
