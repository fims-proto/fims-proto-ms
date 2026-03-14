package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/common/utils"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PostJournalCmd struct {
	JournalId uuid.UUID
	Poster    uuid.UUID
}

type postLedgersCmd struct {
	periodId uuid.UUID
	records  []postLedgersRecordCmd
}

type postLedgersRecordCmd struct {
	accountId uuid.UUID
	amount    decimal.Decimal
}

type PostJournalHandler struct {
	repo domain.Repository
}

func NewPostJournalHandler(repo domain.Repository) PostJournalHandler {
	if repo == nil {
		panic("nil repo")
	}

	return PostJournalHandler{repo: repo}
}

func (h PostJournalHandler) Handle(ctx context.Context, cmd PostJournalCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.postJournal(txCtx, cmd)
	})
}

// postJournal updates journal and triggers ledgers posting
func (h PostJournalHandler) postJournal(ctx context.Context, cmd PostJournalCmd) error {
	return h.repo.UpdateJournalHeader(
		ctx,
		cmd.JournalId,
		func(j *journal.Journal) (*journal.Journal, error) {
			if err := j.Post(cmd.Poster); err != nil {
				return nil, err
			}

			var records []postLedgersRecordCmd
			for _, item := range j.JournalLines() {
				records = append(records, postLedgersRecordCmd{
					accountId: item.AccountId(),
					amount:    item.Amount(),
				})
			}

			if err := h.postLedgers(ctx, postLedgersCmd{
				periodId: j.PeriodId(),
				records:  records,
			}); err != nil {
				return nil, fmt.Errorf("failed to post journal to ledger: %w", err)
			}

			return j, nil
		},
	)
}

func (h PostJournalHandler) postLedgers(ctx context.Context, cmd postLedgersCmd) error {
	accountCommands := cmd.records
	for _, record := range cmd.records {
		// read all superior accounts
		superiorAccounts, err := h.repo.ReadSuperiorAccountsById(ctx, record.accountId)
		if err != nil {
			return fmt.Errorf("failed to read superior accounts: %w", err)
		}

		for _, superiorAccount := range superiorAccounts {
			superiorRecord := postLedgersRecordCmd{
				accountId: superiorAccount.Id(),
				amount:    record.amount,
			}
			accountCommands = append(accountCommands, superiorRecord)
		}
	}

	// merge same accounts
	accountsMap := utils.SliceToMapMerge(accountCommands, func(c postLedgersRecordCmd) uuid.UUID {
		return c.accountId
	}, func(c postLedgersRecordCmd) postLedgersRecordCmd {
		return c
	}, func(existing, replacement postLedgersRecordCmd) postLedgersRecordCmd {
		existing.amount = existing.amount.Add(replacement.amount)
		return existing
	})

	return h.repo.UpdateLedgersByPeriodAndAccountIds(
		ctx,
		cmd.periodId,
		utils.MapToKeySlice(accountsMap),
		func(ledgers []*ledger.Ledger) ([]*ledger.Ledger, error) {
			for _, l := range ledgers {
				record, ok := accountsMap[l.AccountId()]
				if !ok {
					return nil, fmt.Errorf("should not happen, failed to find account %s in accountsMap", l.AccountId())
				}

				l.UpdateBalance(record.amount)
			}

			return ledgers, nil
		},
	)
}
