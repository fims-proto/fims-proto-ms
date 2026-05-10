package http

import (
	"fmt"
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

// PreCloseCheck godoc
//
//	@Text			Validate period against closing conditions
//	@Description	Validate period against closing conditions and return detailed results for each check
//	@Tags			periods
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId		path		string	true	"Sob ID"
//	@Param			periodId	path		string	true	"Period ID"
//	@Success		200			{object}	PreCloseCheckResponse
//	@Failure		500			{object}	Error
//	@Router			/sob/{sobId}/period/{periodId}/pre-close-check [get]
func (h Handler) PreCloseCheck(c *gin.Context) {
	result, err := h.app.Queries.PeriodPreCloseCheck.Handle(
		c,
		uuid.MustParse(c.Param("sobId")),
		uuid.MustParse(c.Param("periodId")),
	)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, preCloseCheckDTOToVO(result))
}

// BatchPreCloseCheck godoc
//
//	@Text			Validate periods against closing conditions for batch close
//	@Description	Validate from the current period to the target period. Checks for unposted journals and trial balance in each existing period. P&L and CYP checks are omitted because auto-transfer handles them during batch close.
//	@Tags			periods
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId			path		string	true	"Sob ID"
//	@Param			targetPeriod	query		string	true	"Target period in YYYY-MM format"
//	@Success		200				{object}	BatchPreCloseCheckResponse
//	@Failure		400				{object}	Error
//	@Failure		500				{object}	Error
//	@Router			/sob/{sobId}/periods/batch-pre-close-check [get]
func (h Handler) BatchPreCloseCheck(c *gin.Context) {
	var targetYear, targetMonth int
	if _, err := fmt.Sscanf(c.Query("targetPeriod"), "%d-%d", &targetYear, &targetMonth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid targetPeriod format, expected YYYY-MM"})
		return
	}

	result, err := h.app.Queries.BatchPeriodPreCloseCheck.Handle(
		c,
		uuid.MustParse(c.Param("sobId")),
		targetYear,
		targetMonth,
	)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.JSON(http.StatusOK, batchPreCloseCheckDTOToVO(result))
}

// ClosePeriods godoc
//
//	@Text			Batch close periods to a target period
//	@Description	Sequentially close all periods from the current period to the target period (inclusive) in a single atomic transaction. For each period, automatically creates monthly and year-end closing journals before closing. Rolls back entirely if any period fails validation.
//	@Tags			periods
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId			path	string	true	"Sob ID"
//	@Param			targetPeriod	query	string	true	"Target period in YYYY-MM format"
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/periods/batch-close [post]
func (h Handler) ClosePeriods(c *gin.Context) {
	var targetYear, targetMonth int
	if _, err := fmt.Sscanf(c.Query("targetPeriod"), "%d-%d", &targetYear, &targetMonth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid targetPeriod format, expected YYYY-MM"})
		return
	}

	if err := h.app.Commands.ClosePeriods.Handle(c, command.ClosePeriodsCmd{
		SobId:       uuid.MustParse(c.Param("sobId")),
		TargetYear:  targetYear,
		TargetMonth: targetMonth,
	}); err != nil {
		_ = c.Error(err)
		return
	}
	c.Status(http.StatusNoContent)
}
