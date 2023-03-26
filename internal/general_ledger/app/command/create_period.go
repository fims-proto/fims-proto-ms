package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/period"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type CreatePeriodCmd struct {
	SobId      uuid.UUID
	PeriodId   uuid.UUID
	FiscalYear int
	Number     int
}

type CreatePeriodHandler struct {
	repo             domain.Repository
	numberingService service.NumberingService
	readModel        query.GeneralLedgerReadModel
}

func NewCreatePeriodHandler(
	repo domain.Repository,
	numberingService service.NumberingService,
	readModel query.GeneralLedgerReadModel,
) CreatePeriodHandler {
	if repo == nil {
		panic("nil account repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	if readModel == nil {
		panic("nil read model")
	}

	return CreatePeriodHandler{
		repo:             repo,
		numberingService: numberingService,
		readModel:        readModel,
	}
}

func (h CreatePeriodHandler) Handle(ctx context.Context, cmd CreatePeriodCmd) error {
	return createPeriod(ctx, cmd, h.repo, h.numberingService)
}

func createPeriod(
	ctx context.Context,
	cmd CreatePeriodCmd,
	repo domain.Repository,
	numberingService service.NumberingService,
) error {
	p, err := period.New(cmd.PeriodId, cmd.SobId, cmd.FiscalYear, cmd.Number)
	if err != nil {
		return errors.Wrap(err, "failed to create period domain model")
	}

	// create numbering configuration for voucher entries in this period
	if err = numberingService.CreateIdentifierConfigurationForVoucher(ctx, cmd.PeriodId); err != nil {
		return errors.Wrap(err, "failed to create numbering configuration for period")
	}

	return repo.CreatePeriod(ctx, p)
}
