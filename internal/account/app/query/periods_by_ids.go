package query

import (
	"context"

	"github.com/google/uuid"
)

type PeriodsByIdsHandler struct {
	readModel AccountReadModel
}

func NewPeriodsByIdsHandler(readModel AccountReadModel) PeriodsByIdsHandler {
	return PeriodsByIdsHandler{readModel: readModel}
}

func (h PeriodsByIdsHandler) Handle(ctx context.Context, periodIds []uuid.UUID) ([]Period, error) {
	return h.readModel.PeriodsByIds(ctx, periodIds)
}
