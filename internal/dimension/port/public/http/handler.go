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

	huma.Register(dimensionApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/dimension/categories",
		Summary:     "Search dimension categories",
		OperationID: "searchDimensionCategories",
	}, h.SearchCategories)

	huma.Register(dimensionApi, huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/dimension/categories",
		Summary:     "Create dimension category",
		OperationID: "createDimensionCategory", DefaultStatus: 204,
	}, h.CreateCategory)

	huma.Register(dimensionApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/dimension/category/{categoryId}",
		Summary:     "Get dimension category by ID",
		OperationID: "readDimensionCategoryById",
	}, h.ReadCategoryById)

	huma.Register(dimensionApi, data.WithResponseLinks(huma.Operation{
		Method:      "PATCH",
		Path:        "/api/v1/sob/{sobId}/dimension/category/{categoryId}",
		Summary:     "Update dimension category",
		OperationID: "updateDimensionCategory", DefaultStatus: 204,
	}, 204, map[string]*huma.Link{
		"readDimensionCategory": {
			OperationID: "readDimensionCategoryById",
			Parameters: map[string]any{
				"sobId":      "$request.path.sobId",
				"categoryId": "$request.path.categoryId",
			},
			Description: "Read updated dimension category.",
		},
	}), h.UpdateCategory)

	huma.Register(dimensionApi, huma.Operation{
		Method:      "DELETE",
		Path:        "/api/v1/sob/{sobId}/dimension/category/{categoryId}",
		Summary:     "Delete dimension category",
		OperationID: "deleteDimensionCategory", DefaultStatus: 204,
	}, h.DeleteCategory)

	huma.Register(dimensionApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sob/{sobId}/dimension/category/{categoryId}/options",
		Summary:     "Search dimension options",
		OperationID: "searchDimensionOptions",
	}, h.SearchOptions)

	huma.Register(dimensionApi, huma.Operation{
		Method:      "POST",
		Path:        "/api/v1/sob/{sobId}/dimension/category/{categoryId}/options",
		Summary:     "Create dimension option",
		OperationID: "createDimensionOption", DefaultStatus: 204,
	}, h.CreateOption)

	huma.Register(dimensionApi, huma.Operation{
		Method:      "PATCH",
		Path:        "/api/v1/sob/{sobId}/dimension/category/{categoryId}/option/{optionId}",
		Summary:     "Update dimension option",
		OperationID: "updateDimensionOption", DefaultStatus: 204,
	}, h.UpdateOption)

	huma.Register(dimensionApi, huma.Operation{
		Method:      "DELETE",
		Path:        "/api/v1/sob/{sobId}/dimension/category/{categoryId}/option/{optionId}",
		Summary:     "Delete dimension option",
		OperationID: "deleteDimensionOption", DefaultStatus: 204,
	}, h.DeleteOption)
}
