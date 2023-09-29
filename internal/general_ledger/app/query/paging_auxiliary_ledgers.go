package query

import (
	"context"

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
	return h.readModel.SearchAuxiliaryLedgers(ctx, pageRequest)
}
