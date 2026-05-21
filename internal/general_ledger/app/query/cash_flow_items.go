package query

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CashFlowItemsHandler struct {
	readModel GeneralLedgerReadModel
}

func NewCashFlowItemsHandler(readModel GeneralLedgerReadModel) CashFlowItemsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return CashFlowItemsHandler{readModel: readModel}
}

func (h CashFlowItemsHandler) Handle(ctx context.Context, sobId uuid.UUID) ([]CashFlowItem, error) {
	items, err := h.readModel.CashFlowItemsBySobId(ctx, sobId)
	if err != nil {
		return nil, fmt.Errorf("error getting cash flow items: %w", err)
	}

	return items, nil
}
