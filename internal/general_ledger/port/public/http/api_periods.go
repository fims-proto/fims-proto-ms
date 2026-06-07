package http

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/data/converter"
	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ReadPeriodsInput struct {
	SobId uuid.UUID `path:"sobId"`
}

type ReadPeriodsOutput struct {
	Body []PeriodResponse
}

// ReadPeriods lists accounting periods for a SoB.
func (h Handler) ReadPeriods(ctx context.Context, input *ReadPeriodsInput) (*ReadPeriodsOutput, error) {
	periods, err := h.app.Queries.AllPeriods.Handle(ctx, input.SobId)
	if err != nil {
		return nil, err
	}
	return &ReadPeriodsOutput{Body: converter.DTOsToVOs(periods, periodDTOToVO)}, nil
}

type ClosePeriodInput struct {
	AcceptLanguage string    `header:"Accept-Language"`
	SobId          uuid.UUID `path:"sobId"`
	PeriodId       uuid.UUID `path:"periodId"`
}

type ClosePeriodOutput struct {
	Body *PeriodCloseWarningResponse
}

// ClosePeriod closes one accounting period.
func (h Handler) ClosePeriod(ctx context.Context, input *ClosePeriodInput) (*ClosePeriodOutput, error) {
	result, err := h.app.Commands.ClosePeriod.Handle(ctx, command.ClosePeriodCmd{
		SobId:    input.SobId,
		PeriodId: input.PeriodId,
	})
	if err != nil {
		return nil, err
	}
	if result.ReportGenerationFailed {
		msg := h.localizer.Get(input.AcceptLanguage, commonErrors.SlugPeriodClosedButReportFailed, nil)
		return &ClosePeriodOutput{Body: &PeriodCloseWarningResponse{
			Slug:    commonErrors.SlugPeriodClosedButReportFailed,
			Message: msg,
		}}, nil
	}
	return &ClosePeriodOutput{}, nil
}

type PreCloseCheckInput struct {
	SobId    uuid.UUID `path:"sobId"`
	PeriodId uuid.UUID `path:"periodId"`
}

type PreCloseCheckOutput struct {
	Body PreCloseCheckResponse
}

// PreCloseCheck checks whether one period can close.
func (h Handler) PreCloseCheck(ctx context.Context, input *PreCloseCheckInput) (*PreCloseCheckOutput, error) {
	result, err := h.app.Queries.PeriodPreCloseCheck.Handle(ctx, input.SobId, input.PeriodId)
	if err != nil {
		return nil, err
	}
	return &PreCloseCheckOutput{Body: preCloseCheckDTOToVO(result)}, nil
}

type BatchPreCloseCheckInput struct {
	SobId        uuid.UUID `path:"sobId"`
	TargetPeriod string    `query:"targetPeriod" doc:"Target period in YYYY-MM format"`
}

type BatchPreCloseCheckOutput struct {
	Body BatchPreCloseCheckResponse
}

// BatchPreCloseCheck checks close readiness through a target period.
func (h Handler) BatchPreCloseCheck(ctx context.Context, input *BatchPreCloseCheckInput) (*BatchPreCloseCheckOutput, error) {
	var targetYear, targetMonth int
	if _, err := fmt.Sscanf(input.TargetPeriod, "%d-%d", &targetYear, &targetMonth); err != nil {
		return nil, huma.Error400BadRequest("invalid targetPeriod format, expected YYYY-MM")
	}
	result, err := h.app.Queries.BatchPeriodPreCloseCheck.Handle(ctx, input.SobId, targetYear, targetMonth)
	if err != nil {
		return nil, err
	}
	return &BatchPreCloseCheckOutput{Body: batchPreCloseCheckDTOToVO(result)}, nil
}

type ClosePeriodsInput struct {
	AcceptLanguage string    `header:"Accept-Language"`
	SobId          uuid.UUID `path:"sobId"`
	TargetPeriod   string    `query:"targetPeriod" doc:"Target period in YYYY-MM format"`
}

type ClosePeriodsOutput struct {
	Body *PeriodCloseWarningResponse
}

// ClosePeriods closes all open periods through a target period.
func (h Handler) ClosePeriods(ctx context.Context, input *ClosePeriodsInput) (*ClosePeriodsOutput, error) {
	var targetYear, targetMonth int
	if _, err := fmt.Sscanf(input.TargetPeriod, "%d-%d", &targetYear, &targetMonth); err != nil {
		return nil, huma.Error400BadRequest("invalid targetPeriod format, expected YYYY-MM")
	}
	result, err := h.app.Commands.ClosePeriods.Handle(ctx, command.ClosePeriodsCmd{
		SobId:       input.SobId,
		TargetYear:  targetYear,
		TargetMonth: targetMonth,
	})
	if err != nil {
		return nil, err
	}
	if result.ReportGenerationFailed {
		msg := h.localizer.Get(input.AcceptLanguage, commonErrors.SlugPeriodClosedButReportFailed, nil)
		return &ClosePeriodsOutput{Body: &PeriodCloseWarningResponse{
			Slug:    commonErrors.SlugPeriodClosedButReportFailed,
			Message: msg,
		}}, nil
	}
	return &ClosePeriodsOutput{}, nil
}
