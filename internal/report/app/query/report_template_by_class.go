package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report/class"

	"github.com/google/uuid"
)

type ReportTemplateByClassHandler struct {
	readModel ReportReadModel
}

func NewReportTemplateByClassHandler(readModel ReportReadModel) ReportTemplateByClassHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return ReportTemplateByClassHandler{readModel: readModel}
}

func (h ReportTemplateByClassHandler) Handle(
	ctx context.Context,
	sobId uuid.UUID,
	reportClass class.Class,
) (Report, error) {
	return h.readModel.TemplateBySobIdAndClass(ctx, sobId, reportClass)
}
