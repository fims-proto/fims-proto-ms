package query

import (
	"context"

	"github.com/google/uuid"
)

type OpenPeriodHandler struct {
	readModel AccountReadModel
}

func NewOpenPeriodHandler(readModel AccountReadModel) OpenPeriodHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return OpenPeriodHandler{readModel: readModel}
}

func (h OpenPeriodHandler) Handle(ctx context.Context, sobId uuid.UUID) (Period, error) {
	return h.readModel.OpenPeriod(ctx, sobId)
}
