package dedicated_datasource

import (
	"context"
	"sync"

	"github/fims-proto/fims-proto-ms/internal/common/config"
	"github/fims-proto/fims-proto-ms/internal/common/datasource"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"gorm.io/gorm"
)

// DedicatedDataSource provide only the datasource that specified in the environment
type DedicatedDataSource struct {
	get func() *gorm.DB
}

func NewDedicatedDataSource() *DedicatedDataSource {
	return &DedicatedDataSource{
		get: sync.OnceValue(func() *gorm.DB {
			connector := datasource.NewConnector()
			connection, err := connector.GetConnection(config.GetString("postgres.dsn"))
			if err != nil {
				panic(err)
			}
			log.DebugWithoutCxt("database connection acquired %s", connection.Name())
			return connection
		}),
	}
}

func (d *DedicatedDataSource) GetConnection(ctx context.Context) *gorm.DB {
	// get connection from context first, since transactional connection is passed via context
	return datasource.GetIfAbsentInContext(ctx, d.get)
}

func (d *DedicatedDataSource) EnableTransaction(ctx context.Context, transactionalFn func(txCtx context.Context) error) error {
	db := d.GetConnection(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		return transactionalFn(datasource.WrapInNewContext(ctx, tx))
	})
}
