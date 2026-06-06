package command

import (
	"context"
	"fmt"

	commonErrors "github/fims-proto/fims-proto-ms/internal/common/errors"
	"github/fims-proto/fims-proto-ms/internal/common/utils"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/journal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type DeleteSystemJournalCmd struct {
	SobId     uuid.UUID
	JournalId uuid.UUID
}

type DeleteSystemJournalHandler struct {
	repo domain.Repository
}

func NewDeleteSystemJournalHandler(repo domain.Repository) DeleteSystemJournalHandler {
	if repo == nil {
		panic("nil repo")
	}

	return DeleteSystemJournalHandler{repo: repo}
}

func (h DeleteSystemJournalHandler) Handle(ctx context.Context, cmd DeleteSystemJournalCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.deleteSystemJournal(txCtx, cmd)
	})
}

// deleteSystemJournal deletes a system journal (CLOSING or YEARLY_CLOSING) and reverses its ledger posts.
func (h DeleteSystemJournalHandler) deleteSystemJournal(ctx context.Context, cmd DeleteSystemJournalCmd) error {
	// 1. Read journal (header + lines + period)
	j, err := h.repo.ReadJournalById(ctx, cmd.JournalId)
	if err != nil {
		return fmt.Errorf("failed to read journal: %w", err)
	}

	// 2. Guard: sobId must match (cross-tenant access prevention)
	if j.SobId() != cmd.SobId {
		return commonErrors.NewInvalidInputError(commonErrors.SlugJournalNotFound)
	}

	// 3. Guard: only CLOSING and YEARLY_CLOSING may be deleted
	if j.JournalType() != journal.TypeClosing && j.JournalType() != journal.TypeYearlyClosing {
		return commonErrors.NewInvalidInputError(commonErrors.SlugJournalDeleteNotSystemJournal)
	}

	// 4. Reverse ledger balances (system journals are always posted, but guard defensively)
	if j.IsPosted() {
		if err = h.reverseLedgers(ctx, j); err != nil {
			return fmt.Errorf("failed to reverse ledgers: %w", err)
		}
	}

	// 5. Delete journal (header + lines + dimension options)
	return h.repo.DeleteJournalById(ctx, cmd.JournalId)
}

// reverseLedgers is the exact inverse of PostJournalHandler.postLedgers.
// It applies amount.Neg() to each affected ledger account (and its superiors).
func (h DeleteSystemJournalHandler) reverseLedgers(ctx context.Context, j *journal.Journal) error {
	type reverseRecord struct {
		accountId uuid.UUID
		amount    decimal.Decimal
	}

	// Build reversal records: negate each journal line amount
	records := make([]reverseRecord, 0, len(j.JournalLines()))
	for _, line := range j.JournalLines() {
		records = append(records, reverseRecord{
			accountId: line.AccountId(),
			amount:    line.Amount().Neg(), // ← only difference from postLedgers
		})
	}

	// Expand to superior accounts (mirrors postLedgers exactly)
	allRecords := make([]reverseRecord, len(records))
	copy(allRecords, records)
	for _, record := range records {
		superiors, err := h.repo.ReadSuperiorAccountsById(ctx, record.accountId)
		if err != nil {
			return fmt.Errorf("failed to read superior accounts: %w", err)
		}
		for _, sup := range superiors {
			allRecords = append(allRecords, reverseRecord{
				accountId: sup.Id(),
				amount:    record.amount,
			})
		}
	}

	// Merge duplicate account IDs (mirrors postLedgers)
	accountsMap := utils.SliceToMapMerge(
		allRecords,
		func(r reverseRecord) uuid.UUID { return r.accountId },
		func(r reverseRecord) reverseRecord { return r },
		func(existing, replacement reverseRecord) reverseRecord {
			existing.amount = existing.amount.Add(replacement.amount)
			return existing
		},
	)

	// Batch-update ledgers with FOR UPDATE locking
	return h.repo.UpdateLedgersByPeriodAndAccountIds(
		ctx,
		j.PeriodId(),
		utils.MapToKeySlice(accountsMap),
		func(ledgers []*ledger.Ledger) ([]*ledger.Ledger, error) {
			for _, l := range ledgers {
				record, ok := accountsMap[l.AccountId()]
				if !ok {
					return nil, fmt.Errorf(
						"should not happen: account %s not found in reversal map",
						l.AccountId(),
					)
				}
				l.UpdateBalance(record.amount)
			}
			return ledgers, nil
		},
	)
}
