package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type PagingTemplateInfosHandler struct {
	readModel ReportReadModel
}

func NewPagingTemplateInfosHandler(readModel ReportReadModel) PagingTemplateInfosHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return PagingTemplateInfosHandler{
		readModel: readModel,
	}
}

func (h PagingTemplateInfosHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[TemplateInfo], error) {
	entriesPage, err := h.readModel.SearchTemplateInfos(ctx, sobId, pageRequest)
	if err != nil {
		return nil, err
	}

	templateInfos := entriesPage.Content()

	return data.NewPage(templateInfos, pageRequest, entriesPage.NumberOfElements())
}
