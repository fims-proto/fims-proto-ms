package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadPagingPeriods godoc
//
//	@Text			List periods
//	@Description	List periods
//	@Tags			periods
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Param			$page	query		int		false	"page number"			default(1)
//	@Param			$size	query		int		false	"page size"				default(40)
//	@Param			$sort	query		string	false	"sort on field(s)"		example(updatedAt desc,createdAt)
//	@Param			$filter	query		string	false	"filter on field(s)"	example(title eq 'something' and amount lt 10)
//	@Success		200		{object}	data.PageResponse[PeriodResponse]
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/periods [get]
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
//
//	@Text			Current period
//	@Description	Current period
//	@Tags			periods
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Success		200		{object}	PeriodResponse
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/periods/current [get]
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
