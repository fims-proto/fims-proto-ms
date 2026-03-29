package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github.com/google/uuid"
)

type FirstPeriodLedgersHandler struct {
	readModel  GeneralLedgerReadModel
	sobService service.SobService
}

func NewFirstPeriodLedgersHandler(readModel GeneralLedgerReadModel, sobService service.SobService) FirstPeriodLedgersHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	return FirstPeriodLedgersHandler{
		readModel:  readModel,
		sobService: sobService,
	}
}

func (h FirstPeriodLedgersHandler) Handle(ctx context.Context, sobId uuid.UUID) (Period, []Ledger, error) {
	firstPeriod, err := h.readModel.FirstPeriod(ctx, sobId)
	if err != nil {
		return Period{}, nil, fmt.Errorf("error getting first period: %w", err)
	}

	periodIdFilter, err := filterable.NewFilter("periodId", filterable.OptEq, firstPeriod.Id)
	if err != nil {
		panic(fmt.Errorf("failed to build filter 'periodId': %w", err))
	}

	sort, err := sortable.NewSort("account.rawAccountNumber", "asc")
	if err != nil {
		panic(fmt.Errorf("failed to build sort 'account.rawAccountNumber': %w", err))
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

	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return Period{}, nil, fmt.Errorf("failed to read sob: %w", err)
	}

	return firstPeriod, enrichLedgerAccountNumbers(sob.AccountsCodeLength, ledgers.Content()), nil
}
