package query

import (
	"context"

	"github.com/google/uuid"
)

type CurrentPeriodHandler struct {
	readModel AccountReadModel
}

func NewCurrentPeriodHandler(readModel AccountReadModel) CurrentPeriodHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return CurrentPeriodHandler{readModel: readModel}
}

func (h CurrentPeriodHandler) Handle(ctx context.Context, sobId uuid.UUID) (Period, error) {
	return h.readModel.CurrentPeriod(ctx, sobId)
}
