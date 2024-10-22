package report

import (
	"context"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/report/app/command"
	"github/fims-proto/fims-proto-ms/internal/report/port/private/intraprocess"
)

type IntraProcessAdapter struct {
	reportInterface intraprocess.ReportInterface
}

func NewIntraProcessAdapter(reportInterface intraprocess.ReportInterface) IntraProcessAdapter {
	return IntraProcessAdapter{reportInterface: reportInterface}
}

func (i IntraProcessAdapter) InitializeReport(ctx context.Context, sobId uuid.UUID) error {
	return i.reportInterface.Initialize(ctx, command.InitializeCmd{SobId: sobId})
}
