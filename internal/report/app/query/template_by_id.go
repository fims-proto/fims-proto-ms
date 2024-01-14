package query

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"

	"github.com/google/uuid"
)

type TemplateByIdHandler struct {
	readModel ReportReadModel
}

func NewTemplateByIdHandler(readModel ReportReadModel) TemplateByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return TemplateByIdHandler{
		readModel: readModel,
	}
}

func (h TemplateByIdHandler) Handle(ctx context.Context, TemplateId uuid.UUID) (Template, error) {
	idFilter, err := filterable.NewFilter("id", filterable.OptEq, TemplateId)
	if err != nil {
		return Template{}, fmt.Errorf("failed to build filter: %w", err)
	}

	templates, err := h.readModel.SearchTemplates(ctx, uuid.Nil, data.NewPageRequest(pageable.Unpaged(), sortable.Unsorted(), filterable.NewFilterableAtom(idFilter)))
	if err != nil {
		return Template{}, fmt.Errorf("failed to search Templates: %w", err)
	}
	if templates.NumberOfElements() != 1 {
		return Template{}, commonErrors.ErrRecordNotFound()
	}

	return templates.Content()[0], nil
}
