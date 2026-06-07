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

	huma.Register(sobsApi, huma.Operation{Method: "GET", Path: "/api/v1/sobs", Summary: "List sobs"}, h.SearchSobs)
	huma.Register(sobsApi, huma.Operation{Method: "GET", Path: "/api/v1/sobs/{sobId}", Summary: "Get sob by ID"}, h.ReadSobById)
	huma.Register(sobsApi, huma.Operation{Method: "PATCH", Path: "/api/v1/sobs/{sobId}", Summary: "Update sob", DefaultStatus: 204}, h.UpdateSob)
	huma.Register(sobsApi, huma.Operation{Method: "POST", Path: "/api/v1/sobs", Summary: "Create sob", DefaultStatus: 201}, h.CreateSob)
}
