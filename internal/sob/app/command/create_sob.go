package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/sob/app/service"
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
	repo           domain.Repository
	accountService service.AccountService
	counterService service.CounterService
}

func NewCreateSobHandler(repo domain.Repository, accountService service.AccountService, counterService service.CounterService) CreateSobHandler {
	if repo == nil {
		panic("nil repo")
	}
	if accountService == nil {
		panic("nil account service")
	}
	if counterService == nil {
		panic("nil counter service")
	}
	return CreateSobHandler{
		repo:           repo,
		accountService: accountService,
		counterService: counterService,
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

	// after sob creation, load counters and accounts
	if err := h.counterService.InitializeCounters(ctx, sob.Id()); err != nil {
		return uuid.Nil, errors.Wrapf(err, "failed to initialize counters")
	}
	if err := h.accountService.InitializeAccounts(ctx, sob.Id()); err != nil {
		return uuid.Nil, errors.Wrapf(err, "failed to initialize accounts")
	}

	return sob.Id(), nil
}
