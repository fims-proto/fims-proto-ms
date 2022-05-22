package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UpdateSobCmd struct {
	Id                 uuid.UUID
	Name               string
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

func (h UpdateSobHandler) Handle(ctx context.Context, cmd UpdateSobCmd) (err error) {
	log.Info(ctx, "handle updating sob")
	log.Debug(ctx, "handle updating sob, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle updating failed")
		}
	}()

	return h.repo.UpdateSob(
		ctx,
		cmd.Id,
		func(s *domain.Sob) (*domain.Sob, error) {
			if cmd.Name != "" {
				if err := s.UpdateName(cmd.Name); err != nil {
					return nil, errors.Wrap(err, "sob name updating failed")
				}
			}
			if cmd.AccountsCodeLength != nil {
				if err := s.UpdateAccountsCodeLength(cmd.AccountsCodeLength); err != nil {
					return nil, errors.Wrap(err, "sob accounts code length updating failed")
				}
			}
			return s, nil
		},
	)
}
