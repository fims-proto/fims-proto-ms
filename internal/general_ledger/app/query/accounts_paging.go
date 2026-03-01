package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type PagingAccountsHandler struct {
	readModel GeneralLedgerReadModel
}

func NewPagingAccountsHandler(readModel GeneralLedgerReadModel) PagingAccountsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingAccountsHandler{readModel: readModel}
}

func (h PagingAccountsHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Account], error) {
	return h.readModel.SearchAccounts(ctx, sobId, pageRequest)
}
