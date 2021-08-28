package db

import (
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tenant struct {
	Id             uuid.UUID `gorm:"type:uuid"`
	Subdomain      string
	DBConnPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func unmarshallToQuery(t *tenant) query.Tenant {
	return query.Tenant{
		TenantId:       t.Id,
		Subdomain:      t.Subdomain,
		DBConnPassword: t.DBConnPassword,
	}
}
