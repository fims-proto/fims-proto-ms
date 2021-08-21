package db

import (
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	ID             uuid.UUID `gorm:"type:uuid"`
	Subdomain      string
	DBConnPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (t Tenant) mapToQuery() query.Tenant {
	return query.Tenant{
		TenantId:       t.ID,
		Subdomain:      t.Subdomain,
		DBConnPassword: t.DBConnPassword,
	}
}
