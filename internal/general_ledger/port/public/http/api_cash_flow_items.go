package http

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"

	"github.com/google/uuid"
)

type ReadCashFlowItemsInput struct {
	SobId uuid.UUID `path:"sobId"`
}

type ReadCashFlowItemsOutput struct {
	Body []CashFlowItemResponse
}

// ReadCashFlowItems lists cash flow items for a SoB.
func (h Handler) ReadCashFlowItems(ctx context.Context, input *ReadCashFlowItemsInput) (*ReadCashFlowItemsOutput, error) {
	items, err := h.app.Queries.CashFlowItems.Handle(ctx, input.SobId)
	if err != nil {
		return nil, err
	}
	return &ReadCashFlowItemsOutput{Body: converter.DTOsToVOs(items, cashFlowItemDTOToVO)}, nil
}
