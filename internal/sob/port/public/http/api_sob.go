package http

import (
	"github/fims-proto/fims-proto-ms/internal/sob/app"
	"github/fims-proto/fims-proto-ms/internal/sob/app/command"
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

// ReadAllSobs godoc
// @Summary List all sobs
// @Description List all sobs
// @Tags sobs
// @Accept application/json
// @Produce application/json
// @Success 200 {array} SobResponse
// @Failure 500 {object} Error
// @Router /sobs/ [get]
func (h Handler) ReadAllSobs(c *gin.Context) {
	sobs, err := h.app.Queries.ReadSobs.HandleReadAll(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	resp := make([]SobResponse, len(sobs))
	for i, sob := range sobs {
		resp[i] = mapFromSobQuery(sob)
	}
	c.JSON(http.StatusOK, resp)
}

// ReadSobById godoc
// @Summary Show sob by id
// @Description Show sob by id
// @Tags sobs
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "ID of a SobId"
// @Success 200 {object} SobResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /sobs/{sobId} [get]
func (h Handler) ReadSobById(c *gin.Context) {
	sob, err := h.app.Queries.ReadSobs.HandleReadById(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	if sob.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, mapFromSobQuery(sob))
}

// CreateSob godoc
// @Summary Create sob
// @Description Create sob
// @Tags sobs
// @Accept application/json
// @Produce application/json
// @Param CreateSobRequest body CreateSobRequest true "CreateSob SobId"
// @Success 201 {object} SobResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sobs/ [post]
func (h Handler) CreateSob(c *gin.Context) {
	var req CreateSobRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	createdId, err := h.app.Commands.CreateSob.Handle(c, req.mapToCommand())
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	createdSob, err := h.app.Queries.ReadSobs.HandleReadById(c, createdId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.JSON(http.StatusCreated, mapFromSobQuery(createdSob))
}

// UpdateSob godoc
// @Summary Update sob
// @Description Update sob
// @Tags sobs
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "SobId ID"
// @Param UpdateSobRequest body UpdateSobRequest true "UpdateLineItems sob request"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sobs/{sobId} [patch]
func (h Handler) UpdateSob(c *gin.Context) {
	var req UpdateSobRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	cmd := command.UpdateSobCmd{
		Id:                 uuid.MustParse(c.Param("sobId")),
		Name:               req.Name,
		AccountsCodeLength: req.AccountsCodeLength,
	}
	if err := h.app.Commands.UpdateSob.Handle(c, cmd); err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.Status(http.StatusNoContent)
}

func wrapErr(e error) Error {
	var slug string
	se, ok := e.(slugErr)
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
	g := r.Group("/sobs/")
	{
		g.GET("", h.ReadAllSobs)
		g.GET(":sobId", h.ReadSobById)
		g.PATCH(":sobId", h.UpdateSob)
		g.POST("", h.CreateSob)
	}
}
