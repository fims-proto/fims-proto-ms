package http

import (
	"context"
	"fmt"
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/command"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type SearchJournalsInput struct {
	SobId uuid.UUID `path:"sobId"`
	data.PaginationInput
}

type SearchJournalsOutput struct {
	Body data.PageResponse[JournalSlimResponse]
}

// SearchJournals lists journals for a SoB.
func (h Handler) SearchJournals(ctx context.Context, input *SearchJournalsInput) (*SearchJournalsOutput, error) {
	pageRequest, err := data.PageRequestFromInput(input.PaginationInput)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	page, err := h.app.Queries.PagingJournals.Handle(ctx, input.SobId, pageRequest)
	if err != nil {
		return nil, err
	}
	vos := make([]JournalSlimResponse, len(page.Content()))
	for i, dto := range page.Content() {
		vos[i] = journalDTOToSlimVO(dto)
	}
	return &SearchJournalsOutput{Body: data.PageResponse[JournalSlimResponse]{
		Content:          vos,
		PageNumber:       page.PageNumber(),
		PageSize:         page.PageSize(),
		TotalPage:        page.TotalPage(),
		NumberOfElements: page.NumberOfElements(),
	}}, nil
}

type ReadJournalByIdInput struct {
	SobId     uuid.UUID `path:"sobId"`
	JournalId uuid.UUID `path:"journalId"`
}

type ReadJournalByIdOutput struct {
	Body JournalDetailResponse
}

// ReadJournalById returns one journal by ID.
func (h Handler) ReadJournalById(ctx context.Context, input *ReadJournalByIdInput) (*ReadJournalByIdOutput, error) {
	j, err := h.app.Queries.JournalById.Handle(ctx, input.JournalId)
	if err != nil {
		return nil, err
	}
	if j.Id == uuid.Nil {
		return nil, huma.Error404NotFound("journal not found")
	}
	return &ReadJournalByIdOutput{Body: journalDTOToDetailVO(j)}, nil
}

type CreateJournalInput struct {
	SobId uuid.UUID `path:"sobId"`
	Body  CreateJournalRequest
}

type CreateJournalOutput struct {
	Status int
	Body   JournalDetailResponse
}

// CreateJournal creates a journal and returns created detail.
func (h Handler) CreateJournal(ctx context.Context, input *CreateJournalInput) (*CreateJournalOutput, error) {
	cmd := input.Body.mapToCommand(input.SobId)
	if err := h.app.Commands.CreateJournal.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	createdJournal, err := h.app.Queries.JournalById.Handle(ctx, cmd.JournalId)
	if err != nil {
		return nil, err
	}
	return &CreateJournalOutput{Status: http.StatusCreated, Body: journalDTOToDetailVO(createdJournal)}, nil
}

type UpdateJournalInput struct {
	SobId     uuid.UUID `path:"sobId"`
	JournalId uuid.UUID `path:"journalId"`
	Body      UpdateJournalRequest
}

// UpdateJournal updates an unaudited journal.
func (h Handler) UpdateJournal(ctx context.Context, input *UpdateJournalInput) (*struct{}, error) {
	var items []command.JournalLineCmd
	for _, itemReq := range input.Body.JournalLines {
		items = append(items, itemReq.mapToCommand())
	}
	cmd := command.UpdateJournalCmd{
		JournalId:       input.JournalId,
		HeaderText:      input.Body.HeaderText,
		JournalLines:    items,
		TransactionDate: transaction_date.TransactionDate(input.Body.TransactionDate),
		Updater:         input.Body.Updater,
	}
	if err := h.app.Commands.UpdateJournal.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}

type AuditJournalInput struct {
	SobId     uuid.UUID `path:"sobId"`
	JournalId uuid.UUID `path:"journalId"`
	Body      AuditJournalRequest
}

// AuditJournal audits a reviewed journal.
func (h Handler) AuditJournal(ctx context.Context, input *AuditJournalInput) (*struct{}, error) {
	cmd := command.AuditJournalCmd{
		JournalId: input.JournalId,
		Auditor:   input.Body.Auditor,
	}
	if err := h.app.Commands.AuditJournal.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}

type CancelAuditJournalInput struct {
	SobId     uuid.UUID `path:"sobId"`
	JournalId uuid.UUID `path:"journalId"`
	Body      AuditJournalRequest
}

// CancelAuditJournal reverts journal audit status.
func (h Handler) CancelAuditJournal(ctx context.Context, input *CancelAuditJournalInput) (*struct{}, error) {
	cmd := command.CancelAuditJournalCmd{
		JournalId: input.JournalId,
		Auditor:   input.Body.Auditor,
	}
	if err := h.app.Commands.CancelAuditJournal.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}

type ReviewJournalInput struct {
	SobId     uuid.UUID `path:"sobId"`
	JournalId uuid.UUID `path:"journalId"`
	Body      ReviewJournalRequest
}

