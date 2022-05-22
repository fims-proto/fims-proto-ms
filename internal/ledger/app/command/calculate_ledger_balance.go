package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// calculate one ledger balance by reading all ledger logs in current period

type CalculateLedgerBalanceCmd struct {
	Ids []uuid.UUID
}

type CalculateLedgerBalanceHandler struct {
	repo            domain.Repository
	ledgerReadModel query.LedgerReadModel
	accountService  service.AccountService
}

func NewCalculateLedgerBalanceHandler(repo domain.Repository, readModel query.LedgerReadModel, accountService service.AccountService) CalculateLedgerBalanceHandler {
	if repo == nil {
		panic("nil ledger repository")
	}
	if readModel == nil {
		panic("nil ledger read model")
	}
	if accountService == nil {
		panic("nil account service")
	}
	return CalculateLedgerBalanceHandler{
		repo:            repo,
		ledgerReadModel: readModel,
		accountService:  accountService,
	}
}

func (h CalculateLedgerBalanceHandler) Handle(ctx context.Context, cmd CalculateLedgerBalanceCmd) (err error) {
	log.Info(ctx, "handle calculating ledger balance, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle calculating ledger balance failed")
		}
	}()

	return h.repo.UpdateLedgers(ctx, cmd.Ids, func(ledgers []*domain.Ledger) ([]*domain.Ledger, error) {
		if len(ledgers) == 0 {
			return nil, errors.Errorf("no ledgers found")
		}

		// verify accounting period: should be same
		log.Info(ctx, "verify accounting period")
		periods := make(map[uuid.UUID]string, 1)
		periods[ledgers[0].PeriodId()] = "dummy"
		for _, ledger := range ledgers {
			_, ok := periods[ledger.PeriodId()]
			if !ok {
				return nil, errors.New("accounting period should be same")
			}
		}
		period, err := h.ledgerReadModel.ReadAccountingPeriodById(ctx, ledgers[0].PeriodId())
		if err != nil {
			return nil, errors.New("read accounting period failed")
		}
		if period.IsClosed {
			return nil, errors.New("accounting period is closed, ledger update not possible")
		}

		// read account
		log.Info(ctx, "read account")
		var accountIds []uuid.UUID
		for _, ledger := range ledgers {
			accountIds = append(accountIds, ledger.AccountId())
		}
		accounts, err := h.accountService.ReadAccountsByIds(ctx, accountIds)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read accounts")
		}

		// read ledger logs
		log.Info(ctx, "read ledger logs")
		ledgerLogs, err := h.ledgerReadModel.ReadLedgerLogsByAccountIdsAndTimes(ctx, accountIds, period.OpeningTime, period.EndingTime)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read ledger logs by AccountId and opening and ending time")
		}

		// ledger log calculation
		log.Info(ctx, "ledger log calculation")
		for _, ledger := range ledgers {
			var totalDebit, totalCredit decimal.Decimal
			for _, ledgerLog := range ledgerLogs {
				if ledgerLog.AccountId == ledger.AccountId() {
					totalDebit = totalDebit.Add(ledgerLog.Debit)
					totalCredit = totalCredit.Add(ledgerLog.Credit)
				}
			}

			account, ok := accounts[ledger.AccountId()]
			if !ok {
				return nil, errors.Errorf("nat able to find account by id %s, should not happen", ledger.AccountId())
			}
			ledger.UpdatePeriodAmount(totalDebit, totalCredit, account.BalanceDirection)
		}

		return ledgers, nil
	})
}
