package app

import "github/fims-proto/fims-proto-ms/internal/tenant/app/query"

type Queries struct {
	ReadTenants query.ReadTenantsHandler
}

type Application struct {
	Queries Queries
}

func NewApplication(readModel query.TenantsReadModel) Application {
	return Application{
		Queries{
			ReadTenants: query.NewReadTenantsHandler(readModel),
		},
	}
}
