package query

import (
	"context"

	"github.com/google/uuid"
)

type AccountsByNumbersHandler struct {
	readModel AccountReadModel
}

func NewAccountsByNumbersHandler(readModel AccountReadModel) AccountsByNumbersHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return AccountsByNumbersHandler{readModel: readModel}
}

func (h AccountsByNumbersHandler) Handle(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]Account, error) {
	return h.readModel.AccountsByNumbers(ctx, sobId, accountNumbers)
}
