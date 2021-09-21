package db

import (
	"github/fims-proto/fims-proto-ms/internal/tenant/app/query"
	"time"

	"github.com/google/uuid"
)

type tenant struct {
	Id              uuid.UUID `gorm:"type:uuid"`
	Subdomain       string
	DBConnPassword  string
	KratosServerUrl string
	CreatedAt       time.Time `gorm:"<-:create"`
	UpdatedAt       time.Time
}

func unmarshallToQuery(t *tenant) query.Tenant {
	return query.Tenant{
		TenantId:        t.Id,
		Subdomain:       t.Subdomain,
		DBConnPassword:  t.DBConnPassword,
		KratosServerUrl: t.KratosServerUrl,
	}
}
