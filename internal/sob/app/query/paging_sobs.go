package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"
)

type PagingSobsHandler struct {
	readModel SobReadModel
}

func NewPagingSobsHandler(readModel SobReadModel) PagingSobsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingSobsHandler{readModel: readModel}
}

func (h PagingSobsHandler) Handle(ctx context.Context, pageRequest data.PageRequest) (data.Page[Sob], error) {
	return h.readModel.SearchSobs(ctx, pageRequest)
}
