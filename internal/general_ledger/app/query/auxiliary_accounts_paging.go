package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
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

func (h PagingAuxiliaryAccountsHandler) Handle(ctx context.Context, sobId uuid.UUID, categoryKey string, pageRequest data.PageRequest) (data.Page[AuxiliaryAccount], error) {
	filter, err := filterable.NewFilter("category.key", filterable.OptEq, categoryKey)
	if err != nil {
		panic(fmt.Errorf("failed to build filter: %w", err))
	}

	pageRequest.AddAndFilterable(filterable.NewFilterableAtom(filter))
	return h.readModel.SearchAuxiliaryAccounts(ctx, sobId, pageRequest)
}
