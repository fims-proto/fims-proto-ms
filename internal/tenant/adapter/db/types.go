package db

import (
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"

	"github.com/google/uuid"
)

type Tenant struct {
	TenantId       uuid.UUID `db:"tenant_id"`
	Subdomain      string    `db:"subdomain"`
	DBConnPassword string    `db:"db_conn_password"`
}

func (t Tenant) mapToQuery() query.Tenant {
	return query.Tenant{
		TenantId:       t.TenantId,
		Subdomain:      t.Subdomain,
		DBConnPassword: t.DBConnPassword,
	}
}
