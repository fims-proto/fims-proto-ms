package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type PagingLedgersByPeriodHandler struct {
	readModel AccountReadModel
}

func NewPagingLedgersByPeriodHandler(readModel AccountReadModel) PagingLedgersByPeriodHandler {
	if readModel == nil {
		panic("nil account read model")
	}

	return PagingLedgersByPeriodHandler{readModel: readModel}
}

func (h PagingLedgersByPeriodHandler) Handle(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID, pageRequest data.PageRequest) (data.Page[Ledger], error) {
	return h.readModel.PagingLedgersByPeriod(ctx, sobId, periodId, pageRequest)
}
