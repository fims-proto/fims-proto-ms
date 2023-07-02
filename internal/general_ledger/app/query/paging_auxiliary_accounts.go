package query

import (
	"context"

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
	return h.readModel.SearchAuxiliaryAccounts(ctx, pageRequest)
}
