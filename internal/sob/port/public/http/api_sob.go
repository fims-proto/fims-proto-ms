package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/sob/app/query"

	"github/fims-proto/fims-proto-ms/internal/sob/app"
	"github/fims-proto/fims-proto-ms/internal/sob/app/command"

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
// @Text List all sobs
// @Description List all sobs
// @Tags sobs
// @Accept application/json
// @Produce application/json
// @Success 200 {array} SobResponse
// @Failure 500 {object} Error
// @Router /sobs/ [get]
func (h Handler) ReadAllSobs(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Sob], error) {
			return h.app.Queries.PagingSobs.Handle(c, pageRequest)
		},
		sobDTOToVO,
	)
}

// ReadSobById godoc
// @Text Show sob by id
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
	sob, err := h.app.Queries.SobById.Handle(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if sob.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, sobDTOToVO(sob))
}

// CreateSob godoc
// @Text Create sob
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
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := req.mapToCommand()
	err := h.app.Commands.CreateSob.Handle(c, cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}
	createdSob, err := h.app.Queries.SobById.Handle(c, cmd.SobId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, sobDTOToVO(createdSob))
}

// UpdateSob godoc
// @Text Update sob
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
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := command.UpdateSobCmd{
		SobId:              uuid.MustParse(c.Param("sobId")),
		Name:               req.Name,
		AccountsCodeLength: req.AccountsCodeLength,
	}
	if err := h.app.Commands.UpdateSob.Handle(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/sobs/", h.ReadAllSobs)
	r.GET("/sobs/:sobId", h.ReadSobById)
	r.PATCH("/sobs/:sobId", h.UpdateSob)
	r.POST("/sobs/", h.CreateSob)
}
