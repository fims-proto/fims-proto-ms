package query

import (
	"context"

	"github.com/google/uuid"
)

type CurrentPeriodHandler struct {
	readModel GeneralLedgerReadModel
}

func NewCurrentPeriodHandler(readModel GeneralLedgerReadModel) CurrentPeriodHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return CurrentPeriodHandler{readModel: readModel}
}

func (h CurrentPeriodHandler) Handle(ctx context.Context, sobId uuid.UUID) (Period, error) {
	return h.readModel.CurrentPeriod(ctx, sobId)
}
