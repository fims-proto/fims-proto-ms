package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type PagingReportInfosHandler struct {
	readModel ReportReadModel
}

func NewPagingReportInfosHandler(readModel ReportReadModel) PagingReportInfosHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingReportInfosHandler{
		readModel: readModel,
	}
}

func (h PagingReportInfosHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[ReportInfo], error) {
	entriesPage, err := h.readModel.SearchReportInfos(ctx, sobId, pageRequest)
	if err != nil {
		return nil, err
	}

	ReportInfos := entriesPage.Content()

	return data.NewPage(ReportInfos, pageRequest, entriesPage.NumberOfElements())
}
