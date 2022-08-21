package tenant

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Tenant struct {
	id        uuid.UUID
	subdomain string
	dsn       string
}

func New(id uuid.UUID, subdomain, dsn string) (*Tenant, error) {
	if id == uuid.Nil {
		return nil, errors.New("nil id")
	}

	if subdomain == "" {
		return nil, errors.New("empty subdomain")
	}

	if dsn == "" {
		return nil, errors.New("empty dsn")
	}

	return &Tenant{
		id:        id,
		subdomain: subdomain,
		dsn:       dsn,
	}, nil
}

func (t Tenant) Id() uuid.UUID {
	return t.id
}

func (t Tenant) Subdomain() string {
	return t.subdomain
}

func (t Tenant) DSN() string {
	return t.dsn
}
