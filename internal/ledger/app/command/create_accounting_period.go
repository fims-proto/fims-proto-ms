package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// create first accounting period in the SoB
// create accounting period, checking number, using ending time as opening time

type CreatePeriodCmd struct {
	PreviousPeriodId uuid.UUID
	SobId            uuid.UUID
	FinancialYear    int
	Number           int
	OpeningTime      time.Time
	EndingTime       time.Time
}

type CreateAccountingPeriodHandler struct {
	repo            domain.Repository
	ledgerReadModel query.LedgerReadModel
}

func NewCreateAccountingPeriodHandler(repo domain.Repository, readModel query.LedgerReadModel) CreateAccountingPeriodHandler {
	if repo == nil {
		panic("nil ledger repo")
	}
	if readModel == nil {
		panic("nil ledger read model")
	}
	return CreateAccountingPeriodHandler{
		repo:            repo,
		ledgerReadModel: readModel,
	}
}

func (h CreateAccountingPeriodHandler) Handle(ctx context.Context, cmd CreatePeriodCmd) (createdId uuid.UUID, err error) {
	log.Info(ctx, "handle initial accounting period creation, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle initial accounting period creation failed")
		}
	}()

	// use previous period ending time as new opening time if previous period provided
	// otherwise using given opening time
	openingTime := cmd.OpeningTime
	if cmd.PreviousPeriodId != uuid.Nil {
		previousAccountingPeriod, err := h.ledgerReadModel.ReadAccountingPeriodById(ctx, cmd.PreviousPeriodId)
		if err != nil {
			return uuid.Nil, errors.Wrap(err, "failed to read previous accounting period")
		}
		if previousAccountingPeriod.SobId != cmd.SobId {
			return uuid.Nil, errors.Wrap(err, "sob id not equals to the one from previous accounting period")
		}
		if !previousAccountingPeriod.IsClosed {
			return uuid.Nil, errors.Wrap(err, "previous accounting period not closed")
		}
		openingTime = previousAccountingPeriod.EndingTime
	}

	createdId = uuid.New()
	accountingPeriod, err := domain.NewAccountingPeriod(createdId, cmd.SobId, cmd.PreviousPeriodId, cmd.FinancialYear, cmd.Number, openingTime, cmd.EndingTime, false)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "failed to create accounting period domain model")
	}
	if err := h.repo.CreateAccountingPeriod(ctx, accountingPeriod); err != nil {
		return uuid.Nil, err
	}
	return createdId, nil

	// TODO create ledgers for new accounting period
}
