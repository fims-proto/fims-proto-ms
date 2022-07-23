package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/ledger/app/service"

	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// create first period in the SoB
// create period, checking number, using ending time as opening time

type CreatePeriodCmd struct {
	PreviousPeriodId uuid.UUID
	SobId            uuid.UUID
	FinancialYear    int
	Number           int
	OpeningTime      time.Time
	EndingTime       time.Time
}

type CreatePeriodHandler struct {
	repo             domain.Repository
	readModel        query.LedgerReadModel
	selfService      service.SelfService
	numberingService service.NumberingService
}

func NewCreatePeriodHandler(repo domain.Repository, readModel query.LedgerReadModel, selfService service.SelfService, numberingService service.NumberingService) CreatePeriodHandler {
	if repo == nil {
		panic("nil ledger repo")
	}
	if readModel == nil {
		panic("nil ledger read model")
	}
	if selfService == nil {
		panic("nil ledger self service")
	}
	if numberingService == nil {
		panic("nil numbering service")
	}
	return CreatePeriodHandler{
		repo:             repo,
		readModel:        readModel,
		selfService:      selfService,
		numberingService: numberingService,
	}
}

func (h CreatePeriodHandler) Handle(ctx context.Context, cmd CreatePeriodCmd) (uuid.UUID, error) {
	// use previous period ending time as new opening time if previous period provided
	// otherwise using given opening time
	openingTime := cmd.OpeningTime
	if cmd.PreviousPeriodId != uuid.Nil {
		previousPeriod, err := h.readModel.ReadPeriodById(ctx, cmd.PreviousPeriodId)
		if err != nil {
			return uuid.Nil, errors.Wrap(err, "failed to read previous period")
		}
		if previousPeriod.SobId != cmd.SobId {
			return uuid.Nil, errors.Wrap(err, "sob id not equals to the one from previous period")
		}
		if !previousPeriod.IsClosed {
			return uuid.Nil, errors.Wrap(err, "previous period not closed")
		}
		openingTime = previousPeriod.EndingTime
	}

	createdId := uuid.New()
	period, err := domain.NewPeriod(createdId, cmd.SobId, cmd.PreviousPeriodId, cmd.FinancialYear, cmd.Number, openingTime, cmd.EndingTime, false)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to create period domain model")
	}
	if err = h.repo.CreatePeriod(ctx, period); err != nil {
		return uuid.Nil, err
	}

	// create ledgers for this period
	if err = h.selfService.CreateLedgersForPeriod(ctx, createdId); err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to create ledgers for period")
	}

	// create numbering configuration for vouchers in this period
	if err = h.numberingService.InitializeIdentifierConfigurationForVoucher(ctx, createdId); err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to create numbering configuration for period")
	}

	return createdId, nil
}
