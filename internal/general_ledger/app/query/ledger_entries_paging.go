package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type LedgerEntriesHandler struct {
	readModel            GeneralLedgerReadModel
	periodRangeValidator periodRangeValidator
}

func NewLedgerEntriesHandler(readModel GeneralLedgerReadModel) LedgerEntriesHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return LedgerEntriesHandler{
		readModel:            readModel,
		periodRangeValidator: newPeriodRangeValidator(readModel),
	}
}

func (h LedgerEntriesHandler) Handle(
	ctx context.Context,
	sobId uuid.UUID,
	accountId *uuid.UUID,
	fromPeriod, toPeriod string,
	dimensionOptionId *uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[LedgerEntry], error) {
	// Validate period continuity
	fromFiscalYear, fromPeriodNumber, toFiscalYear, toPeriodNumber, err := h.periodRangeValidator.validate(ctx, sobId, fromPeriod, toPeriod)
	if err != nil {
		return nil, fmt.Errorf("invalid period range: %w", err)
	}

	// Query ledger entries for the period range with pagination
	entriesPage, err := h.readModel.LedgerEntriesByPeriodRange(
		ctx,
		sobId,
		accountId,
		fromFiscalYear,
		fromPeriodNumber,
		toFiscalYear,
		toPeriodNumber,
		dimensionOptionId,
		pageRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query ledger entries: %w", err)
	}

	return entriesPage, nil
}
