package db

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TenantPostgresRepository struct {
	db *gorm.DB
}

func NewTenantPostgresRepository(db *gorm.DB) *TenantPostgresRepository {
	if db == nil {
		panic("nil db connection")
	}
	return &TenantPostgresRepository{db: db}
}

func (t TenantPostgresRepository) ReadById(ctx context.Context, tenantId uuid.UUID) (query.Tenant, error) {
	dbTenant := tenant{}

	if err := t.db.WithContext(ctx).First(&dbTenant, "id = ?", tenantId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return query.Tenant{}, errors.Wrapf(err, "tenant %s does not exist", tenantId)
		} else {
			return query.Tenant{}, errors.Wrapf(err, "unknown error when get tenant %s", tenantId)
		}
	}

	return unmarshalToQuery(&dbTenant), nil
}

func (t TenantPostgresRepository) ReadBySubdomain(ctx context.Context, subdomain string) (query.Tenant, error) {
	dbTenant := tenant{}

	if err := t.db.WithContext(ctx).Where("subdomain = ?", subdomain).First(&dbTenant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return query.Tenant{}, errors.Wrapf(err, "tenant %s does not exist", subdomain)
		} else {
			return query.Tenant{}, errors.Wrapf(err, "unknown error when get tenant %s", subdomain)
		}
	}

	return unmarshalToQuery(&dbTenant), nil
}
