package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/common/log"
	"github/fims-proto/fims-proto-ms/internal/ledger/app/service"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// append new logs in to ledger log. given list of account Ids, get all superior Ids, and append into logs table

type AppendLedgerLogCmd struct {
	PostingId       uuid.UUID
	AccountId       uuid.UUID
	VoucherId       uuid.UUID
	TransactionTime time.Time
	Debit           decimal.Decimal
	Credit          decimal.Decimal
}

type AppendLedgerLogsHandler struct {
	repo           domain.Repository
	accountService service.AccountService
}

func NewAppendLedgerLogsHandler(repo domain.Repository, accountService service.AccountService) AppendLedgerLogsHandler {
	if repo == nil {
		panic("nil ledger repo")
	}
	if accountService == nil {
		panic("nil account service")
	}
	return AppendLedgerLogsHandler{
		repo:           repo,
		accountService: accountService,
	}
}

func (h AppendLedgerLogsHandler) Handle(ctx context.Context, commands []AppendLedgerLogCmd) (err error) {
	log.Info(ctx, "handle append ledger log")
	log.Debug(ctx, "handle append ledger log, commands: %+v", commands)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle calculating ledger balance failed")
		}
	}()

	var ledgerLogs []*domain.LedgerLog
	for _, cmd := range commands {
		ledgerLog, err := domain.NewLedgerLog(uuid.New(), cmd.PostingId, cmd.AccountId, cmd.VoucherId, cmd.TransactionTime, cmd.Debit, cmd.Credit)
		if err != nil {
			return errors.Wrap(err, "create ledger log domain model failed")
		}
		ledgerLogs = append(ledgerLogs, ledgerLog)

		// for superior account
		superiorAccountIds, err := h.accountService.ReadSuperiorAccountIds(ctx, cmd.AccountId)
		if err != nil {
			return errors.Wrapf(err, "failed to read superior account id of %s", cmd.AccountId)
		}
		for _, accountId := range superiorAccountIds {
			ledgerLog, err = domain.NewLedgerLog(uuid.New(), cmd.PostingId, accountId, cmd.VoucherId, cmd.TransactionTime, cmd.Debit, cmd.Credit)
			if err != nil {
				return errors.Wrap(err, "create ledger log domain model failed")
			}
			ledgerLogs = append(ledgerLogs, ledgerLog)
		}
	}

	return h.repo.CreateLedgerLogs(ctx, ledgerLogs)
}
