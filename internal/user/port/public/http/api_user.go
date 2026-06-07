package http

import (
	"context"
	"encoding/json"

	"github/fims-proto/fims-proto-ms/internal/user/app/command"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type ReadUserByIdInput struct {
	UserId uuid.UUID `path:"userId"`
}

type ReadUserByIdOutput struct {
	Body UserResponse
}

// ReadUserById returns one user by ID.
func (h Handler) ReadUserById(ctx context.Context, input *ReadUserByIdInput) (*ReadUserByIdOutput, error) {
	user, err := h.app.Queries.UserById.Handle(ctx, input.UserId)
	if err != nil {
		return nil, err
	}
	if user.Id == uuid.Nil {
		return nil, huma.Error404NotFound("user not found")
	}
	return &ReadUserByIdOutput{Body: userDTOToVO(user)}, nil
}

type UpdateUserInput struct {
	UserId uuid.UUID `path:"userId"`
	Body   UpdateUserRequest
}

// UpdateUser updates user traits JSON.
func (h Handler) UpdateUser(ctx context.Context, input *UpdateUserInput) (*struct{}, error) {
	var traits json.RawMessage
	if err := traits.UnmarshalJSON([]byte(input.Body.Traits)); err != nil {
		return nil, huma.Error400BadRequest("invalid traits JSON")
	}
	cmd := command.UpdateUserCmd{
		Id:     input.UserId,
		Traits: traits,
	}
	if err := h.app.Commands.UpdateUser.Handle(ctx, cmd); err != nil {
		return nil, err
	}
	return nil, nil
}
