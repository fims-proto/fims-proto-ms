package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type AllAccountsHandler struct {
	readModel GeneralLedgerReadModel
}

func NewAllAccountsHandler(readModel GeneralLedgerReadModel) AllAccountsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return AllAccountsHandler{readModel: readModel}
}

func (h AllAccountsHandler) Handle(ctx context.Context, sobId uuid.UUID) ([]Account, error) {
	sort, err := sortable.NewSort("accountNumber", "asc")
	if err != nil {
		panic(fmt.Errorf("failed to build sort 'accountNumber': %w", err))
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
	return accounts.Content(), nil
}
