package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github.com/google/uuid"
)

type LedgerDimensionSummaryHandler struct {
	readModel            GeneralLedgerReadModel
	periodRangeValidator periodRangeValidator
	dimensionService     service.DimensionService
}

func NewLedgerDimensionSummaryHandler(
	readModel GeneralLedgerReadModel,
	dimensionService service.DimensionService,
) LedgerDimensionSummaryHandler {
	if readModel == nil {
		panic("nil read model")
	}
	if dimensionService == nil {
		panic("nil dimension service")
	}
	return LedgerDimensionSummaryHandler{
		readModel:            readModel,
		periodRangeValidator: newPeriodRangeValidator(readModel),
		dimensionService:     dimensionService,
	}
}

func (h LedgerDimensionSummaryHandler) Handle(
	ctx context.Context,
	sobId, accountId, dimensionCategoryId uuid.UUID,
	fromPeriod, toPeriod string,
) ([]LedgerDimensionSummaryItem, error) {
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber, err := h.periodRangeValidator.validate(ctx, sobId, fromPeriod, toPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid period range: %w", err)
	}

	items, err := h.readModel.LedgerDimensionSummary(
		ctx,
		sobId,
		accountId,
		dimensionCategoryId,
		fromFiscalYear,
		fromPeriodNumber,
		toFiscalYear,
		toPeriodNumber,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query ledger dimension summary: %w", err)
	}

	// Collect all dimension option IDs for batch enrichment
	optionIds := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		optionIds = append(optionIds, item.DimensionOptionId)
	}

	optionsMap, err := h.dimensionService.FetchOptionsByIds(ctx, optionIds)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dimension options: %w", err)
	}

	// Enrich items with option names
	for i, item := range items {
		if opt, ok := optionsMap[item.DimensionOptionId]; ok {
			items[i].DimensionOptionName = opt.Name
		}
	}

	return items, nil
}
