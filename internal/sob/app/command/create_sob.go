package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/sob/domain/sob"

	"github/fims-proto/fims-proto-ms/internal/sob/app/service"

	"github/fims-proto/fims-proto-ms/internal/sob/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreateSobCmd struct {
	SobId               uuid.UUID
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
}

func NewCreateSobHandler(repo domain.Repository, accountService service.AccountService) CreateSobHandler {
	if repo == nil {
		panic("nil repo")
	}

	if accountService == nil {
		panic("nil account service")
	}

	return CreateSobHandler{
		repo:           repo,
		accountService: accountService,
	}
}

func (h CreateSobHandler) Handle(ctx context.Context, cmd CreateSobCmd) error {
	sobBO, err := sob.New(
		cmd.SobId,
		cmd.Name,
		cmd.Description,
		cmd.BaseCurrency,
		cmd.StartingPeriodYear,
		cmd.StartingPeriodMonth,
		cmd.AccountsCodeLength,
	)
	if err != nil {
		return errors.Wrap(err, "create sob failed")
	}

	if err = h.repo.CreateSob(ctx, sobBO); err != nil {
		return err
	}

	// triggers
	// create accounts
	if err = h.accountService.InitializeAccounts(ctx, cmd.SobId); err != nil {
		return errors.Wrapf(err, "failed to initialize accounts")
	}

	// create period
	if err = h.accountService.InitializeFirstPeriod(ctx, cmd.SobId, sobBO.StartingPeriodYear(), sobBO.StartingPeriodMonth()); err != nil {
		return errors.Wrapf(err, "failed to initialize first period")
	}

	return nil
}
