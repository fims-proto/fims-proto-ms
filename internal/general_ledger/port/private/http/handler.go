package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app"

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

func (h Handler) Migrate(c *gin.Context) {
	if err := h.app.Commands.Migrate.Handle(c); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	g := r.Group("/general-ledgers/")
	{
		g.POST("migrate", h.Migrate)
	}
}
