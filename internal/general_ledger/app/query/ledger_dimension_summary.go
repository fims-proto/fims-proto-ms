package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type LedgerDimensionSummaryHandler struct {
	readModel            GeneralLedgerReadModel
	periodRangeValidator periodRangeValidator
}

func NewLedgerDimensionSummaryHandler(
	readModel GeneralLedgerReadModel,
) LedgerDimensionSummaryHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return LedgerDimensionSummaryHandler{
		readModel:            readModel,
		periodRangeValidator: newPeriodRangeValidator(readModel),
	}
}

func (h LedgerDimensionSummaryHandler) Handle(
	ctx context.Context,
	sobId, accountId, dimensionCategoryId uuid.UUID,
	fromPeriod, toPeriod string,
	pageRequest data.PageRequest,
) (data.Page[LedgerDimensionSummaryItem], error) {
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber, err := h.periodRangeValidator.validate(ctx, sobId, fromPeriod, toPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid period range: %w", err)
	}

	return h.readModel.LedgerDimensionSummary(
		ctx,
		sobId,
		accountId,
		dimensionCategoryId,
		fromFiscalYear,
		fromPeriodNumber,
		toFiscalYear,
		toPeriodNumber,
		pageRequest,
	)
}
