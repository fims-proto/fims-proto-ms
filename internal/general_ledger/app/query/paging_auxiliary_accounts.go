package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"

	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type PagingAuxiliaryAccountsHandler struct {
	readModel GeneralLedgerReadModel
}

func NewPagingAuxiliaryAccountsHandler(readModel GeneralLedgerReadModel) PagingAuxiliaryAccountsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingAuxiliaryAccountsHandler{readModel: readModel}
}

func (h PagingAuxiliaryAccountsHandler) Handle(ctx context.Context, categoryKey string, pageRequest data.PageRequest) (data.Page[AuxiliaryAccount], error) {
	filter, err := filterable.NewFilter("category.key", "eq", categoryKey)
	if err != nil {
		return nil, fmt.Errorf("failed to build filter: %w", err)
	}

	pageRequest.AddFilter(filter)
	return h.readModel.SearchAuxiliaryAccounts(ctx, pageRequest)
}
