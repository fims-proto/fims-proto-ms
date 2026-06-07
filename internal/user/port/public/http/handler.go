package http

import (
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/user/app"

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
	usersApi := data.WithTag(api, "users")

	huma.Register(usersApi, huma.Operation{
		Method:      "GET",
		Path:        "/api/v1/user/{userId}",
		Summary:     "Get user by ID",
		OperationID: "readUserById",
	}, h.ReadUserById)

	huma.Register(usersApi, data.WithResponseLinks(huma.Operation{Method: "PATCH", Path: "/api/v1/user/{userId}", Summary: "Update user", OperationID: "updateUser", DefaultStatus: 204}, 204, map[string]*huma.Link{
		"readUser": {
			OperationID: "readUserById",
			Parameters:  map[string]any{"userId": "$request.path.userId"},
			Description: "Read updated user.",
		},
	}), h.UpdateUser)
}
