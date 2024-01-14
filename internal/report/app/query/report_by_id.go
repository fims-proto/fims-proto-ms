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

type ReportByIdHandler struct {
	readModel ReportReadModel
}

func NewReportByIdHandler(readModel ReportReadModel) ReportByIdHandler {
	if readModel == nil {
		panic("nil read model")
	}

	return ReportByIdHandler{
		readModel: readModel,
	}
}

func (h ReportByIdHandler) Handle(ctx context.Context, ReportId uuid.UUID) (Report, error) {
	idFilter, err := filterable.NewFilter("id", filterable.OptEq, ReportId)
	if err != nil {
		return Report{}, fmt.Errorf("failed to build filter: %w", err)
	}

	reports, err := h.readModel.SearchReports(ctx, uuid.Nil, data.NewPageRequest(pageable.Unpaged(), sortable.Unsorted(), filterable.NewFilterableAtom(idFilter)))
	if err != nil {
		return Report{}, fmt.Errorf("failed to search Reports: %w", err)
	}
	if reports.NumberOfElements() != 1 {
		return Report{}, commonErrors.ErrRecordNotFound()
	}

	return reports.Content()[0], nil
}
