package adapter

import (
	"context"
	"database/sql"
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type TenantPostgresRepository struct {
	db *sqlx.DB
}

func NewTenantPostgresRepository(db *sqlx.DB) *TenantPostgresRepository {
	if db == nil {
		panic("nil db connection")
	}
	return &TenantPostgresRepository{db: db}
}

func (t TenantPostgresRepository) ReadByUUID(ctx context.Context, tenantId uuid.UUID) (query.Tenant, error) {
	tenant := Tenant{}
	err := t.db.GetContext(ctx, &tenant, "SELECT * FROM tenant WHERE tenant_id = $1", tenantId)
	if err == sql.ErrNoRows {
		return query.Tenant{}, errors.Wrapf(err, "tenant %s does not exist", tenantId)
	} else if err != nil {
		return query.Tenant{}, errors.Wrapf(err, "unknown error when get tenant %s", tenantId)
	}

	return tenant.mapToQuery(), nil
}
