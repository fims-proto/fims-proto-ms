package query

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PeriodByTimeHandler struct {
	readModel AccountReadModel
}

func NewPeriodByTimeHandler(readModel AccountReadModel) PeriodByTimeHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PeriodByTimeHandler{readModel: readModel}
}

func (h PeriodByTimeHandler) Handle(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (Period, error) {
	return h.readModel.PeriodByTime(ctx, sobId, timePoint)
}
