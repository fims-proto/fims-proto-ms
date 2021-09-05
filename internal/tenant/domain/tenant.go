package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Tenant struct {
	tenantId       uuid.UUID
	subdomain      string
	dbConnPassword string
}

func NewTenant(tenantId uuid.UUID, subdomain, dbConnPassword string) (*Tenant, error) {
	if tenantId == uuid.Nil {
		return nil, errors.New("nil tenantId")
	}
	if subdomain == "" {
		return nil, errors.New("empty subdomain")
	}
	if dbConnPassword == "" {
		return nil, errors.New("empty DB connection password")
	}

	return &Tenant{
		tenantId:       tenantId,
		subdomain:      subdomain,
		dbConnPassword: dbConnPassword,
	}, nil
}

func (t Tenant) TenantId() uuid.UUID {
	return t.tenantId
}

func (t Tenant) Subdomain() string {
	return t.subdomain
}

func (t Tenant) DBConnPassword() string {
	return t.dbConnPassword
}
