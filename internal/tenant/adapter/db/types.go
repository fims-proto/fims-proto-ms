package db

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"

	"github.com/google/uuid"
)

type tenantPO struct {
	Id        uuid.UUID `gorm:"type:uuid"`
	Subdomain string
	DSN       string
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

// table names

func (t tenantPO) TableName() string {
	return "tenants"
}

// mappers

func tenantPOToDTO(po tenantPO) query.Tenant {
	return query.Tenant{
		TenantId:  po.Id,
		Subdomain: po.Subdomain,
		DSN:       po.DSN,
	}
}
