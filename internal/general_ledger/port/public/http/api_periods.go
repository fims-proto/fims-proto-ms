package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadPeriods godoc
//
//	@Text			List all periods
//	@Description	List all periods
//	@Tags			periods
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Success		200		{array}		PeriodResponse
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/periods [get]
func (h Handler) ReadPeriods(c *gin.Context) {
	periods, err := h.app.Queries.AllPeriods.Handle(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, converter.DTOsToVOs(periods, periodDTOToVO))
}

// ClosePeriod godoc
//
//	@Text			Close period
//	@Description	Close period
//	@Tags			periods
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId		path	string	true	"Sob ID"
//	@Param			periodId	path	string	true	"Period ID"
//	@Success		204
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/period/{periodId}/close [post]
func (h Handler) ClosePeriod(c *gin.Context) {
	if err := h.app.Commands.ClosePeriod.Handle(c, command.ClosePeriodCmd{
		SobId:    uuid.MustParse(c.Param("sobId")),
		PeriodId: uuid.MustParse(c.Param("periodId")),
	}); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
