package app

import "github/fims-proto/fims-proto-ms/internal/tenant/app/query"

type Queries struct {
	ReadTenants query.ReadTenantsHandler
}

type Application struct {
	Queries Queries
}

func NewApplication() Application {
	return Application{}
}

func (a *Application) Inject(readModel query.TenantsReadModel) {
	a.Queries = Queries{
		ReadTenants: query.NewReadTenantsHandler(readModel),
	}
}
