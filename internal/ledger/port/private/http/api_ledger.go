package http

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/app"
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

func (h Handler) AllLedgers(c *gin.Context) {
	ls, err := h.app.Queries.ReadLedgers.HandleReadAll(c, c.Param("sob"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, ls)
}

func (h Handler) Migrate(c *gin.Context) {
	if err := h.app.Commands.Migrate.Handle(c); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func InitRouter(h Handler, r *gin.Engine) {
	g1 := r.Group("/private/ledgers")
	{
		g1.POST("/migrate", h.Migrate)
	}

	g2 := r.Group("/private/ledgers/:sob")
	{
		g2.GET("/", h.AllLedgers)
	}
}
