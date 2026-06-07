package http

import (
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/dimension/app"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	app *app.Application
}

func NewHandler(app *app.Application) Handler {
	if app == nil {
		panic("nil application")
	}
	return Handler{app: app}
}

func InitRouter(h Handler, api huma.API) {
	dimensionApi := data.WithTag(api, "dimension")

	huma.Register(dimensionApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/dimension/categories", Summary: "Search dimension categories"}, h.SearchCategories)
	huma.Register(dimensionApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/dimension/categories", Summary: "Create dimension category", DefaultStatus: 204}, h.CreateCategory)
	huma.Register(dimensionApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/dimension/category/{categoryId}", Summary: "Get dimension category by ID"}, h.ReadCategoryById)
	huma.Register(dimensionApi, huma.Operation{Method: "PATCH", Path: "/api/v1/sob/{sobId}/dimension/category/{categoryId}", Summary: "Update dimension category", DefaultStatus: 204}, h.UpdateCategory)
	huma.Register(dimensionApi, huma.Operation{Method: "DELETE", Path: "/api/v1/sob/{sobId}/dimension/category/{categoryId}", Summary: "Delete dimension category", DefaultStatus: 204}, h.DeleteCategory)

	huma.Register(dimensionApi, huma.Operation{Method: "GET", Path: "/api/v1/sob/{sobId}/dimension/category/{categoryId}/options", Summary: "Search dimension options"}, h.SearchOptions)
	huma.Register(dimensionApi, huma.Operation{Method: "POST", Path: "/api/v1/sob/{sobId}/dimension/category/{categoryId}/options", Summary: "Create dimension option", DefaultStatus: 204}, h.CreateOption)
	huma.Register(dimensionApi, huma.Operation{Method: "PATCH", Path: "/api/v1/sob/{sobId}/dimension/category/{categoryId}/option/{optionId}", Summary: "Update dimension option", DefaultStatus: 204}, h.UpdateOption)
	huma.Register(dimensionApi, huma.Operation{Method: "DELETE", Path: "/api/v1/sob/{sobId}/dimension/category/{categoryId}/option/{optionId}", Summary: "Delete dimension option", DefaultStatus: 204}, h.DeleteOption)
}
