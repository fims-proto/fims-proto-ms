package query

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/google/uuid"
)

type LedgerReadModel interface {
	ReadLedgerById(ctx context.Context, id uuid.UUID) (Ledger, error)
	ReadAllLedgersByAccountingPeriod(ctx context.Context, periodId uuid.UUID) ([]Ledger, error)
	ReadAllAccountingPeriods(ctx context.Context, sobId uuid.UUID) ([]AccountingPeriod, error)
	ReadAccountingPeriodById(ctx context.Context, id uuid.UUID) (AccountingPeriod, error)
	ReadOpenAccountingPeriod(ctx context.Context, sobId uuid.UUID) (AccountingPeriod, error)
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

func (h ReadLedgerHandler) HandleReadAllAccountingPeriods(ctx context.Context, sobId uuid.UUID) ([]AccountingPeriod, error) {
	return h.readModel.ReadAllAccountingPeriods(ctx, sobId)
}

func (h ReadLedgerHandler) HandleReadOpenAccountingPeriod(ctx context.Context, sobId uuid.UUID) (AccountingPeriod, error) {
	return h.readModel.ReadOpenAccountingPeriod(ctx, sobId)
}

func (h ReadLedgerHandler) HandleReadAccountingPeriodById(ctx context.Context, id uuid.UUID) (AccountingPeriod, error) {
	return h.readModel.ReadAccountingPeriodById(ctx, id)
}

func (h ReadLedgerHandler) HandleReadAllLedgersByAccountingPeriod(ctx context.Context, periodId uuid.UUID) ([]Ledger, error) {
	ledgers, err := h.readModel.ReadAllLedgersByAccountingPeriod(ctx, periodId)
	if err != nil {
		return nil, errors.Wrap(err, "failed on reading ledgers by period")
	}

	var accountIds []uuid.UUID
	for _, ledger := range ledgers {
		accountIds = append(accountIds, ledger.Account.Id)
	}

	accounts, err := h.accountService.ReadAccountsByIds(ctx, accountIds)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read account by ids")
	}
	for i := range ledgers {
		account := accounts[ledgers[i].Account.Id]

		ledgers[i].Account.SuperiorAccountId = account.SuperiorAccountId
		ledgers[i].Account.AccountNumber = account.AccountNumber
		ledgers[i].Account.Title = account.Title
		ledgers[i].Account.Level = account.Level
		ledgers[i].Account.AccountType = account.AccountType
		ledgers[i].Account.BalanceDirection = account.BalanceDirection
	}
	return ledgers, nil
}
