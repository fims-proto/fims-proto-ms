package http

import (
	"github/fims-proto/fims-proto-ms/internal/account/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h Handler) DataLoad(c *gin.Context) {
	if err := h.app.Commands.LoadAccounts.Handle(c, uuid.MustParse(c.Param("sobId"))); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) Migrate(c *gin.Context) {
	if err := h.app.Commands.Migrate.Handle(c); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	g1 := r.Group("/accounts/")
	{
		g1.POST("migrate", h.Migrate)
	}

	g2 := r.Group("/accounts/:sobId/")
	{
		g2.POST("dataload", h.DataLoad)
	}
}
