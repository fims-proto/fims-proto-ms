package http

import (
	"github/fims-proto/fims-proto-ms/internal/sob/app"

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
	r.GET("/sobs", h.SearchSobs)
	r.GET("/sobs/:sobId", h.ReadSobById)
	r.PATCH("/sobs/:sobId", h.UpdateSob)
	r.POST("/sobs", h.CreateSob)
}
