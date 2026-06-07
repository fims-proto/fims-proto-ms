package http

import (
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/sob/app"

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
	sobsApi := data.WithTag(api, "sobs")

	huma.Register(sobsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sobs",
		Summary:     "List sobs",
		OperationID: "searchSobs",
	}, h.SearchSobs)

	huma.Register(sobsApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/sobs/{sobId}",
		Summary:     "Get sob by ID",
		OperationID: "readSobById",
	}, h.ReadSobById)

	huma.Register(sobsApi, data.WithResponseLinks(huma.Operation{
		Method:        "PATCH",
		Path:          "/api/v1/sobs/{sobId}",
		Summary:       "Update sob",
		OperationID:   "updateSob",
		DefaultStatus: 204,
	}, 204, map[string]*huma.Link{
		"readSob": {
			OperationID: "readSobById",
			Parameters:  map[string]any{"sobId": "$request.path.sobId"},
			Description: "Read updated SoB.",
		},
	}), h.UpdateSob)

	huma.Register(sobsApi, data.WithResponseLinks(huma.Operation{
		Method:        "POST",
		Path:          "/api/v1/sobs",
		Summary:       "Create sob",
		OperationID:   "createSob",
		DefaultStatus: 201,
	}, 201, map[string]*huma.Link{
		"readSob": {
			OperationID: "readSobById",
			Parameters:  map[string]any{"sobId": "$response.body#/id"},
			Description: "Read created SoB.",
		},
	}), h.CreateSob)
}
