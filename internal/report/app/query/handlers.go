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
	return PagingReportsHandler{readModel: readModel}
}

func (h PagingReportsHandler) Handle(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Report], error) {
	return h.readModel.SearchReport(ctx, sobId, pageRequest)
}

type ReportByIdHandler struct {
	readModel ReportReadModel
}

func NewReportByIdHandler(readModel ReportReadModel) ReportByIdHandler {
	return ReportByIdHandler{readModel: readModel}
}

func (h ReportByIdHandler) Handle(ctx context.Context, reportId uuid.UUID) (Report, error) {
	return h.readModel.ReportById(ctx, reportId)
}

type ReportByClassAndPeriodHandler struct {
	readModel ReportReadModel
}

func NewReportByClassAndPeriodHandler(readModel ReportReadModel) ReportByClassAndPeriodHandler {
	return ReportByClassAndPeriodHandler{readModel: readModel}
}

func (h ReportByClassAndPeriodHandler) Handle(ctx context.Context, sobId uuid.UUID, reportClass string, fiscalYear int, periodNumber int) (Report, error) {
	return h.readModel.ReportInstanceByClassAndPeriod(ctx, sobId, reportClass, fiscalYear, periodNumber)
}

type ReportTemplateByClassHandler struct {
	readModel ReportReadModel
}

func NewReportTemplateByClassHandler(readModel ReportReadModel) ReportTemplateByClassHandler {
	return ReportTemplateByClassHandler{readModel: readModel}
}

func (h ReportTemplateByClassHandler) Handle(ctx context.Context, sobId uuid.UUID, reportClass string) (Report, error) {
	return h.readModel.TemplateBySobIdAndClass(ctx, sobId, reportClass)
}
