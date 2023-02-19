package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/account/app/query"
	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github/fims-proto/fims-proto-ms/internal/account/app"

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

// ReadPagingAccounts godoc
// @Text List all account configurations
// @Description List all account configurations
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $filter query string false "filter on field(s)" example(title eq 'something' and amount lt 10)
// @Success 200 {array} AccountResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/accounts [get]
func (h Handler) ReadPagingAccounts(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Account], error) {
			return h.app.Queries.PagingAccounts.Handle(c, uuid.MustParse(c.Param("sobId")), pageRequest)
		},
		accountDTOToVO,
	)
}

// ReadPagingLodgersByPeriod godoc
// @Text List accounts in period
// @Description List accounts in period
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param periodId path string true "Period ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $filter query string false "filter on field(s)" example(title eq 'something' and amount lt 10)
// @Success 200 {array} LedgerResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/period/{periodId}/ledgers [get]
func (h Handler) ReadPagingLodgersByPeriod(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Ledger], error) {
			return h.app.Queries.PagingLedgersByPeriod.Handle(c, uuid.MustParse(c.Param("sobId")), uuid.MustParse(c.Param("periodId")), pageRequest)
		},
		ledgerDTOToVO,
	)
}

// ReadPagingPeriods godoc
// @Text List periods
// @Description List periods
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param $page query int false "page number" default(1)
// @Param $size query int false "page size" default(40)
// @Param $sort query string false "sort on field(s)" example(updatedAt desc,createdAt)
// @Param $filter query string false "filter on field(s)" example(title eq 'something' and amount lt 10)
// @Success 200 {array} PeriodResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/periods [get]
func (h Handler) ReadPagingPeriods(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Period], error) {
			return h.app.Queries.PagingPeriods.Handle(c, uuid.MustParse(c.Param("sobId")), pageRequest)
		},
		periodDTOToVO,
	)
}

// ReadSobCurrentPeriod godoc
// @Text Open period in SoB
// @Description Open period in SoB
// @Tags accounts
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Success 200 {object} PeriodResponse
// @Failure 500 {object} Error
// @Router /sob/{sobId}/periods/current [get]
func (h Handler) ReadSobCurrentPeriod(c *gin.Context) {
	periodDTO, err := h.app.Queries.CurrentPeriod.Handle(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		_ = c.Error(err)
		return
	}
	if periodDTO.Id == uuid.Nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, periodDTOToVO(periodDTO))
}

func InitRouter(h Handler, r *gin.RouterGroup) {
	r.GET("/sob/:sobId/accounts", h.ReadPagingAccounts)
	r.GET("/sob/:sobId/periods", h.ReadPagingPeriods)
	r.GET("/sob/:sobId/periods/current", h.ReadSobCurrentPeriod)
	r.GET("/sob/:sobId/period/:periodId/ledgers", h.ReadPagingLodgersByPeriod)
}
