package http

import (
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadLedgersByPeriodRange godoc
//
//	@Text			List ledgers aggregated across a period range
//	@Description	List all ledgers for a SoB aggregated across a period range. Returns one entry per account with opening amount from the first period, summed period debit/credit/amount, and ending amount from the last period. When dimensionOptionId is provided, only accounts that have journal lines tagged with that dimension option are returned.
//	@Tags			ledgers
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId				path		string	true	"Sob ID"
//	@Param			fromPeriod			query		string	true	"From period (YYYY-MM)"
//	@Param			toPeriod			query		string	true	"To period (YYYY-MM)"
//	@Param			dimensionOptionId	query		string	false	"Dimension Option ID (optional filter)"
//	@Param			$page				query		int		false	"page number"	default(1)
//	@Param			$size				query		int		false	"page size"		default(40)
//	@Success		200					{object}	data.PageResponse[LedgerResponse]
//	@Failure		400					{object}	Error
//	@Failure		500					{object}	Error
//	@Router			/sob/{sobId}/ledgers [get]
func (h Handler) ReadLedgersByPeriodRange(c *gin.Context) {
	var dimensionOptionId *uuid.UUID
	if raw := c.Query("dimensionOptionId"); raw != "" {
		dimensionOptionId = new(uuid.MustParse(raw))
	}

	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.Ledger], error) {
			return h.app.Queries.PagingLedgersByPeriod.Handle(
				c,
				uuid.MustParse(c.Param("sobId")),
				c.Query("fromPeriod"),
				c.Query("toPeriod"),
				dimensionOptionId,
				pageRequest,
			)
		},
		ledgerDTOToVO,
	)
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

// ReadLedgerTransactions godoc
//
//	@Text			Get ledger transaction entries across period range
//	@Description	Get detailed ledger transaction entries across a period range. At least one of accountId or dimensionOptionId must be provided. When both are provided, entries matching both filters are returned.
//	@Tags			ledgers
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId				path		string	true	"Sob ID"
//	@Param			fromPeriod			query		string	true	"From period (YYYY-MM)"
//	@Param			toPeriod			query		string	true	"To period (YYYY-MM)"
//	@Param			accountId			query		string	false	"Account ID (optional filter — must provide at least one of accountId or dimensionOptionId)"
//	@Param			dimensionOptionId	query		string	false	"Dimension Option ID (optional filter — must provide at least one of accountId or dimensionOptionId)"
//	@Param			$page				query		int		false	"page number"			default(1)
//	@Param			$size				query		int		false	"page size"				default(40)
//	@Param			$sort				query		string	false	"sort on field(s)"		example(updatedAt desc,createdAt)
//	@Param			$filter				query		string	false	"filter on field(s)"	example(text eq 'something' and amount lt 10)
//	@Success		200					{object}	data.PageResponse[LedgerEntryResponse]
//	@Failure		400					{object}	Error
//	@Failure		404
//	@Failure		500	{object}	Error
//	@Router			/sob/{sobId}/ledgers/transactions [get]
func (h Handler) ReadLedgerTransactions(c *gin.Context) {
	var accountId *uuid.UUID
	if raw := c.Query("accountId"); raw != "" {
		accountId = new(uuid.MustParse(raw))
	}

	var dimensionOptionId *uuid.UUID
	if raw := c.Query("dimensionOptionId"); raw != "" {
		dimensionOptionId = new(uuid.MustParse(raw))
	}

	if accountId == nil && dimensionOptionId == nil {
		_ = c.Error(commonErrors.NewInvalidInputError(commonErrors.SlugLedgerTransactionsMissingFilter))
		return
	}

	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.LedgerEntry], error) {
			return h.app.Queries.LedgerEntries.Handle(
				c,
				uuid.MustParse(c.Param("sobId")),
				accountId,
				c.Query("fromPeriod"),
				c.Query("toPeriod"),
				dimensionOptionId,
				pageRequest,
			)
		},
		ledgerEntryDTOToVO,
	)
}

// ReadLedgerByDimensionCategory godoc
//
//	@Text			Get ledger amounts aggregated by dimension option
//	@Description	Get total amounts from journal lines for a specific dimension category across a period range, grouped by dimension option. When accountId is provided, results are scoped to that account only; otherwise all accounts are aggregated.
//	@Tags			ledgers
//	@Accept			application/json
//	@Produce		application/json
//	@Param			sobId				path		string	true	"Sob ID"
//	@Param			dimensionCategoryId	path		string	true	"Dimension Category ID"
//	@Param			accountId			query		string	false	"Account ID (optional — omit to aggregate all accounts)"
//	@Param			fromPeriod			query		string	true	"From period (YYYY-MM)"
//	@Param			toPeriod			query		string	true	"To period (YYYY-MM)"
//	@Param			$page				query		int		false	"page number"	default(1)
//	@Param			$size				query		int		false	"page size"		default(40)
//	@Success		200					{object}	data.PageResponse[LedgerDimensionOptionResponse]
//	@Failure		400					{object}	Error
//	@Failure		500					{object}	Error
//	@Router			/sob/{sobId}/ledgers/dimension-category/{dimensionCategoryId}/options [get]
func (h Handler) ReadLedgerByDimensionCategory(c *gin.Context) {
	var accountId *uuid.UUID
	if raw := c.Query("accountId"); raw != "" {
		accountId = new(uuid.MustParse(raw))
	}

	data.PagingResponseProcessor(
		c,
		func(pageRequest data.PageRequest) (data.Page[query.LedgerDimensionSummaryItem], error) {
			return h.app.Queries.LedgersByDimensionCategory.Handle(
				c,
				uuid.MustParse(c.Param("sobId")),
				uuid.MustParse(c.Param("dimensionCategoryId")),
				accountId,
				c.Query("fromPeriod"),
				c.Query("toPeriod"),
				pageRequest,
			)
		},
		ledgerDimensionSummaryItemToVO,
	)
}
