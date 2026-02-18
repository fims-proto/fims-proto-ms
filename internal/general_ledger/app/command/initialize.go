package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
)

type InitializeCmd struct {
	SobId uuid.UUID
}

type InitializeHandler struct {
	repo             domain.Repository
	sobService       service.SobService
	numberingService service.NumberingService
}

func NewInitializeHandler(repo domain.Repository, sobService service.SobService, numberingService service.NumberingService) InitializeHandler {
	if repo == nil {
		panic("nil repo")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return InitializeHandler{
		repo:             repo,
		sobService:       sobService,
		numberingService: numberingService,
	}
}

func (h InitializeHandler) Handle(ctx context.Context, cmd InitializeCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		sob, err := h.sobService.ReadById(txCtx, cmd.SobId)
		if err != nil {
			return fmt.Errorf("failed to read sob: %w", err)
		}

		// create all accounts
		if err = initializeAccounts(txCtx, sob, h.repo); err != nil {
			return fmt.Errorf("failed to create accounts: %w", err)
		}

		// create first period
		periodId := uuid.New()
		if err = initializePeriod(txCtx, initializePeriodCmd{
			SobId:      sob.Id,
			PeriodId:   periodId,
			FiscalYear: sob.StartingPeriodYear,
			Number:     sob.StartingPeriodMonth,
		}, h.repo, h.numberingService); err != nil {
			return fmt.Errorf("failed to create first period: %w", err)
		}

		// create all ledgers for this period
		if err = initializeAllLedgers(txCtx, h.repo, sob.Id); err != nil {
			return fmt.Errorf("failed to initialize ledgers: %w", err)
		}

		return nil
	})
}
