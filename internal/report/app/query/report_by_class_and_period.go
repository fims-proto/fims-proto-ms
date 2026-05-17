package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report/class"

	"github.com/google/uuid"
)

type ReportByClassAndPeriodHandler struct {
	readModel ReportReadModel
}

func NewReportByClassAndPeriodHandler(readModel ReportReadModel) ReportByClassAndPeriodHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return ReportByClassAndPeriodHandler{readModel: readModel}
}

func (h ReportByClassAndPeriodHandler) Handle(
	ctx context.Context,
	sobId uuid.UUID,
	reportClass class.Class,
	fiscalYear int,
	periodNumber int,
) (Report, error) {
	return h.readModel.ReportInstanceByClassAndPeriod(ctx, sobId, reportClass, fiscalYear, periodNumber)
}
