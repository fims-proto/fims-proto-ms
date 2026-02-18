package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"

	"github.com/google/uuid"

	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type PagingAuxiliaryLedgersHandler struct {
	readModel GeneralLedgerReadModel
}

func NewPagingAuxiliaryLedgersHandler(readModel GeneralLedgerReadModel) PagingAuxiliaryLedgersHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingAuxiliaryLedgersHandler{readModel: readModel}
}

func (h PagingAuxiliaryLedgersHandler) Handle(ctx context.Context, sobId, periodId uuid.UUID, pageRequest data.PageRequest) (data.Page[AuxiliaryLedger], error) {
	filter, err := filterable.NewFilter("periodId", filterable.OptEq, periodId)
	if err != nil {
		panic(fmt.Errorf("failed to build filter: %w", err))
	}

	pageRequest.AddAndFilterable(filterable.NewFilterableAtom(filter))
	return h.readModel.SearchAuxiliaryLedgers(ctx, sobId, pageRequest)
}
