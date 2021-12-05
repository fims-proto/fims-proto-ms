package domain

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Tenant struct {
	tenantId  uuid.UUID
	subdomain string
	dsn       string
}

func NewTenant(tenantId uuid.UUID, subdomain, dsn string) (*Tenant, error) {
	if tenantId == uuid.Nil {
		return nil, errors.New("nil tenantId")
	}
	if subdomain == "" {
		return nil, errors.New("empty subdomain")
	}
	if dsn == "" {
		return nil, errors.New("empty dsn")
	}

	return &Tenant{
		tenantId:  tenantId,
		subdomain: subdomain,
		dsn:       dsn,
	}, nil
}

func (t Tenant) TenantId() uuid.UUID {
	return t.tenantId
}

func (t Tenant) Subdomain() string {
	return t.subdomain
}

func (t Tenant) DSN() string {
	return t.dsn
}
