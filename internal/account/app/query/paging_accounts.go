package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/datav3"

	"github.com/google/uuid"
)

type PagingAccountsHandler struct {
	readModel AccountReadModel
}

func NewPagingAccountsHandler(readModel AccountReadModel) PagingAccountsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingAccountsHandler{readModel: readModel}
}

func (h PagingAccountsHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest datav3.PageRequest) (datav3.Page[Account], error) {
	return h.readModel.SearchAccounts(ctx, sobId, pageRequest)
}
