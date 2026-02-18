package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/report/app"
	"github/fims-proto/fims-proto-ms/internal/report/app/command"

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

func (h Handler) Migrate(c *gin.Context) {
	if err := h.app.Commands.Migrate.Handle(c); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

// Initialize TODO internal test
func (h Handler) Initialize(c *gin.Context) {
	if err := h.app.Commands.Initialize.Handle(c, command.InitializeCmd{SobId: uuid.MustParse(c.Query("sobId"))}); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	g := r.Group("/reports/")
	{
		g.POST("migrate", h.Migrate)
		g.POST("init", h.Initialize)
	}
}
