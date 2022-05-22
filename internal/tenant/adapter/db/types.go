package db

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"

	"github.com/google/uuid"
)

type tenant struct {
	Id        uuid.UUID `gorm:"type:uuid"`
	Subdomain string
	DSN       string
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
}

func unmarshallToQuery(t *tenant) query.Tenant {
	return query.Tenant{
		TenantId:  t.Id,
		Subdomain: t.Subdomain,
		DSN:       t.DSN,
	}
}
