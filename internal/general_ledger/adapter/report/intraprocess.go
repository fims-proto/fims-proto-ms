package report

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/port/private/intraprocess"

	"github.com/google/uuid"
)

type IntraProcessAdapter struct {
	reportInterface intraprocess.ReportInterface
}

func NewIntraProcessAdapter(reportInterface intraprocess.ReportInterface) IntraProcessAdapter {
	return IntraProcessAdapter{reportInterface: reportInterface}
}

func (a IntraProcessAdapter) GenerateForPeriod(ctx context.Context, sobId uuid.UUID, periodId uuid.UUID) error {
	return a.reportInterface.GenerateForPeriod(ctx, sobId, periodId)
}
