package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type ReportReadModel interface {
	SearchReport(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Report], error)
	ReportById(ctx context.Context, reportId uuid.UUID) (Report, error)
	ReportInstanceByClassAndPeriod(ctx context.Context, sobId uuid.UUID, reportClass string, fiscalYear int, periodNumber int) (Report, error)
	TemplateBySobIdAndClass(ctx context.Context, sobId uuid.UUID, reportClass string) (Report, error)
}
