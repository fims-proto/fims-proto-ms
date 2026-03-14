package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/dimension/domain"
	"github/fims-proto/fims-proto-ms/internal/dimension/domain/option"

	"github.com/google/uuid"
)

type UpdateOptionCmd struct {
	OptionId uuid.UUID
	NewName  string
}

type UpdateOptionHandler struct {
	repo domain.Repository
}

func NewUpdateOptionHandler(repo domain.Repository) UpdateOptionHandler {
	if repo == nil {
		panic("nil repo")
	}

	return UpdateOptionHandler{repo: repo}
}

func (h UpdateOptionHandler) Handle(ctx context.Context, cmd UpdateOptionCmd) error {
	return h.repo.UpdateOption(ctx, cmd.OptionId, func(o *option.DimensionOption) (*option.DimensionOption, error) {
		if err := o.Rename(cmd.NewName); err != nil {
			return nil, fmt.Errorf("failed to rename dimension option: %w", err)
		}

		return o, nil
	})
}
