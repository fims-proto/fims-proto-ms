package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadCashFlowItems godoc
//
//	@Text			List cash flow items
//	@Description	List all cash flow items for a Set of Books
//	@Tags			cash-flow-items
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Success		200		{array}		CashFlowItemResponse
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/cash-flow-items [get]
func (h Handler) ReadCashFlowItems(c *gin.Context) {
	items, err := h.app.Queries.CashFlowItems.Handle(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, converter.DTOsToVOs(items, cashFlowItemDTOToVO))
}
