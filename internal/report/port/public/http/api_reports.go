package http

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/report/app/command"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type SearchReportsInput struct {
	SobId uuid.UUID `path:"sobId"`
	data.PaginationInput
}

type SearchReportsOutput struct {
	Body data.PageResponse[ReportResponse]
}

// SearchReports lists reports for a SoB.
func (h Handler) SearchReports(ctx context.Context, input *SearchReportsInput) (*SearchReportsOutput, error) {
	pageRequest, err := data.PageRequestFromInput(input.PaginationInput)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	page, err := h.app.Queries.PagingReports.Handle(ctx, input.SobId, pageRequest)
	if err != nil {
		return nil, err
	}
	return &SearchReportsOutput{Body: data.MapPageResponse(page, reportDTOToVO)}, nil
}

type ReadReportTemplateByClassInput struct {
	SobId uuid.UUID `path:"sobId"`
	Class string    `query:"class" doc:"Report class"`
}

type ReadReportTemplateByClassOutput struct {
	Body ReportResponse
}

// ReadReportTemplateByClass returns report template by class.
func (h Handler) ReadReportTemplateByClass(ctx context.Context, input *ReadReportTemplateByClassInput) (*ReadReportTemplateByClassOutput, error) {
	r, err := h.app.Queries.ReportTemplateByClass.Handle(ctx, input.SobId, input.Class)
	if err != nil {
		return nil, err
	}
	if r.Id == uuid.Nil {
		return nil, huma.Error404NotFound("report template not found")
	}
	return &ReadReportTemplateByClassOutput{Body: reportDTOToVO(r)}, nil
}

type ReadReportByClassAndPeriodInput struct {
	SobId  uuid.UUID `path:"sobId"`
	Class  string    `query:"class" doc:"Report class"`
	Period string    `query:"period" doc:"Period in YYYY-MM format"`
}

type ReadReportByClassAndPeriodOutput struct {
	Body ReportResponse
}

// ReadReportByClassAndPeriod returns generated report by class and period.
func (h Handler) ReadReportByClassAndPeriod(ctx context.Context, input *ReadReportByClassAndPeriodInput) (*ReadReportByClassAndPeriodOutput, error) {
	t, err := time.Parse("2006-01", input.Period)
	if err != nil {
		return nil, huma.Error400BadRequest("period must be YYYY-MM")
	}
	r, err := h.app.Queries.ReportByClassAndPeriod.Handle(ctx, input.SobId, input.Class, t.Year(), int(t.Month()))
	if err != nil {
		return nil, err
	}
	if r.Id == uuid.Nil {
		return nil, huma.Error404NotFound("report not found")
	}
	return &ReadReportByClassAndPeriodOutput{Body: reportDTOToVO(r)}, nil
}

type GenerateReportInput struct {
	SobId  uuid.UUID `path:"sobId"`
	Class  string    `query:"class" doc:"Report class"`
	Period string    `query:"period" doc:"Period in YYYY-MM format"`
}

type GenerateReportOutput struct {
	Body ReportResponse
}

// GenerateReport generates a report for one class and period.
func (h Handler) GenerateReport(ctx context.Context, input *GenerateReportInput) (*GenerateReportOutput, error) {
	t, err := time.Parse("2006-01", input.Period)
	if err != nil {
		return nil, huma.Error400BadRequest("period must be YYYY-MM")
	}
	actualId, err := h.app.Commands.Generate.Handle(ctx, command.GenerateReportCmd{
		SobId:        input.SobId,
		Class:        input.Class,
		FiscalYear:   t.Year(),
		PeriodNumber: int(t.Month()),
	})
	if err != nil {
		return nil, err
	}
	generatedReport, err := h.app.Queries.ReportById.Handle(ctx, actualId)
	if err != nil {
		return nil, err
	}
	return &GenerateReportOutput{Body: reportDTOToVO(generatedReport)}, nil
}

type RecalculateReportInput struct {
	SobId    uuid.UUID `path:"sobId"`
	ReportId uuid.UUID `path:"reportId"`
}

// RecalculateReport recalculates report values in place.
func (h Handler) RecalculateReport(ctx context.Context, input *RecalculateReportInput) (*struct{}, error) {
	if err := h.app.Commands.Recalculate.Handle(ctx, command.RecalculateReportCmd{ReportId: input.ReportId}); err != nil {
		return nil, err
	}
	return nil, nil
}

type RegenerateReportInput struct {
	SobId    uuid.UUID `path:"sobId"`
	ReportId uuid.UUID `path:"reportId"`
}

// RegenerateReport regenerates report structure and values.
func (h Handler) RegenerateReport(ctx context.Context, input *RegenerateReportInput) (*struct{}, error) {
	if err := h.app.Commands.Regenerate.Handle(ctx, command.RegenerateReportCmd{ReportId: input.ReportId}); err != nil {
		return nil, err
	}
	return nil, nil
}

type UpdateReportInput struct {
	SobId    uuid.UUID `path:"sobId"`
	ReportId uuid.UUID `path:"reportId"`
	Body     UpdateReportRequest
}

// UpdateReport updates report rows and formulas.
func (h Handler) UpdateReport(ctx context.Context, input *UpdateReportInput) (*struct{}, error) {
	cmd, err := input.Body.mapToCommand(input.ReportId, input.SobId)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	if err = h.app.Commands.UpdateReport.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}
