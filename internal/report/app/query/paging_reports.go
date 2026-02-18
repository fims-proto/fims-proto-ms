package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type PagingReportsHandler struct {
	readModel ReportReadModel
}

func NewPagingReportsHandler(readModel ReportReadModel) PagingReportsHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingReportsHandler{
		readModel: readModel,
	}
}

func (h PagingReportsHandler) Handle(
	ctx context.Context,
	sobId uuid.UUID,
	pageRequest data.PageRequest,
) (data.Page[Report], error) {
	return h.readModel.SearchReport(ctx, sobId, pageRequest)
}
