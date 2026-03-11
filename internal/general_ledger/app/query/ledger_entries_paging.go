package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type PagingLedgerEntriesHandler struct {
	readModel            GeneralLedgerReadModel
	periodRangeValidator periodRangeValidator
}

func NewPagingLedgerEntriesHandler(readModel GeneralLedgerReadModel) PagingLedgerEntriesHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return PagingLedgerEntriesHandler{
		readModel:            readModel,
		periodRangeValidator: newPeriodRangeValidator(readModel),
	}
}

func (h PagingLedgerEntriesHandler) Handle(
	ctx context.Context,
	sobId, accountId uuid.UUID,
	fromPeriod, toPeriod string,
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
		pageRequest,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query ledger entries: %w", err)
	}

	return entriesPage, nil
}
