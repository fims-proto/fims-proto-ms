package domain

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report"

	"github.com/google/uuid"
)

type Repository interface {
	Migrate(ctx context.Context) error
	EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error

	CreateReports(ctx context.Context, reports []*report.Report) error
	UpdateReport(ctx context.Context, reportId uuid.UUID, updateFn func(r *report.Report) (*report.Report, error)) error
	ReadReportById(ctx context.Context, reportId uuid.UUID) (*report.Report, error)
	ReadTemplatesBySobId(ctx context.Context, sobId uuid.UUID) ([]*report.Report, error)
	ReadTemplateBySobIdAndClass(ctx context.Context, sobId uuid.UUID, reportClass string) (*report.Report, error)
	ReadInstanceBySobClassAndPeriod(ctx context.Context, sobId uuid.UUID, reportClass string, periodId uuid.UUID) (*report.Report, error)
}
