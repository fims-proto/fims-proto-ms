package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/sob/domain/sob"

	"github/fims-proto/fims-proto-ms/internal/sob/app/service"

	"github/fims-proto/fims-proto-ms/internal/sob/domain"

	"github.com/google/uuid"
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
	repo                 domain.Repository
	generalLedgerService service.GeneralLedgerService
}

func NewCreateSobHandler(repo domain.Repository, generalLedgerService service.GeneralLedgerService) CreateSobHandler {
	if repo == nil {
		panic("nil repo")
	}

	if generalLedgerService == nil {
		panic("nil account service")
	}

	return CreateSobHandler{
		repo:                 repo,
		generalLedgerService: generalLedgerService,
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
		return fmt.Errorf("failed to create sob: %w", err)
	}

	if err = h.repo.CreateSob(ctx, sobBO); err != nil {
		return err
	}

	// initialize general ledger for sob
	if err = h.generalLedgerService.InitializeForSob(ctx, cmd.SobId); err != nil {
		return fmt.Errorf("failed to initialzie account: %w", err)
	}

	return nil
}
