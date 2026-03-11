package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadLedgerSummary godoc
//
//	@Text			Get ledger summary by account across period range
//	@Description	Get aggregated ledger summary for a single account across a period range
//	@Tags			ledgers
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId		path		string	true	"Sob ID"
//	@Param			accountId	path		string	true	"Account ID"
//	@Param			fromPeriod	query		string	true	"From period (YYYY-MM)"
//	@Param			toPeriod	query		string	true	"To period (YYYY-MM)"
//	@Success		200			{object}	LedgerSummaryResponse
//	@Failure		400			{object}	Error
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/ledger/{accountId} [get]
func (h Handler) ReadLedgerSummary(c *gin.Context) {
	summary, err := h.app.Queries.LedgerSummary.Handle(
		c,
		uuid.MustParse(c.Param("sobId")),
		uuid.MustParse(c.Param("accountId")),
		c.Query("fromPeriod"),
		c.Query("toPeriod"),
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, ledgerSummaryToVO(summary))
}

// ReadLedgersByPeriodRange godoc
//
//	@Text			List ledgers aggregated across a period range
//	@Description	List all ledgers for a SoB aggregated across a period range. Returns one entry per account with opening amount from the first period, summed period debit/credit/amount, and ending amount from the last period.
//	@Tags			ledgers
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId		path		string	true	"Sob ID"
//	@Param			fromPeriod	query		string	true	"From period (YYYY-MM)"
//	@Param			toPeriod	query		string	true	"To period (YYYY-MM)"
//	@Success		200			{array}		LedgerResponse
//	@Failure		400			{object}	Error
//	@Failure		500			{object}	Error
//	@Router			/sob/{sobId}/ledgers [get]
func (h Handler) ReadLedgersByPeriodRange(c *gin.Context) {
	ledgers, err := h.app.Queries.PagingLedgersByPeriod.Handle(
		c,
		uuid.MustParse(c.Param("sobId")),
		c.Query("fromPeriod"),
		c.Query("toPeriod"),
	)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, converter.DTOsToVOs(ledgers, ledgerDTOToVO))
}

// ReadFirstPeriodLedgers godoc
//
//	@Text			List ledgers in first period
//	@Description	List ledgers in first period
//	@Tags			ledgers
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId	path		string	true	"Sob ID"
//	@Success		200		{object}	PeriodAndLedgersResponse
//	@Failure		500		{object}	Error
//	@Router			/sob/{sobId}/first-period/ledgers [get]
func (h Handler) ReadFirstPeriodLedgers(c *gin.Context) {
	period, ledgers, err := h.app.Queries.FirstPeriodLedgers.Handle(c, uuid.MustParse(c.Param("sobId")))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, PeriodAndLedgersResponse{
		Period:  periodDTOToVO(period),
		Ledgers: converter.DTOsToVOs(ledgers, ledgerDTOToVO),
	})
}

// InitializeLedgers godoc
//
//	@Text			Initialize ledgers in first period of current SoB
//	@Description	Initialize ledgers in first period of current SoB
//	@Tags			ledgers
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId							path	string							true	"Sob ID"
//	@Param			InitializeLedgersBalanceRequest	body	InitializeLedgersBalanceRequest	true	"Ledgers with opening balance"
//	@Success		204
//	@Failure		400	{object}	Error
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/ledgers/initialize [post]
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

// ReadAuxiliaryLedgerSummary godoc
//
//	@Text			Get auxiliary ledger summary by account across period range
//	@Description	Get aggregated auxiliary ledger summary grouped by auxiliary account for a specific category
//	@Tags			ledgers
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId		path		string	true	"Sob ID"
//	@Param			accountId	path		string	true	"Account ID"
//	@Param			categoryKey	query		string	true	"Category key (e.g., customer, project)"
//	@Param			fromPeriod	query		string	true	"From period (YYYY-MM)"
//	@Param			toPeriod	query		string	true	"To period (YYYY-MM)"
//	@Param			$page		query		int		false	"page number"			default(1)
//	@Param			$size		query		int		false	"page size"				default(40)
//	@Param			$sort		query		string	false	"sort on field(s)"		example(updatedAt desc,createdAt)
//	@Param			$filter		query		string	false	"filter on field(s)"	example(title eq 'something' and amount lt 10)
//	@Success		200			{object}	data.PageResponse[AuxiliaryLedgerSummaryResponse]
//	@Failure		400			{object}	Error
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/ledger/{accountId}/auxiliary [get]
func (h Handler) ReadAuxiliaryLedgerSummary(c *gin.Context) {
	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.AuxiliaryLedgerSummary], error) {
			return h.app.Queries.AuxiliaryLedgerSummary.Handle(
				c,
				uuid.MustParse(c.Param("sobId")),
				uuid.MustParse(c.Param("accountId")),
				c.Query("categoryKey"),
				c.Query("fromPeriod"),
				c.Query("toPeriod"),
				pageRequest,
			)
		},
		auxiliaryLedgerSummaryToVO,
	)
}

// ReadLedgerEntries godoc
//
//	@Text			Get ledger entries by account across period range
//	@Description	Get detailed ledger entries for a single account across a period range
//	@Tags			ledgers
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId				path		string	true	"Sob ID"
//	@Param			accountId			path		string	true	"Account ID"
//	@Param			fromPeriod			query		string	true	"From period (YYYY-MM)"
//	@Param			toPeriod			query		string	true	"To period (YYYY-MM)"
//	@Param			auxiliaryAccountId	query		string	false	"Auxiliary Account ID (optional)"
//	@Param			$page				query		int		false	"page number"			default(1)
//	@Param			$size				query		int		false	"page size"				default(40)
//	@Param			$sort				query		string	false	"sort on field(s)"		example(updatedAt desc,createdAt)
//	@Param			$filter				query		string	false	"filter on field(s)"	example(text eq 'something' and amount lt 10)
//	@Success		200					{object}	data.PageResponse[LedgerEntryResponse]
//	@Failure		400					{object}	Error
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/ledger/{accountId}/entries [get]
func (h Handler) ReadLedgerEntries(c *gin.Context) {
	// Parse optional auxiliaryAccountId
	var auxiliaryAccountId *uuid.UUID
	if auxiliaryAccountIdStr := c.Query("auxiliaryAccountId"); auxiliaryAccountIdStr != "" {
		id := uuid.MustParse(auxiliaryAccountIdStr)
		auxiliaryAccountId = &id
	}

	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.LedgerEntry], error) {
			return h.app.Queries.PagingLedgerEntries.Handle(
				c,
				uuid.MustParse(c.Param("sobId")),
				uuid.MustParse(c.Param("accountId")),
				auxiliaryAccountId,
				c.Query("fromPeriod"),
				c.Query("toPeriod"),
				pageRequest,
			)
		},
		ledgerEntryDTOToVO,
	)
}
