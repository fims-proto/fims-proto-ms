package http

import (
	"fmt"
	"github/fims-proto/fims-proto-ms/internal/sob/app"
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

func (h Handler) AllSobs(c *gin.Context) {
	sobs, err := h.app.Queries.ReadSobs.HandleReadAll(c.Request.Context())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	resp := SobsResponse{}
	for _, sob := range sobs {
		resp = append(resp, mapFromSobQuery(sob))
	}
	c.JSON(http.StatusOK, resp)
}

func (h Handler) SobById(c *gin.Context) {
	sob, err := h.app.Queries.ReadSobs.HandleReadById(c.Request.Context(), c.Param("sob"))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, mapFromSobQuery(sob))
}

func (h Handler) Create(c *gin.Context) {
	var req CreateSobRequest
	if err := c.ShouldBind(&req); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if err := h.app.Commands.CreateSob.Handle(c.Request.Context(), req.mapToCommand()); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusCreated)
	// TODO figure out a way to avoid hard coded string of url
	c.Writer.Header().Set("Content-Location", fmt.Sprintf("/sobs/%s", req.Id))
}

func InitRouter(h Handler, r *gin.Engine) {
	g := r.Group("/sobs")
	{
		g.GET("/", h.AllSobs)
		g.GET("/:sob", h.SobById)
		g.POST("/", h.Create)
	}
}
