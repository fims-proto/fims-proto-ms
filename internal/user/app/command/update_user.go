package command

import (
	"context"
	"encoding/json"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/user/domain/user"

	"github/fims-proto/fims-proto-ms/internal/user/domain"

	"github.com/google/uuid"
)

type UpdateUserCmd struct {
	Id     uuid.UUID
	Traits json.RawMessage
}

type UpdateUserHandler struct {
	repo domain.Repository
}

func NewUpdateUserHandler(repo domain.Repository) UpdateUserHandler {
	if repo == nil {
		panic("nil repository")
	}
	return UpdateUserHandler{repo: repo}
}

func (h UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUserCmd) error {
	return h.repo.UpsertUser(ctx, cmd.Id, func(user *user.User) (*user.User, error) {
		if err := user.Update(cmd.Traits); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
		return user, nil
	})
}
