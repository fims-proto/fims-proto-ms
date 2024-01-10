package domain

import (
	"context"
	"text/template"

	"github/fims-proto/fims-proto-ms/internal/report/domain/report"

	"github.com/google/uuid"
)

type Repository interface {
	Migrate(ctx context.Context) error
	EnableTx(ctx context.Context, txFn func(txCtx context.Context) error) error

	InitializeTemplates(ctx context.Context, templates []*template.Template) error
	ReadAllTemplatesId(ctx context.Context) ([]uuid.UUID, error)
	ReadTemplateById(ctx context.Context, templateId uuid.UUID) (*template.Template, error)

	ApplyTemplateToReport(ctx context.Context, templateId uuid.UUID, reportId uuid.UUID) (report.Report, error)
	CreatetTemplate(ctx context.Context, template *template.Template) (*template.Template, error)
	UpdateTemplate(
		ctx context.Context,
		template *template.Template,
		updateFn func(t *template.Template) (*template.Template, error),
	) error

	ReadAllReportsId(ctx context.Context) ([]uuid.UUID, error)
	ReadReportById(ctx context.Context, reportId uuid.UUID) (*report.Report, error)
	CreateReport(ctx context.Context, template *template.Template) (*report.Report, error)
	SaveNewTemplateFromReport(ctx context.Context, report report.Report) (*template.Template, error)
	UpdateReport(
		ctx context.Context,
		report *report.Report,
		updateFn func(t *template.Template) (*template.Template, error),
	) error
}
