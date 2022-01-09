package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/sob/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateSobCmd struct {
	Name                string
	Description         string
	BaseCurrency        string
	StartingPeriodYear  int
	StartingPeriodMonth int
	AccountsCodeLength  []int
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

func (h CreateSobHandler) Handle(ctx context.Context, cmd CreateSobCmd) (createdId uuid.UUID, err error) {
	log.Info(ctx, "handle creating sob")
	log.Debug(ctx, "handle creating sob, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle creating sob failed")
		}
	}()

	sob, err := domain.NewSob(uuid.New(), cmd.Name, cmd.Description, cmd.BaseCurrency, cmd.StartingPeriodYear, cmd.StartingPeriodMonth, cmd.AccountsCodeLength)
	if err != nil {
		return uuid.Nil, errors.Wrapf(err, "create sob failed")
	}
	if err := h.repo.CreateSob(ctx, sob); err != nil {
		return uuid.Nil, err
	}
	return sob.Id(), nil
}
