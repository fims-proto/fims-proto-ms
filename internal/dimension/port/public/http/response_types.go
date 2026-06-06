package http

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/dimension/app/query"

	"github.com/google/uuid"
)

// CategoryResponse is the JSON view model for a DimensionCategory.
type CategoryResponse struct {
	Id        uuid.UUID `json:"id"`
	SobId     uuid.UUID `json:"sobId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// OptionResponse is the JSON view model for a DimensionOption.
type OptionResponse struct {
	Id         uuid.UUID `json:"id"`
	CategoryId uuid.UUID `json:"categoryId"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// mappers
func categoryDTOToVO(dto query.DimensionCategory) CategoryResponse {
	return CategoryResponse{
		Id:        dto.Id,
		SobId:     dto.SobId,
		Name:      dto.Name,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

func optionDTOToVO(dto query.DimensionOption) OptionResponse {
	return OptionResponse{
		Id:         dto.Id,
		CategoryId: dto.CategoryId,
		Name:       dto.Name,
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.UpdatedAt,
	}
}
