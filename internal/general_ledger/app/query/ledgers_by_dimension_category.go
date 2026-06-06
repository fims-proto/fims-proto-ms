package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type LedgersByDimensionCategoryHandler struct {
	readModel            GeneralLedgerReadModel
	periodRangeValidator periodRangeValidator
}

func NewLedgersByDimensionCategoryHandler(readModel GeneralLedgerReadModel) LedgersByDimensionCategoryHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return LedgersByDimensionCategoryHandler{
		readModel:            readModel,
		periodRangeValidator: newPeriodRangeValidator(readModel),
	}
}

func (h LedgersByDimensionCategoryHandler) Handle(
	ctx context.Context,
	sobId, dimensionCategoryId uuid.UUID,
	accountId *uuid.UUID,
	fromPeriod, toPeriod string,
	pageRequest data.PageRequest,
) (data.Page[LedgerDimensionSummaryItem], error) {
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber, err := h.periodRangeValidator.validate(ctx, sobId, fromPeriod, toPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid period range: %w", err)
	}

	return h.readModel.LedgersByAccountAndDimensionOption(
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
