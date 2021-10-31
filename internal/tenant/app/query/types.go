package query

import "github.com/google/uuid"

type Tenant struct {
	TenantId       uuid.UUID
	Subdomain      string
	DBConnPassword string
}
