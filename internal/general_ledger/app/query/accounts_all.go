package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github.com/google/uuid"
)

type AllAccountsHandler struct {
	readModel  GeneralLedgerReadModel
	sobService service.SobService
}

func NewAllAccountsHandler(readModel GeneralLedgerReadModel, sobService service.SobService) AllAccountsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	return AllAccountsHandler{
		readModel:  readModel,
		sobService: sobService,
	}
}

func (h AllAccountsHandler) Handle(ctx context.Context, sobId uuid.UUID) ([]Account, error) {
	sort, err := sortable.NewSort("rawAccountNumber", "asc")
	if err != nil {
		panic(fmt.Errorf("failed to build sort 'rawAccountNumber': %w", err))
	}

	pageRequest := data.NewPageRequest(
		pageable.Unpaged(),
		sortable.New(sort),
		filterable.Unfiltered(),
	)
	accounts, err := h.readModel.SearchAccounts(ctx, sobId, pageRequest)
	if err != nil {
		return nil, fmt.Errorf("error getting accounts: %w", err)
	}

	sob, err := h.sobService.ReadById(ctx, sobId)
	if err != nil {
		return nil, fmt.Errorf("failed to read sob: %w", err)
	}

	return enrichAccountNumbers(sob.AccountsCodeLength, accounts.Content()), nil
}