// ReviewJournal reviews a journal before audit.
func (h Handler) ReviewJournal(ctx context.Context, input *ReviewJournalInput) (*struct{}, error) {
	cmd := command.ReviewJournalCmd{
		JournalId: input.JournalId,
		Reviewer:  input.Body.Reviewer,
	}
	if err := h.app.Commands.ReviewJournal.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}

type CancelReviewJournalInput struct {
	SobId     uuid.UUID `path:"sobId"`
	JournalId uuid.UUID `path:"journalId"`
	Body      ReviewJournalRequest
}

// CancelReviewJournal reverts journal review status.
func (h Handler) CancelReviewJournal(ctx context.Context, input *CancelReviewJournalInput) (*struct{}, error) {
	cmd := command.CancelReviewJournalCmd{
		JournalId: input.JournalId,
		Reviewer:  input.Body.Reviewer,
	}
	if err := h.app.Commands.CancelReviewJournal.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}

type PostJournalInput struct {
	SobId     uuid.UUID `path:"sobId"`
	JournalId uuid.UUID `path:"journalId"`
	Body      PostJournalRequest
}

// PostJournal posts an audited journal to ledgers.
func (h Handler) PostJournal(ctx context.Context, input *PostJournalInput) (*struct{}, error) {
	cmd := command.PostJournalCmd{
		JournalId: input.JournalId,
		Poster:    input.Body.Poster,
	}
	if err := h.app.Commands.PostJournal.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}

type CreateMonthlyClosingJournalInput struct {
	SobId uuid.UUID `path:"sobId"`
}

type CreateMonthlyClosingJournalOutput struct {
	Status int
	Body   ClosingJournalResponse
}

// CreateMonthlyClosingJournal creates a monthly closing journal.
func (h Handler) CreateMonthlyClosingJournal(ctx context.Context, input *CreateMonthlyClosingJournalInput) (*CreateMonthlyClosingJournalOutput, error) {
	journalId, err := h.app.Commands.CreateMonthlyClosingJournal.Handle(ctx,
		command.CreateMonthlyClosingJournalCmd{SobId: input.SobId},
	)
	if err != nil {
		return nil, err
	}
	return &CreateMonthlyClosingJournalOutput{Status: http.StatusCreated, Body: ClosingJournalResponse{JournalId: journalId}}, nil
}

type CreateYearEndClosingJournalInput struct {
	SobId uuid.UUID `path:"sobId"`
}

type CreateYearEndClosingJournalOutput struct {
	Status int
	Body   ClosingJournalResponse
}

// CreateYearEndClosingJournal creates a year-end closing journal.
func (h Handler) CreateYearEndClosingJournal(ctx context.Context, input *CreateYearEndClosingJournalInput) (*CreateYearEndClosingJournalOutput, error) {
	journalId, err := h.app.Commands.CreateYearEndClosingJournal.Handle(ctx,
		command.CreateYearEndClosingJournalCmd{SobId: input.SobId},
	)
	if err != nil {
		return nil, err
	}
	return &CreateYearEndClosingJournalOutput{Status: http.StatusCreated, Body: ClosingJournalResponse{JournalId: journalId}}, nil
}

type GetClosingJournalInput struct {
	SobId  uuid.UUID `path:"sobId"`
	Period string    `query:"period" doc:"Period in YYYY-MM format"`
}

type GetClosingJournalOutput struct {
	Body ClosingJournalIdsResponse
}

// GetClosingJournal returns closing journal IDs for a period.
func (h Handler) GetClosingJournal(ctx context.Context, input *GetClosingJournalInput) (*GetClosingJournalOutput, error) {
	var fiscalYear, periodNumber int
	if _, err := fmt.Sscanf(input.Period, "%d-%d", &fiscalYear, &periodNumber); err != nil {
		return nil, huma.Error400BadRequest("invalid period format, expected YYYY-MM")
	}
	result, err := h.app.Queries.ClosingJournalIdsByPeriod.Handle(ctx, input.SobId, fiscalYear, periodNumber)
	if err != nil {
		return nil, err
	}
	return &GetClosingJournalOutput{Body: ClosingJournalIdsResponse{
		MonthlyClosingJournalId: result.MonthlyClosingJournalId,
		YearEndClosingJournalId: result.YearEndClosingJournalId,
	}}, nil
}

type DeleteSystemJournalInput struct {
	SobId     uuid.UUID `path:"sobId"`
	JournalId uuid.UUID `path:"journalId"`
}

// DeleteSystemJournal deletes a system closing journal and reverses ledger impact.
func (h Handler) DeleteSystemJournal(ctx context.Context, input *DeleteSystemJournalInput) (*struct{}, error) {
	cmd := command.DeleteSystemJournalCmd{
		SobId:     input.SobId,
		JournalId: input.JournalId,
	}
	if err := h.app.Commands.DeleteSystemJournal.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}
