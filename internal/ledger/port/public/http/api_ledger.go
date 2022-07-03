package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github/fims-proto/fims-proto-ms/internal/ledger/app"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"

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

// ReadCurrentPeriod godoc
// @Summary Current period
// @Description Current period
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Success 200 {object} PeriodResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/period/current [get]
func (h Handler) ReadCurrentPeriod(c *gin.Context) {
	period, err := h.app.Queries.ReadLedgers.HandleReadOpenPeriod(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if period.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, mapFromPeriodQuery(period))
}

// ReadAllPeriods godoc
// @Summary All periods
// @Description All periods
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $choose query string false "choose only field(s)"
// @Param $filter query string false "filter on field(s)" example(title eq 'some thing' and amount lt 10)
// @Success 200 {array} PeriodResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/periods/ [get]
func (h Handler) ReadAllPeriods(c *gin.Context) {
	pageable, err := data.NewPageableFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	periodsPage, err := h.app.Queries.ReadLedgers.HandleReadAllPeriods(c, uuid.MustParse(c.Param("sobId")), pageable)
	if err != nil {
		_ = c.Error(err)
		return
	}
	periods := make([]PeriodResponse, len(periodsPage.Content))
	for i, period := range periodsPage.Content {
		periods[i] = mapFromPeriodQuery(period)
	}
	resp, _ := data.NewPage(periods, periodsPage.Page, periodsPage.Size, periodsPage.NumberOfElements)
	c.JSON(http.StatusOK, resp)
}

// ReadPeriodById godoc
// @Summary Show period by sob and id
// @Description Show period by sob and id
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param periodId path string true "Period ID"
// @Success 200 {object} PeriodResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /sob/{sobId}/period/{periodId} [get]
func (h Handler) ReadPeriodById(c *gin.Context) {
	period, err := h.app.Queries.ReadLedgers.HandleReadPeriodById(c, uuid.MustParse(c.Param("periodId")))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if period.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, mapFromPeriodQuery(period))
}

// CreatePeriod godoc
// @Summary Create period
// @Description Create period
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param CreatePeriodRequest body CreatePeriodRequest true "Create period request"
// @Success 201 {object} PeriodResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/periods/ [post]
func (h Handler) CreatePeriod(c *gin.Context) {
	var req CreatePeriodRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	cmd := req.mapToCommand()
	cmd.SobId = uuid.MustParse(c.Param("sobId"))
	createdId, err := h.app.Commands.CreatePeriod.Handle(c, cmd)
	if err != nil {
		_ = c.Error(err)
		return
	}

	// also create ledgers
	createLedgersCmd := command.CreatePeriodLedgersCmd{
		PeriodId: createdId,
	}
	if err = h.app.Commands.CreatePeriodLedgers.Handle(c, createLedgersCmd); err != nil {
		_ = c.Error(err)
		return
	}

	createdPeriod, err := h.app.Queries.ReadLedgers.HandleReadPeriodById(c, createdId)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, createdPeriod)
}

// ReadAllLedgersByPeriod godoc
// @Summary All ledgers in an period
// @Description All ledgers in an period
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param periodId path string true "Period ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $choose query string false "choose only field(s)"
// @Param $filter query string false "filter on field(s)" example(title eq 'some thing' and amount lt 10)
// @Success 200 {array} LedgerResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/period/{periodId}/ledgers/ [get]
func (h Handler) ReadAllLedgersByPeriod(c *gin.Context) {
	pageable, err := data.NewPageableFromRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	ledgersPage, err := h.app.Queries.ReadLedgers.HandleReadAllLedgersByPeriod(c, uuid.MustParse(c.Param("periodId")), pageable)
	if err != nil {
		_ = c.Error(err)
		return
	}
	ledgers := make([]LedgerResponse, len(ledgersPage.Content))
	for i, ledger := range ledgersPage.Content {
		ledgers[i] = mapFromLedgerQuery(ledger)
	}
	resp, _ := data.NewPage(ledgers, ledgersPage.Page, ledgersPage.Size, ledgersPage.NumberOfElements)
	c.JSON(http.StatusOK, resp)
}

// CalculatePeriodLedgers godoc
// @Summary Calculate ledger balance in period
// @Description Calculate ledger balance in period
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param periodId path string true "Period ID"
// @Success 204
// @Failure 500 {object} Error
// @Router /sob/{sobId}/period/{periodId}/ledgers/calculate [post]
func (h Handler) CalculatePeriodLedgers(c *gin.Context) {
	periodId := uuid.MustParse(c.Param("periodId"))
	cmd := command.CalculateBalanceByPeriodCmd{
		PeriodId: periodId,
	}
	if err := h.app.Commands.CalculateLedgerBalance.HandleCalculateBalanceByPeriod(c, cmd); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/sob/:sobId/period/current", h.ReadCurrentPeriod)
	r.GET("/sob/:sobId/periods/", h.ReadAllPeriods)
	r.GET("/sob/:sobId/period/:periodId", h.ReadPeriodById)
	r.POST("/sob/:sobId/periods/", h.CreatePeriod)
	r.GET("/sob/:sobId/period/:periodId/ledgers/", h.ReadAllLedgersByPeriod)
	r.POST("/sob/:sobId/period/rperiodId/ledgers/calculate", h.CalculatePeriodLedgers)
}
