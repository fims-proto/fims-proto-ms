package query

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
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

func (h PagingAccountsHandler) Handle(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Account], error) {
	return h.readModel.PagingAccounts(ctx, sobId, pageable)
}
