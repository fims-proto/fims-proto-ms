package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
)

type InitializeCmd struct {
	SobId uuid.UUID
}

type InitializeHandler struct {
	repo             domain.Repository
	readModel        query.GeneralLedgerReadModel
	sobService       service.SobService
	numberingService service.NumberingService
}

func NewInitializeHandler(
	repo domain.Repository,
	readModel query.GeneralLedgerReadModel,
	sobService service.SobService,
	numberingService service.NumberingService,
) InitializeHandler {
	if repo == nil {
		panic("nil repo")
	}

	if readModel == nil {
		panic("nil read model")
	}

	if sobService == nil {
		panic("nil sob service")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return InitializeHandler{
		repo:             repo,
		readModel:        readModel,
		sobService:       sobService,
		numberingService: numberingService,
	}
}

func (h InitializeHandler) Handle(ctx context.Context, cmd InitializeCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		sob, err := h.sobService.ReadById(txCtx, cmd.SobId)
		if err != nil {
			return errors.Wrap(err, "read sob failed")
		}

		// create all accounts
		if err = initializeAccounts(txCtx, sob, h.repo); err != nil {
			return errors.Wrap(err, "failed to create accounts")
		}

		// create first period
		periodId := uuid.New()
		if err = createPeriod(txCtx, CreatePeriodCmd{
			SobId:      sob.Id,
			PeriodId:   periodId,
			FiscalYear: sob.StartingPeriodYear,
			Number:     sob.StartingPeriodMonth,
		}, h.repo, h.numberingService); err != nil {
			return errors.Wrap(err, "failed to create first period")
		}

		// create all ledgers for this period
		if err = initializeLedgers(txCtx, InitializeLedgersCmd{
			SobId:    sob.Id,
			PeriodId: periodId,
		}, h.repo, h.readModel); err != nil {
			return errors.Wrap(err, "failed to initialize ledgers")
		}

		return nil
	})
}
