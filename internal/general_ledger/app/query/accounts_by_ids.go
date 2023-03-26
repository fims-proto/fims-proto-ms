package query

import (
	"context"

	"github.com/google/uuid"
)

type AccountsByIdsHandler struct {
	readModel GeneralLedgerReadModel
}

func NewAccountsByIdsHandler(readModel GeneralLedgerReadModel) AccountsByIdsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return AccountsByIdsHandler{readModel: readModel}
}

func (h AccountsByIdsHandler) Handle(ctx context.Context, accountIds []uuid.UUID) ([]Account, error) {
	return h.readModel.AccountsByIds(ctx, accountIds)
}
