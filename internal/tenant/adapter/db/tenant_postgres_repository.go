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

func (t TenantPostgresRepository) TenantById(ctx context.Context, tenantId uuid.UUID) (query.Tenant, error) {
	po := tenantPO{}

	if err := t.db.WithContext(ctx).First(&po, "id = ?", tenantId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return query.Tenant{}, errors.Wrapf(err, "tenant %s does not exist", tenantId)
		} else {
			return query.Tenant{}, errors.Wrapf(err, "unknown error when get tenant %s", tenantId)
		}
	}

	return tenantPOToDTO(po), nil
}

func (t TenantPostgresRepository) TenantBySubdomain(ctx context.Context, subdomain string) (query.Tenant, error) {
	po := tenantPO{}

	if err := t.db.WithContext(ctx).Where("subdomain = ?", subdomain).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return query.Tenant{}, errors.Wrapf(err, "tenant %s does not exist", subdomain)
		} else {
			return query.Tenant{}, errors.Wrapf(err, "unknown error when get tenant %s", subdomain)
		}
	}

	return tenantPOToDTO(po), nil
}
