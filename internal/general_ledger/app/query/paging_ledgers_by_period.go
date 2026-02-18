package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"

	"github.com/google/uuid"
)

type PagingLedgersByPeriodHandler struct {
	readModel GeneralLedgerReadModel
}

func NewPagingLedgersByPeriodHandler(readModel GeneralLedgerReadModel) PagingLedgersByPeriodHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingLedgersByPeriodHandler{readModel: readModel}
}

func (h PagingLedgersByPeriodHandler) Handle(ctx context.Context, sobId, periodId uuid.UUID, pageRequest data.PageRequest) (data.Page[Ledger], error) {
	periodIdFilter, err := filterable.NewFilter("periodId", filterable.OptEq, periodId)
	if err != nil {
		panic(fmt.Errorf("failed to build filter 'periodId': %w", err))
	}
	pageRequest.AddAndFilterable(filterable.NewFilterableAtom(periodIdFilter))
	return h.readModel.SearchLedgers(ctx, sobId, pageRequest)
}
