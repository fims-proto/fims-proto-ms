package http

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/dimension/app/command"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type SearchCategoriesInput struct {
	SobId uuid.UUID `path:"sobId"`
	data.PaginationInput
}

type SearchCategoriesOutput struct {
	Body data.PageResponse[CategoryResponse]
}

// SearchCategories lists dimension categories for a SoB.
func (h Handler) SearchCategories(ctx context.Context, input *SearchCategoriesInput) (*SearchCategoriesOutput, error) {
	pageRequest, err := data.PageRequestFromInput(input.PaginationInput)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	page, err := h.app.Queries.PagingCategories.Handle(ctx, input.SobId, pageRequest)
	if err != nil {
		return nil, err
	}
	return &SearchCategoriesOutput{Body: data.MapPageResponse(page, categoryDTOToVO)}, nil
}

type ReadCategoryByIdInput struct {
	SobId      uuid.UUID `path:"sobId"`
	CategoryId uuid.UUID `path:"categoryId"`
}

type ReadCategoryByIdOutput struct {
	Body CategoryResponse
}

// ReadCategoryById returns one dimension category by ID.
func (h Handler) ReadCategoryById(ctx context.Context, input *ReadCategoryByIdInput) (*ReadCategoryByIdOutput, error) {
	v, err := h.app.Queries.CategoryById.Handle(ctx, input.CategoryId)
	if err != nil {
		return nil, err
	}
	if v.Id == uuid.Nil {
		return nil, huma.Error404NotFound("category not found")
	}
	return &ReadCategoryByIdOutput{Body: categoryDTOToVO(v)}, nil
}

type CreateCategoryInput struct {
	SobId uuid.UUID `path:"sobId"`
	Body  CreateCategoryRequest
}

// CreateCategory creates a dimension category in a SoB.
func (h Handler) CreateCategory(ctx context.Context, input *CreateCategoryInput) (*struct{}, error) {
	if err := h.app.Commands.CreateCategory.Handle(ctx, command.CreateCategoryCmd{
		CategoryId: uuid.New(),
		SobId:      input.SobId,
		Name:       input.Body.Name,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}

type UpdateCategoryInput struct {
	SobId      uuid.UUID `path:"sobId"`
	CategoryId uuid.UUID `path:"categoryId"`
	Body       UpdateCategoryRequest
}

// UpdateCategory renames a dimension category.
func (h Handler) UpdateCategory(ctx context.Context, input *UpdateCategoryInput) (*struct{}, error) {
	if err := h.app.Commands.UpdateCategory.Handle(ctx, command.UpdateCategoryCmd{
		CategoryId: input.CategoryId,
		NewName:    input.Body.Name,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}

type DeleteCategoryInput struct {
	SobId      uuid.UUID `path:"sobId"`
	CategoryId uuid.UUID `path:"categoryId"`
}

// DeleteCategory deletes a dimension category.
func (h Handler) DeleteCategory(ctx context.Context, input *DeleteCategoryInput) (*struct{}, error) {
	if err := h.app.Commands.DeleteCategory.Handle(ctx, command.DeleteCategoryCmd{
		CategoryId: input.CategoryId,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}

type SearchOptionsInput struct {
	SobId      uuid.UUID `path:"sobId"`
	CategoryId uuid.UUID `path:"categoryId"`
	data.PaginationInput
}

type SearchOptionsOutput struct {
	Body data.PageResponse[OptionResponse]
}

// SearchOptions lists options under a dimension category.
func (h Handler) SearchOptions(ctx context.Context, input *SearchOptionsInput) (*SearchOptionsOutput, error) {
	pageRequest, err := data.PageRequestFromInput(input.PaginationInput)
	if err != nil {
		return nil, huma.Error400BadRequest(err.Error())
	}
	page, err := h.app.Queries.PagingOptions.Handle(ctx, input.CategoryId, pageRequest)
	if err != nil {
		return nil, err
	}
	return &SearchOptionsOutput{Body: data.MapPageResponse(page, optionDTOToVO)}, nil
}

type CreateOptionInput struct {
	SobId      uuid.UUID `path:"sobId"`
	CategoryId uuid.UUID `path:"categoryId"`
	Body       CreateOptionRequest
}

// CreateOption creates an option under a dimension category.
func (h Handler) CreateOption(ctx context.Context, input *CreateOptionInput) (*struct{}, error) {
	if err := h.app.Commands.CreateOption.Handle(ctx, command.CreateOptionCmd{
		OptionId:   uuid.New(),
		CategoryId: input.CategoryId,
		Name:       input.Body.Name,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}

type UpdateOptionInput struct {
	SobId      uuid.UUID `path:"sobId"`
	CategoryId uuid.UUID `path:"categoryId"`
	OptionId   uuid.UUID `path:"optionId"`
	Body       UpdateOptionRequest
}

// UpdateOption renames a dimension option.
func (h Handler) UpdateOption(ctx context.Context, input *UpdateOptionInput) (*struct{}, error) {
	if err := h.app.Commands.UpdateOption.Handle(ctx, command.UpdateOptionCmd{
		OptionId: input.OptionId,
		NewName:  input.Body.Name,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}

type DeleteOptionInput struct {
	SobId      uuid.UUID `path:"sobId"`
	CategoryId uuid.UUID `path:"categoryId"`
	OptionId   uuid.UUID `path:"optionId"`
}

// DeleteOption deletes a dimension option.
func (h Handler) DeleteOption(ctx context.Context, input *DeleteOptionInput) (*struct{}, error) {
	if err := h.app.Commands.DeleteOption.Handle(ctx, command.DeleteOptionCmd{
		OptionId: input.OptionId,
	}); err != nil {
		return nil, err
	}
	return nil, nil
}
