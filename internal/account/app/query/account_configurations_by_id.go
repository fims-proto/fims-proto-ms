package query

import (
	"context"

	"github.com/google/uuid"
)

type AccountConfigurationsByIdsHandler struct {
	readModel AccountReadModel
}

func NewAccountConfigurationsByIdsHandler(readModel AccountReadModel) AccountConfigurationsByIdsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return AccountConfigurationsByIdsHandler{readModel: readModel}
}

func (h AccountConfigurationsByIdsHandler) Handle(ctx context.Context, accountIds []uuid.UUID) ([]AccountConfiguration, error) {
	return h.readModel.AccountConfigurationsByIds(ctx, accountIds)
}
