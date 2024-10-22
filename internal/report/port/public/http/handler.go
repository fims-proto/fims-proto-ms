package http

import (
	"github.com/gin-gonic/gin"
	"github/fims-proto/fims-proto-ms/internal/report/app"
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
	r.GET("/sob/:sobId/reports", h.ReadAllReports)
	r.GET("/sob/:sobId/report/:reportId", h.ReadReportById)
	r.POST("/sob/:sobId/report/:reportId/generate", h.GenerateReport)
	r.POST("/sob/:sobId/report/:reportId/regenerate", h.RegenerateReport)
}
