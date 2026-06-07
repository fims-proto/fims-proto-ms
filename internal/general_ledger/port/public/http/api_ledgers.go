package http

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ReadLedgersByPeriodRangeInput struct {
	SobId             uuid.UUID                     `path:"sobId"`
	FromPeriod        string                        `query:"fromPeriod" doc:"From period (YYYY-MM)"`
	ToPeriod          string                        `query:"toPeriod" doc:"To period (YYYY-MM)"`
	DimensionOptionId data.OptionalParam[uuid.UUID] `query:"dimensionOptionId" doc:"Optional dimension option filter"`
	data.PaginationInput
}

type ReadLedgersByPeriodRangeOutput struct {
	Body data.PageResponse[LedgerResponse]
}

// ReadLedgersByPeriodRange lists ledger balances across a period range.
func (h Handler) ReadLedgersByPeriodRange(ctx context.Context, input *ReadLedgersByPeriodRangeInput) (*ReadLedgersByPeriodRangeOutput, error) {
	pageRequest, err := data.PageRequestFromInput(input.PaginationInput)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	page, err := h.app.Queries.PagingLedgersByPeriod.Handle(ctx, input.SobId, input.FromPeriod, input.ToPeriod, input.DimensionOptionId.Ptr(), pageRequest)
	if err != nil {
		return nil, err
	}
	return &ReadLedgersByPeriodRangeOutput{Body: data.MapPageResponse(page, ledgerDTOToVO)}, nil
}

type ReadFirstPeriodLedgersInput struct {
	SobId uuid.UUID `path:"sobId"`
}

type ReadFirstPeriodLedgersOutput struct {
	Body PeriodAndLedgersResponse
}

// ReadFirstPeriodLedgers returns opening ledgers for first period.
func (h Handler) ReadFirstPeriodLedgers(ctx context.Context, input *ReadFirstPeriodLedgersInput) (*ReadFirstPeriodLedgersOutput, error) {
	period, ledgers, err := h.app.Queries.FirstPeriodLedgers.Handle(ctx, input.SobId)
	if err != nil {
		return nil, err
	}
	vos := make([]LedgerResponse, len(ledgers))
	for i, l := range ledgers {
		vos[i] = ledgerDTOToVO(l)
	}
	return &ReadFirstPeriodLedgersOutput{Body: PeriodAndLedgersResponse{
		Period:  periodDTOToVO(period),
		Ledgers: vos,
	}}, nil
}

type InitializeLedgersInput struct {
	SobId uuid.UUID `path:"sobId"`
	Body  InitializeLedgersBalanceRequest
}

// InitializeLedgers sets initial ledger balances for a SoB.
func (h Handler) InitializeLedgers(ctx context.Context, input *InitializeLedgersInput) (*struct{}, error) {
	if err := h.app.Commands.InitializeLedgersBalance.Handle(ctx, input.Body.mapToCommand(input.SobId)); err != nil {
		return nil, err
	}
	return nil, nil
}

type ReadLedgerTransactionsInput struct {
	SobId             uuid.UUID                     `path:"sobId"`
	FromPeriod        string                        `query:"fromPeriod" doc:"From period (YYYY-MM)"`
	ToPeriod          string                        `query:"toPeriod" doc:"To period (YYYY-MM)"`
	AccountId         data.OptionalParam[uuid.UUID] `query:"accountId" doc:"Optional account filter"`
	DimensionOptionId data.OptionalParam[uuid.UUID] `query:"dimensionOptionId" doc:"Optional dimension option filter"`
	data.PaginationInput
}

type ReadLedgerTransactionsOutput struct {
	Body data.PageResponse[LedgerEntryResponse]
}

// ReadLedgerTransactions lists ledger entries by account or dimension filter.
func (h Handler) ReadLedgerTransactions(ctx context.Context, input *ReadLedgerTransactionsInput) (*ReadLedgerTransactionsOutput, error) {
	accountId := input.AccountId.Ptr()
	dimensionOptionId := input.DimensionOptionId.Ptr()
	if accountId == nil && dimensionOptionId == nil {
		return nil, commonErrors.NewInvalidInputError(commonErrors.SlugLedgerTransactionsMissingFilter)
	}
	pageRequest, err := data.PageRequestFromInput(input.PaginationInput)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	page, err := h.app.Queries.LedgerEntries.Handle(ctx, input.SobId, accountId, input.FromPeriod, input.ToPeriod, dimensionOptionId, pageRequest)
	if err != nil {
		return nil, err
	}
	return &ReadLedgerTransactionsOutput{Body: data.MapPageResponse(page, ledgerEntryDTOToVO)}, nil
}

type ReadLedgerByDimensionCategoryInput struct {
	SobId               uuid.UUID                     `path:"sobId"`
	DimensionCategoryId uuid.UUID                     `path:"dimensionCategoryId"`
	AccountId           data.OptionalParam[uuid.UUID] `query:"accountId" doc:"Optional account filter"`
	FromPeriod          string                        `query:"fromPeriod" doc:"From period (YYYY-MM)"`
	ToPeriod            string                        `query:"toPeriod" doc:"To period (YYYY-MM)"`
	data.PaginationInput
}

type ReadLedgerByDimensionCategoryOutput struct {
	Body data.PageResponse[LedgerDimensionOptionResponse]
}

// ReadLedgerByDimensionCategory summarizes ledger amounts by dimension option.
func (h Handler) ReadLedgerByDimensionCategory(ctx context.Context, input *ReadLedgerByDimensionCategoryInput) (*ReadLedgerByDimensionCategoryOutput, error) {
	pageRequest, err := data.PageRequestFromInput(input.PaginationInput)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	page, err := h.app.Queries.LedgersByDimensionCategory.Handle(ctx, input.SobId, input.DimensionCategoryId, input.AccountId.Ptr(), input.FromPeriod, input.ToPeriod, pageRequest)
	if err != nil {
		return nil, err
	}
	return &ReadLedgerByDimensionCategoryOutput{Body: data.MapPageResponse(page, ledgerDimensionSummaryItemToVO)}, nil
}
