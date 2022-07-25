package command

import (
	"context"

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
	ledgerService  service.LedgerService
}

func NewCreateSobHandler(repo domain.Repository, accountService service.AccountService, ledgerService service.LedgerService) CreateSobHandler {
	if repo == nil {
		panic("nil repo")
	}
	if accountService == nil {
		panic("nil account service")
	}
	if ledgerService == nil {
		panic("nil ledger service")
	}
	return CreateSobHandler{
		repo:           repo,
		accountService: accountService,
		ledgerService:  ledgerService,
	}
}

func (h CreateSobHandler) Handle(ctx context.Context, cmd CreateSobCmd) error {
	sob, err := domain.NewSob(cmd.SobId, cmd.Name, cmd.Description, cmd.BaseCurrency, cmd.StartingPeriodYear, cmd.StartingPeriodMonth, cmd.AccountsCodeLength)
	if err != nil {
		return errors.Wrap(err, "create sob failed")
	}
	if err := h.repo.CreateSob(ctx, sob); err != nil {
		return err
	}

	// triggers
	// create accounts
	if err = h.accountService.InitializeAccounts(ctx, cmd.SobId); err != nil {
		return errors.Wrapf(err, "failed to initialize accounts")
	}

	// create period
	if err = h.ledgerService.InitializeFirstPeriod(ctx, cmd.SobId, sob.StartingPeriodYear(), sob.StartingPeriodMonth()); err != nil {
		return errors.Wrapf(err, "failed to initialize first period")
	}

	return nil
}
