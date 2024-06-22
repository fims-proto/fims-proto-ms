package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
)

// ReadPagingLedgersByPeriod godoc
// @Text List ledgers in period
// @Description List ledgers in period
// @Tags ledgers
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
func (h Handler) ReadPagingLedgersByPeriod(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Ledger], error) {
			return h.app.Queries.PagingLedgersByPeriod.Handle(c, uuid.MustParse(c.Param("sobId")), uuid.MustParse(c.Param("periodId")), pageRequest)
		},
		ledgerDTOToVO,
	)
}

// InitializeLedgers godoc
// @Text Initialize ledgers in first period of current SoB
// @Description Initialize ledgers in first period of current SoB
// @Tags ledgers
// @Accept application/json
// @Produce application/json
// @Param sobId path string true "Sob ID"
// @Param InitializeLedgersBalanceRequest body InitializeLedgersBalanceRequest true "Ledgers with opening balance"
// @Success 204
// @Failure 400 {object} Error
// @Failure 500 {object} Error
// @Router /sob/{sobId}/ledgers/initialize [post]
func (h Handler) InitializeLedgers(c *gin.Context) {
	var req InitializeLedgersBalanceRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := h.app.Commands.InitializeLedgersBalance.Handle(c, req.mapToCommand(uuid.MustParse(c.Param("sobId")))); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
