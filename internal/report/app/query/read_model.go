package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/report/domain/report/class"

	"github.com/google/uuid"
)

type ReportReadModel interface {
	SearchReport(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Report], error)
	ReportInstanceByClassAndPeriod(ctx context.Context, sobId uuid.UUID, reportClass class.Class, fiscalYear int, periodNumber int) (Report, error)
	TemplateBySobIdAndClass(ctx context.Context, sobId uuid.UUID, reportClass class.Class) (Report, error)
}
