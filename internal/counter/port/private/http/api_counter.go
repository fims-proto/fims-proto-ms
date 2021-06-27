package http

import (
	"github/fims-proto/fims-proto-ms/internal/counter/app"
	"net/http"

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

func (h Handler) Dataload(c *gin.Context) {
	if err := h.app.Commands.LoadCounters.Handle(c.Request.Context(), c.Param("sob")); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func InitRouter(h Handler, r *gin.Engine) {
	g := r.Group("/private/counters/:sob")
	{
		g.POST("/dataload", h.Dataload)
	}
}
