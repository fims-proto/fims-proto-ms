package http

import (
	"time"

	"github.com/google/uuid"
	"github/fims-proto/fims-proto-ms/internal/user/app/query"
)

type Error struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}

type UserResponse struct {
	Id        uuid.UUID `json:"id"`
	Traits    any       `json:"traits"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateUserRequest struct {
	Traits string `json:"traits"`
}

func userDTOToVO(u query.User) UserResponse {
	return UserResponse{
		Id:        u.Id,
		Traits:    u.Traits,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
