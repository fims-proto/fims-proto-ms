package query

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type PagingLedgersByPeriodHandler struct {
	readModel GeneralLedgerReadModel
}

func NewPagingLedgersByPeriodHandler(readModel GeneralLedgerReadModel) PagingLedgersByPeriodHandler {
	if readModel == nil {
		panic("nil account read model")
	}

	return PagingLedgersByPeriodHandler{readModel: readModel}
}

func (h PagingLedgersByPeriodHandler) Handle(ctx context.Context, sobId, periodId uuid.UUID, pageRequest data.PageRequest) (data.Page[Ledger], error) {
	return h.readModel.PagingLedgersByPeriod(ctx, sobId, periodId, pageRequest)
}
