package adapter

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type tenantManager interface {
	GetDBConn(tenantId uuid.UUID) (*sqlx.DB, error)
}

type VoucherPostgresRepository struct {
	tenantManager tenantManager
}
