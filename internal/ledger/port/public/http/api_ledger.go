package http

import (
	"github/fims-proto/fims-proto-ms/internal/ledger/app"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/command"
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

// ReadAllAccountingPeriods godoc
// @Summary All accounting periods
// @Description All accounting periods
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Success 200 {array} AccountingPeriodResponse
// @Failure 500 {object} Error
// @Router /periods/{sobId}/ [get]
func (h Handler) ReadAllAccountingPeriods(c *gin.Context) {
	periods, err := h.app.Queries.ReadLedgers.HandleReadAllAccountingPeriods(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	var res []AccountingPeriodResponse
	for _, period := range periods {
		res = append(res, mapFromPeriodQuery(period))
	}
	c.JSON(http.StatusOK, res)
}

// ReadAccountingPeriodById godoc
// @Summary Show accounting period by sob and id
// @Description Show accounting period by sob and id
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param periodId path string true "Accounting Period ID"
// @Success 200 {object} AccountingPeriodResponse
// @Failure 404
// @Failure 500 {object} Error
// @Router /periods/{sobId}/{periodId} [get]
func (h Handler) ReadAccountingPeriodById(c *gin.Context) {
	period, err := h.app.Queries.ReadLedgers.HandleReadAccountingPeriodById(c, uuid.MustParse(c.Param("periodId")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	if period.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, mapFromPeriodQuery(period))
}

// CreateAccountingPeriod godoc
// @Summary Create accounting period
// @Description Create accounting period
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param CreateAccountingPeriodRequest body CreateAccountingPeriodRequest true "Create accounting period request"
// @Success 201 {object} AccountingPeriodResponse
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /periods/{sobId}/ [post]
func (h Handler) CreateAccountingPeriod(c *gin.Context) {
	var req CreateAccountingPeriodRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapErr(err))
		return
	}
	cmd := req.mapToCommand()
	cmd.SobId = uuid.MustParse(c.Param("sobId"))
	createdId, err := h.app.Commands.CreateAccountingPeriod.Handle(c, cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}

	// also create ledgers
	createLedgersCmd := command.CreatePeriodLedgersCmd{
		PeriodId: createdId,
	}
	if err = h.app.Commands.CreatePeriodLedgers.Handle(c, createLedgersCmd); err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}

	createdPeriod, err := h.app.Queries.ReadLedgers.HandleReadAccountingPeriodById(c, createdId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	c.JSON(http.StatusCreated, createdPeriod)
}

// ReadAllLedgersByAccountingPeriod godoc
// @Summary All ledgers in an accounting period
// @Description All ledgers in an accounting period
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param periodId path string true "Accounting Period ID"
// @Success 200 {array} LedgerResponse
// @Failure 500 {object} Error
// @Router /periods/{sobId}/{periodId}/ledgers/ [get]
func (h Handler) ReadAllLedgersByAccountingPeriod(c *gin.Context) {
	ledgers, err := h.app.Queries.ReadLedgers.HandleReadAllLedgersByAccountingPeriod(c, uuid.MustParse(c.Param("periodId")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	var res []LedgerResponse
	for _, ledger := range ledgers {
		res = append(res, mapFromLedgerQuery(ledger))
	}
	c.JSON(http.StatusOK, res)
}

// CalculatePeriodLedgers godoc
// @Summary Calculate ledger balance in accounting period
// @Description Calculate ledger balance in accounting period
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param periodId path string true "Accounting Period ID"
// @Success 204
// @Failure 500 {object} Error
// @Router /periods/{sobId}/{periodId}/ledgers/calculate [post]
func (h Handler) CalculatePeriodLedgers(c *gin.Context) {
	ledgers, err := h.app.Queries.ReadLedgers.HandleReadAllLedgersByAccountingPeriod(c, uuid.MustParse(c.Param("periodId")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapErr(err))
		return
	}
	cmd := command.CalculateLedgerBalanceCmd{
		Ids: make([]uuid.UUID, len(ledgers)),
	}
	for i, ledger := range ledgers {
		cmd.Ids[i] = ledger.Id
	}
	if err = h.app.Commands.CalculateLedgerBalance.Handle(c, cmd); err != nil {
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
	g := r.Group("/periods/:sobId/")
	{
		g.GET("", h.ReadAllAccountingPeriods)
		g.GET(":periodId", h.ReadAccountingPeriodById)
		g.POST("", h.CreateAccountingPeriod)
		g.GET(":periodId/ledgers/", h.ReadAllLedgersByAccountingPeriod)
		g.POST(":periodId/ledgers/calculate", h.CalculatePeriodLedgers)
	}
}
