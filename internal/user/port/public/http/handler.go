package http

import (
	"github/fims-proto/fims-proto-ms/internal/user/app"

	"github.com/gin-gonic/gin"
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

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/user/:userId", h.ReadUserById)
	r.PATCH("/user/:userId", h.UpdateUser)
}
