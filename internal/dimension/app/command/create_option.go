package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/dimension/domain"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain/option"

	"github.com/google/uuid"
)

type CreateOptionCmd struct {
	OptionId   uuid.UUID
	CategoryId uuid.UUID
	Name       string
}

type CreateOptionHandler struct {
	repo domain.Repository
}

func NewCreateOptionHandler(repo domain.Repository) CreateOptionHandler {
	if repo == nil {
		panic("nil repo")
	}

	return CreateOptionHandler{repo: repo}
}

func (h CreateOptionHandler) Handle(ctx context.Context, cmd CreateOptionCmd) error {
	o, err := option.New(cmd.OptionId, cmd.CategoryId, cmd.Name)
	if err != nil {
		return fmt.Errorf("failed to create dimension option: %w", err)
	}

	return h.repo.CreateOption(ctx, o)
}
