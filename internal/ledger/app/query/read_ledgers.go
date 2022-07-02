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
	ReadAllLedgersByAccountingPeriod(ctx context.Context, periodId uuid.UUID, pageable data.Pageable) (data.Page[Ledger], error)
	ReadAllAccountingPeriods(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[AccountingPeriod], error)
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

func (h ReadLedgerHandler) HandleReadAllAccountingPeriods(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[AccountingPeriod], error) {
	return h.readModel.ReadAllAccountingPeriods(ctx, sobId, pageable)
}

func (h ReadLedgerHandler) HandleReadOpenAccountingPeriod(ctx context.Context, sobId uuid.UUID) (AccountingPeriod, error) {
	return h.readModel.ReadOpenAccountingPeriod(ctx, sobId)
}

func (h ReadLedgerHandler) HandleReadAccountingPeriodById(ctx context.Context, id uuid.UUID) (AccountingPeriod, error) {
	return h.readModel.ReadAccountingPeriodById(ctx, id)
}

func (h ReadLedgerHandler) HandleReadAllLedgersByAccountingPeriod(ctx context.Context, periodId uuid.UUID, pageable data.Pageable) (data.Page[Ledger], error) {
	ledgersPage, err := h.readModel.ReadAllLedgersByAccountingPeriod(ctx, periodId, pageable)
	if err != nil {
		return data.Page[Ledger]{}, errors.Wrap(err, "failed on reading ledgers by period")
	}

	var accountIds []uuid.UUID
	for _, ledger := range ledgersPage.Content {
		accountIds = append(accountIds, ledger.Account.Id)
	}

	accounts, err := h.accountService.ReadAccountsByIds(ctx, accountIds)
	if err != nil {
		return data.Page[Ledger]{}, errors.Wrap(err, "failed to read account by ids")
	}
	for i := range ledgersPage.Content {
		account := accounts[ledgersPage.Content[i].Account.Id]

		ledgersPage.Content[i].Account.SuperiorAccountId = account.SuperiorAccountId
		ledgersPage.Content[i].Account.AccountNumber = account.AccountNumber
		ledgersPage.Content[i].Account.Title = account.Title
		ledgersPage.Content[i].Account.AccountType = account.AccountType
		ledgersPage.Content[i].Account.BalanceDirection = account.BalanceDirection
	}
	return ledgersPage, nil
}
