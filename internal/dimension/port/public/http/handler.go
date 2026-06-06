package http

import (
	"github/fims-proto/fims-proto-ms/internal/dimension/app"

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
	r.GET("/sob/:sobId/dimension/categories", h.SearchCategories)
	r.POST("/sob/:sobId/dimension/categories", h.CreateCategory)
	r.GET("/sob/:sobId/dimension/category/:categoryId", h.ReadCategoryById)
	r.PATCH("/sob/:sobId/dimension/category/:categoryId", h.UpdateCategory)
	r.DELETE("/sob/:sobId/dimension/category/:categoryId", h.DeleteCategory)

	r.GET("/sob/:sobId/dimension/category/:categoryId/options", h.SearchOptions)
	r.POST("/sob/:sobId/dimension/category/:categoryId/options", h.CreateOption)
	r.PATCH("/sob/:sobId/dimension/category/:categoryId/option/:optionId", h.UpdateOption)
	r.DELETE("/sob/:sobId/dimension/category/:categoryId/option/:optionId", h.DeleteOption)
}
