package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/sob/domain/sob"

	"github/fims-proto/fims-proto-ms/internal/sob/domain"

	"github.com/google/uuid"
)

type UpdateSobCmd struct {
	SobId              uuid.UUID
	Name               string
	Description        *string
	AccountsCodeLength []int
}

type UpdateSobHandler struct {
	repo domain.Repository
}

func NewUpdateSobHandler(repo domain.Repository) UpdateSobHandler {
	if repo == nil {
		panic("nil repo")
	}
	return UpdateSobHandler{repo: repo}
}

func (h UpdateSobHandler) Handle(ctx context.Context, cmd UpdateSobCmd) error {
	return h.repo.UpdateSob(
		ctx,
		cmd.SobId,
		func(s *sob.Sob) (*sob.Sob, error) {
			if cmd.Name != "" {
				if err := s.UpdateName(cmd.Name); err != nil {
					return nil, fmt.Errorf("failed to update sob name: %w", err)
				}
			}

			if cmd.Description != nil {
				if err := s.UpdateDescription(*cmd.Description); err != nil {
					return nil, fmt.Errorf("failed to update sob description: %w", err)
				}
			}

			if cmd.AccountsCodeLength != nil {
				if err := s.UpdateAccountsCodeLength(cmd.AccountsCodeLength); err != nil {
					return nil, fmt.Errorf("failed to update sob account code length: %w", err)
				}
			}
			return s, nil
		},
	)
}
