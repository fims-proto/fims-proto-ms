package app

import "github/fims-proto/fims-proto-ms/internal/tenant/app/query"

type Queries struct {
	TenantById        query.TenantByIdHandler
	TenantBySubdomain query.TenantBySubdomainHandler
}

type Application struct {
	Queries Queries
}

func NewApplication(readModel query.TenantReadModel) Application {
	return Application{
		Queries{
			TenantById:        query.NewTenantByIdHandler(readModel),
			TenantBySubdomain: query.NewTenantBySubdomain(readModel),
		},
	}
}
