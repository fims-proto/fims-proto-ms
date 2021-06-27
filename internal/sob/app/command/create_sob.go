package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"

	"github.com/pkg/errors"
)

type CreateSobCmd struct {
	Id          string
	Name        string
	Description string
}

type CreateSobHandler struct {
	repo domain.Repository
}

func NewCreateSobHandler(repo domain.Repository) CreateSobHandler {
	if repo == nil {
		panic("nil repo")
	}
	return CreateSobHandler{
		repo: repo,
	}
}

func (h CreateSobHandler) Handle(ctx context.Context, cmd CreateSobCmd) error {
	sob, err := domain.NewSob(cmd.Id, cmd.Name, cmd.Description)
	if err != nil {
		return errors.Wrapf(err, "create sob failed")
	}
	return h.repo.CreateSob(ctx, sob)
}
