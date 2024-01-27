package db

import (
	"context"
	"errors"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"

	"github.com/google/uuid"
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
			return query.Tenant{}, fmt.Errorf("tenant %s does not exist: %w", tenantId, err)
		} else {
			return query.Tenant{}, fmt.Errorf("unknown error when get tenant %s: %w", tenantId, err)
		}
	}

	return tenantPOToDTO(po), nil
}

func (t TenantPostgresRepository) TenantBySubdomain(ctx context.Context, subdomain string) (query.Tenant, error) {
	po := tenantPO{}

	if err := t.db.WithContext(ctx).Where("subdomain = ?", subdomain).First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return query.Tenant{}, fmt.Errorf("tenant %s does not exist: %w", subdomain, err)
		} else {
			return query.Tenant{}, fmt.Errorf("unknown error when get tenant %s: %w", subdomain, err)
		}
	}

	return tenantPOToDTO(po), nil
}
