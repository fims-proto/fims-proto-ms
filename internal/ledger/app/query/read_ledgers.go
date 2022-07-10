package query

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type LedgerReadModel interface {
	ReadLedgerById(ctx context.Context, id uuid.UUID) (Ledger, error)
	ReadAllLedgersByPeriod(ctx context.Context, periodId uuid.UUID, pageable data.Pageable) (data.Page[Ledger], error)
	ReadAllPeriods(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Period], error)
	ReadPeriodById(ctx context.Context, id uuid.UUID) (Period, error)
	ReadOpenPeriod(ctx context.Context, sobId uuid.UUID) (Period, error)
	ReadLedgerLogsByAccountIdsAndTimes(ctx context.Context, accountId []uuid.UUID, openingTime, endingTime time.Time) ([]LedgerLog, error)
}

type ReadLedgerHandler struct {
	readModel      LedgerReadModel
	accountService AccountService
}

func NewReadLedgerHandler(readModel LedgerReadModel, accountService AccountService) ReadLedgerHandler {
	if readModel == nil {
		panic("nil read model")
	}
	if accountService == nil {
		panic("nil account service")
	}
	return ReadLedgerHandler{
		readModel:      readModel,
		accountService: accountService,
	}
}

func (h ReadLedgerHandler) HandleReadAllPeriods(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Period], error) {
	return h.readModel.ReadAllPeriods(ctx, sobId, pageable)
}

func (h ReadLedgerHandler) HandleReadOpenPeriod(ctx context.Context, sobId uuid.UUID) (Period, error) {
	return h.readModel.ReadOpenPeriod(ctx, sobId)
}

func (h ReadLedgerHandler) HandleReadPeriodById(ctx context.Context, id uuid.UUID) (Period, error) {
	return h.readModel.ReadPeriodById(ctx, id)
}

func (h ReadLedgerHandler) HandleReadAllLedgersByPeriod(ctx context.Context, periodId uuid.UUID, pageable data.Pageable) (data.Page[Ledger], error) {
	ledgersPage, err := h.readModel.ReadAllLedgersByPeriod(ctx, periodId, pageable)
	if err != nil {
		return nil, errors.Wrap(err, "failed on reading ledgers by period")
	}

	var accountIds []uuid.UUID
	for _, ledger := range ledgersPage.Content() {
		accountIds = append(accountIds, ledger.Account.Id)
	}

	accounts, err := h.accountService.ReadAccountsByIds(ctx, accountIds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read account by ids")
	}
	for i := range ledgersPage.Content() {
		account := accounts[ledgersPage.Content()[i].Account.Id]

		ledgersPage.Content()[i].Account.SuperiorAccountId = account.SuperiorAccountId
		ledgersPage.Content()[i].Account.AccountNumber = account.AccountNumber
		ledgersPage.Content()[i].Account.Title = account.Title
		ledgersPage.Content()[i].Account.AccountType = account.AccountType
		ledgersPage.Content()[i].Account.BalanceDirection = account.BalanceDirection
	}
	return ledgersPage, nil
}
