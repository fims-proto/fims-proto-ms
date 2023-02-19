package query

import (
	"context"

	"github.com/google/uuid"
)

type PeriodByIdHandler struct {
	readModel AccountReadModel
}

func NewPeriodByIdHandler(readModel AccountReadModel) PeriodByIdHandler {
	return PeriodByIdHandler{readModel: readModel}
}

func (h PeriodByIdHandler) Handle(ctx context.Context, periodId uuid.UUID) (Period, error) {
	return h.readModel.PeriodById(ctx, periodId)
}
