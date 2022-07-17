package query

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type LedgerReadModel interface {
	ReadPeriods(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[Period], error)
	ReadPeriodsByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]Period, error)
	ReadPeriodById(ctx context.Context, id uuid.UUID) (Period, error)
	ReadPeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (Period, error)
	ReadOpenPeriod(ctx context.Context, sobId uuid.UUID) (Period, error)

	ReadLedgersByPeriod(ctx context.Context, periodId uuid.UUID, pageable data.Pageable) (data.Page[Ledger], error)
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
	return h.readModel.ReadPeriods(ctx, sobId, pageable)
}

func (h ReadLedgerHandler) HandleReadOpenPeriod(ctx context.Context, sobId uuid.UUID) (Period, error) {
	return h.readModel.ReadOpenPeriod(ctx, sobId)
}

func (h ReadLedgerHandler) HandleReadPeriodById(ctx context.Context, id uuid.UUID) (Period, error) {
	return h.readModel.ReadPeriodById(ctx, id)
}

func (h ReadLedgerHandler) HandleReadPeriodsByIds(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]Period, error) {
	return h.readModel.ReadPeriodsByIds(ctx, ids)
}

func (h ReadLedgerHandler) HandleReadPeriodByTime(ctx context.Context, sobId uuid.UUID, timePoint time.Time) (Period, error) {
	return h.readModel.ReadPeriodByTime(ctx, sobId, timePoint)
}

func (h ReadLedgerHandler) HandleReadAllLedgersByPeriod(ctx context.Context, periodId uuid.UUID, pageable data.Pageable) (data.Page[Ledger], error) {
	ledgersPage, err := h.readModel.ReadLedgersByPeriod(ctx, periodId, pageable)
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

		ledgersPage.Content()[i].Account.AccountNumber = account.AccountNumber
		ledgersPage.Content()[i].Account.Title = account.Title
		ledgersPage.Content()[i].Account.AccountType = account.AccountType
		ledgersPage.Content()[i].Account.BalanceDirection = account.BalanceDirection
	}
	return ledgersPage, nil
}
