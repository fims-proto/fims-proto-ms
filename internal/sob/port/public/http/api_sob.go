package http

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/sob/app/command"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type SearchSobsInput struct {
	data.PaginationInput
}

type SearchSobsOutput struct {
	Body data.PageResponse[SobResponse]
}

// SearchSobs lists sets of books.
func (h Handler) SearchSobs(ctx context.Context, input *SearchSobsInput) (*SearchSobsOutput, error) {
	pageRequest, err := data.PageRequestFromInput(input.PaginationInput)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	page, err := h.app.Queries.PagingSobs.Handle(ctx, pageRequest)
	if err != nil {
		return nil, err
	}
	return &SearchSobsOutput{Body: data.MapPageResponse(page, sobDTOToVO)}, nil
}

type ReadSobByIdInput struct {
	SobId uuid.UUID `path:"sobId"`
}

type ReadSobByIdOutput struct {
	Body SobResponse
}

// ReadSobById returns one set of books by ID.
func (h Handler) ReadSobById(ctx context.Context, input *ReadSobByIdInput) (*ReadSobByIdOutput, error) {
	sob, err := h.app.Queries.SobById.Handle(ctx, input.SobId)
	if err != nil {
		return nil, err
	}
	if sob.Id == uuid.Nil {
		return nil, huma.Error404NotFound("sob not found")
	}
	return &ReadSobByIdOutput{Body: sobDTOToVO(sob)}, nil
}

type CreateSobInput struct {
	Body CreateSobRequest
}

type CreateSobOutput struct {
	Status int
	Body   SobResponse
}

// CreateSob creates a set of books and returns created detail.
func (h Handler) CreateSob(ctx context.Context, input *CreateSobInput) (*CreateSobOutput, error) {
	cmd := input.Body.mapToCommand()
	if err := h.app.Commands.CreateSob.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	createdSob, err := h.app.Queries.SobById.Handle(ctx, cmd.SobId)
	if err != nil {
		return nil, err
	}
	return &CreateSobOutput{Status: 201, Body: sobDTOToVO(createdSob)}, nil
}

type UpdateSobInput struct {
	SobId uuid.UUID `path:"sobId"`
	Body  UpdateSobRequest
}

// UpdateSob updates set of books metadata and account code length.
func (h Handler) UpdateSob(ctx context.Context, input *UpdateSobInput) (*struct{}, error) {
	cmd := command.UpdateSobCmd{
		SobId:              input.SobId,
		Name:               input.Body.Name,
		Description:        input.Body.Description,
		AccountsCodeLength: input.Body.AccountsCodeLength,
	}
	if err := h.app.Commands.UpdateSob.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}
