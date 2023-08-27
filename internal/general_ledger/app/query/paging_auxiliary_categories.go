package query

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type PagingAuxiliaryCategoriesHandler struct {
	readModel GeneralLedgerReadModel
}

func NewPagingAuxiliaryCategoriesHandler(readModel GeneralLedgerReadModel) PagingAuxiliaryCategoriesHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingAuxiliaryCategoriesHandler{readModel: readModel}
}

func (h PagingAuxiliaryCategoriesHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[AuxiliaryCategory], error) {
	return h.readModel.SearchAuxiliaryCategories(ctx, sobId, pageRequest)
}
