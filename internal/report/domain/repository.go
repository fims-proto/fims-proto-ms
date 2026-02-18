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
	UpdateReport(
		ctx context.Context,
		reportId uuid.UUID,
		updateFn func(r *report.Report) (*report.Report, error),
	) error
	ReadReportById(ctx context.Context, reportId uuid.UUID) (*report.Report, error)

	UpdateItem(
		ctx context.Context,
		itemId uuid.UUID,
		updateFn func(i *report.Item) (*report.Item, error),
	) error
}
