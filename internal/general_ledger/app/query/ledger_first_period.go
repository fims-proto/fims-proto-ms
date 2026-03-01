package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"

	"github.com/google/uuid"
)

type FirstPeriodLedgersHandler struct {
	readModel GeneralLedgerReadModel
}

func NewFirstPeriodLedgersHandler(readModel GeneralLedgerReadModel) FirstPeriodLedgersHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return FirstPeriodLedgersHandler{readModel: readModel}
}

func (h FirstPeriodLedgersHandler) Handle(ctx context.Context, sobId uuid.UUID) (Period, []Ledger, error) {
	// find first period
	firstPeriod, err := h.readModel.FirstPeriod(ctx, sobId)
	if err != nil {
		return Period{}, nil, fmt.Errorf("error getting first period: %w", err)
	}

	// find ledgers
	periodIdFilter, err := filterable.NewFilter("periodId", filterable.OptEq, firstPeriod.Id)
	if err != nil {
		panic(fmt.Errorf("failed to build filter 'periodId': %w", err))
	}

	sort, err := sortable.NewSort("account.accountNumber", "asc")
	if err != nil {
		panic(fmt.Errorf("failed to build sort 'account.accountNumber': %w", err))
	}

	ledgerRequest := data.NewPageRequest(
		pageable.Unpaged(),
		sortable.New(sort),
		filterable.NewFilterableAtom(periodIdFilter),
	)
	ledgers, err := h.readModel.SearchLedgers(ctx, sobId, ledgerRequest)
	if err != nil {
		return Period{}, nil, fmt.Errorf("error getting first period ledgers: %w", err)
	}
	return firstPeriod, ledgers.Content(), nil
}
