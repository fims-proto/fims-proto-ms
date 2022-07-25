package command

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/user/domain"
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
	return h.repo.UpsertUser(ctx, cmd.Id, func(user *domain.User) (*domain.User, error) {
		if err := user.Update(cmd.Traits); err != nil {
			return nil, errors.Wrap(err, "failed to update user")
		}
		return user, nil
	})
}
