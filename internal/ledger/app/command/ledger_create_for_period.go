package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/ledger/app/service"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// create ledgers for given period
// using ending balance from previous period

type CreatePeriodLedgersCmd struct {
	PeriodId uuid.UUID
}

type CreatePeriodLedgersHandler struct {
	repo           domain.Repository
	readModel      query.LedgerReadModel
	accountService service.AccountService
}

func NewCreatePeriodLedgersHandler(repo domain.Repository, readModel query.LedgerReadModel, accountService service.AccountService) CreatePeriodLedgersHandler {
	if repo == nil {
		panic("nil ledger repository")
	}
	if readModel == nil {
		panic("nil ledger read model")
	}
	if accountService == nil {
		panic("nil account service")
	}
	return CreatePeriodLedgersHandler{
		repo:           repo,
		readModel:      readModel,
		accountService: accountService,
	}
}

func (h CreatePeriodLedgersHandler) Handle(ctx context.Context, cmd CreatePeriodLedgersCmd) (err error) {
	log.Info(ctx, "handle create ledgers for period, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle create ledgers for period")
		}
	}()

	period, err := h.readModel.ReadPeriodById(ctx, cmd.PeriodId)
	if err != nil {
		return errors.Wrap(err, "failed to read period")
	}

	// read all accounts in sob
	accounts, err := h.accountService.ReadAccountsBySobId(ctx, period.SobId)
	if err != nil {
		return errors.Wrap(err, "failed to read all account Ids by SoB")
	}

	// read all ledgers from previous period
	previousLedgers := make(map[uuid.UUID]query.Ledger) // key: AccountId, value: ledger
	if period.PreviousPeriodId != uuid.Nil {
		ledgersPage, err := h.readModel.ReadLedgersByPeriod(ctx, period.PreviousPeriodId, data.Unpaged())
		if err != nil {
			return errors.Wrap(err, "failed to read ledgers by account period")
		}
		for _, ledger := range ledgersPage.Content() {
			previousLedgers[ledger.Account.Id] = ledger
		}
	}

	// prepare ledgers in this period
	var ledgers []*domain.Ledger
	for _, account := range accounts {
		// use ending balance from previous period as opening balance if available
		openingBalance := decimal.Zero
		previousLedger, ok := previousLedgers[account.Id]
		if ok {
			// TODO how to ensure the endingBalance is up-to-date?
			openingBalance = previousLedger.EndingBalance
		}

		ledger, err := domain.NewLedger(uuid.New(), cmd.PeriodId, account.Id, account.AccountNumber, openingBalance, decimal.Zero, decimal.Zero, decimal.Zero)
		if err != nil {
			return errors.Wrap(err, "failed to create ledger domain model")
		}

		ledgers = append(ledgers, ledger)
	}

	// save
	return h.repo.CreateLedgers(ctx, ledgers)
}
