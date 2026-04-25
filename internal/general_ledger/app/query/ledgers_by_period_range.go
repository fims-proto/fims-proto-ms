package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type LedgersByPeriodRangeHandler struct {
	readModel            GeneralLedgerReadModel
	periodRangeValidator periodRangeValidator
}

func NewLedgersByPeriodRangeHandler(readModel GeneralLedgerReadModel) LedgersByPeriodRangeHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return LedgersByPeriodRangeHandler{
		readModel:            readModel,
		periodRangeValidator: newPeriodRangeValidator(readModel),
	}
}

func (h LedgersByPeriodRangeHandler) Handle(
	ctx context.Context,
	sobId uuid.UUID,
	fromPeriod, toPeriod string,
	dimensionOptionId *uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[Ledger], error) {
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber, err := h.periodRangeValidator.validate(ctx, sobId, fromPeriod, toPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid period range: %w", err)
	}

	ledgers, err := h.readModel.LedgersByPeriodRange(ctx, sobId, fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber, dimensionOptionId, pageRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ledgers: %w", err)
	}

	return ledgers, nil
}
