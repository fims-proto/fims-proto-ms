package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type ReportReadModel interface {
	SearchReportInfos(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[ReportInfo], error)
	SearchTemplateInfos(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[TemplateInfo], error)
	SearchReports(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Report], error)
	SearchTemplates(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Template], error)

	SearchPeriods(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Period], error)

	ReportById(ctx context.Context, reportId uuid.UUID) (Report, error)
	TemplateById(ctx context.Context, templateId uuid.UUID) (Template, error)
	SearchAccounts(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Account], error)
}
