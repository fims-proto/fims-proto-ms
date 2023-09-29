package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
)

// ReadPagingAuxiliaryLedgers godoc
// @Text List all auxiliary ledgers
// @Description List all auxiliary ledgers
// @Tags auxiliary ledgers
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
// @Router /sob/{sobId}/period/{periodId}/auxiliary-ledgers [get]
func (h Handler) ReadPagingAuxiliaryLedgers(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.AuxiliaryLedger], error) {
			return h.app.Queries.PagingAuxiliaryLedgers.Handle(c, uuid.MustParse(c.Param("sobId")), uuid.MustParse(c.Param("periodId")), pageRequest)
		},
		auxiliaryLedgerDTOToVO,
	)
}
