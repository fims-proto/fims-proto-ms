package http

import (
	"fmt"
	"github/fims-proto/fims-proto-ms/internal/sob/app"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"
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

// AllSobs godoc
// @Summary List all sobs
// @Description List all sobs
// @Tags sobs
// @Accept application/json
// @Produce application/json
// @Success 200 {array} SobResponse
// @Failure 500 {object} Error
// @Router /sobs [get]
func (h Handler) AllSobs(c *gin.Context) {
	sobs, err := h.app.Queries.ReadSobs.HandleReadAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	resp := []SobResponse{}
	for _, sob := range sobs {
		resp = append(resp, mapFromSobQuery(sob))
	}
	c.JSON(http.StatusOK, resp)
}

// SobById godoc
// @Summary Show sob by id
// @Description Show sob by id
// @Tags sobs
// @Accept application/json
// @Produce application/json
// @Param sob path string true "Id of a Sob"
// @Success 200 {object} SobResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /sobs/{sob} [get]
func (h Handler) SobById(c *gin.Context) {
	sob, err := h.app.Queries.ReadSobs.HandleReadById(c, c.Param("sob"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	if sob == (query.Sob{}) {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, mapFromSobQuery(sob))
}

// Create godoc
// @Summary Create sob
// @Description Create sob
// @Tags sobs
// @Accept application/json
// @Produce application/json
// @Param sob body CreateSobRequest true "Create Sob"
// @Success 201
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sobs [post]
func (h Handler) Create(c *gin.Context) {
	var req CreateSobRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	if err := h.app.Commands.CreateSob.Handle(c, req.mapToCommand()); err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.Status(http.StatusCreated)
	// TODO figure out a way to avoid hard coded string of url
	c.Writer.Header().Set("Content-Location", fmt.Sprintf("/sobs/%s", req.Id))
}

func wrapErr(e error) Error {
	var slug string
	se, ok := e.(sluggableErr)
	if ok {
		slug = se.Slug()
	} else {
		slug = "unknown-error"
	}
	return Error{
		Slug:    slug,
		Message: e.Error(),
	}
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	g := r.Group("/sobs")
	{
		g.GET("/", h.AllSobs)
		g.GET("/:sob", h.SobById)
		g.POST("/", h.Create)
	}
}
