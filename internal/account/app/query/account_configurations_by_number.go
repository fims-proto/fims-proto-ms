package query

import (
	"context"

	"github.com/google/uuid"
)

type AccountConfigurationsByNumbersHandler struct {
	readModel AccountReadModel
}

func NewAccountConfigurationsByNumbersHandler(readModel AccountReadModel) AccountConfigurationsByNumbersHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return AccountConfigurationsByNumbersHandler{readModel: readModel}
}

func (h AccountConfigurationsByNumbersHandler) Handle(ctx context.Context, sobId uuid.UUID, accountNumbers []string) ([]AccountConfiguration, error) {
	return h.readModel.AccountConfigurationsByNumbers(ctx, sobId, accountNumbers)
}
