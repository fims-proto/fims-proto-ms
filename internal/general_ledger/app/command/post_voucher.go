package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger_entry"

	"github/fims-proto/fims-proto-ms/internal/common/utils"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type PostVoucherCmd struct {
	VoucherId uuid.UUID
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

type postAuxiliaryLedgersCmd struct {
	sobId    uuid.UUID
	periodId uuid.UUID
	records  []postAuxiliaryLedgersRecordCmd
}

type postAuxiliaryLedgersRecordCmd struct {
	accountId           uuid.UUID
	auxiliaryCategoryId uuid.UUID
	auxiliaryAccount    auxiliary_account.AuxiliaryAccount
	amount              decimal.Decimal
}

type PostVoucherHandler struct {
	repo domain.Repository
}

func NewPostVoucherHandler(repo domain.Repository) PostVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}

	return PostVoucherHandler{repo: repo}
}

func (h PostVoucherHandler) Handle(ctx context.Context, cmd PostVoucherCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.postVoucher(txCtx, cmd)
	})
}

// postVoucher updates voucher, creates ledger entries, and triggers ledgers and auxiliary ledgers (if applicable) posting
func (h PostVoucherHandler) postVoucher(ctx context.Context, cmd PostVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherId,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			if err := v.Post(cmd.Poster); err != nil {
				return nil, err
			}

			if err := h.insertLedgerEntries(v, ctx); err != nil {
				return nil, fmt.Errorf("failed to insert ledger entries: %w", err)
			}

			var records []postLedgersRecordCmd
			for _, item := range v.LineItems() {
				records = append(records, postLedgersRecordCmd{
					accountId: item.AccountId(),
					amount:    item.Amount(),
				})
			}

			if err := h.postLedgers(ctx, postLedgersCmd{
				periodId: v.PeriodId(),
				records:  records,
			}); err != nil {
				return nil, fmt.Errorf("failed to post voucher to ledger: %w", err)
			}

			var auxiliaryRecords []postAuxiliaryLedgersRecordCmd
			for _, item := range v.LineItems() {
				for _, auxiliaryAccount := range item.AuxiliaryAccounts() {
					auxiliaryRecords = append(auxiliaryRecords, postAuxiliaryLedgersRecordCmd{
						accountId:           item.AccountId(),
						auxiliaryCategoryId: auxiliaryAccount.Category().Id(),
						auxiliaryAccount:    *auxiliaryAccount,
						amount:              item.Amount(),
					})
				}
			}

			if err := h.postAuxiliaryLedgers(ctx, postAuxiliaryLedgersCmd{
				sobId:    v.SobId(),
				periodId: v.PeriodId(),
				records:  auxiliaryRecords,
			}); err != nil {
				return nil, fmt.Errorf("failed to post voucher to auxiliary ledger: %w", err)
			}

			return v, nil
		},
	)
}

func (h PostVoucherHandler) insertLedgerEntries(v *voucher.Voucher, ctx context.Context) error {
	var ledgerEntries []*ledger_entry.LedgerEntry
	for _, item := range v.LineItems() {
		entry, err := ledger_entry.New(
			uuid.New(),
			v.SobId(),
			v.PeriodId(),
			v.Id(),
			item.Id(),
			item.AccountId(),
			item.AuxiliaryAccounts(),
			v.TransactionDate(),
			item.Amount(),
		)
		if err != nil {
			return fmt.Errorf("failed to create ledger entry for line item %s: %w", item.Id(), err)
		}
		ledgerEntries = append(ledgerEntries, entry)
	}

	if err := h.repo.CreateLedgerEntries(ctx, ledgerEntries); err != nil {
		return fmt.Errorf("failed to create ledger entries: %w", err)
	}
	return nil
}

func (h PostVoucherHandler) postLedgers(ctx context.Context, cmd postLedgersCmd) error {
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

func (h PostVoucherHandler) postAuxiliaryLedgers(ctx context.Context, cmd postAuxiliaryLedgersCmd) error {
	if len(cmd.records) == 0 {
		return nil
	}

	// Merge records by composite key to sum debits/credits for duplicate combinations
	auxiliaryLedgerMap := utils.SliceToMapMerge(cmd.records,
		func(c postAuxiliaryLedgersRecordCmd) domain.AuxiliaryLedgerKey {
			return domain.AuxiliaryLedgerKey{
				AccountId:           c.accountId,
				AuxiliaryCategoryId: c.auxiliaryCategoryId,
				AuxiliaryAccountId:  c.auxiliaryAccount.Id(),
			}
		},
		func(c postAuxiliaryLedgersRecordCmd) postAuxiliaryLedgersRecordCmd {
			return c
		},
		func(existing, replacement postAuxiliaryLedgersRecordCmd) postAuxiliaryLedgersRecordCmd {
			existing.amount = existing.amount.Add(replacement.amount)
			return existing
		},
	)

	requiredKeys := make([]domain.AuxiliaryLedgerKey, 0, len(auxiliaryLedgerMap))
	for key := range auxiliaryLedgerMap {
		requiredKeys = append(requiredKeys, key)
	}

	return h.repo.UpsertAuxiliaryLedgersByPeriodAndAccounts(
		ctx,
		cmd.sobId,
		cmd.periodId,
		requiredKeys,
		func(auxiliaryLedgers []*auxiliary_ledger.AuxiliaryLedger) ([]*auxiliary_ledger.AuxiliaryLedger, error) {
			ledgerMap := make(map[domain.AuxiliaryLedgerKey]*auxiliary_ledger.AuxiliaryLedger)
			for _, l := range auxiliaryLedgers {
				key := domain.AuxiliaryLedgerKey{
					AccountId:           l.AccountId(),
					AuxiliaryCategoryId: l.AuxiliaryCategoryId(),
					AuxiliaryAccountId:  l.AuxiliaryAccountId(),
				}
				ledgerMap[key] = l
			}

			for key, record := range auxiliaryLedgerMap {
				auxiliaryLedger, ok := ledgerMap[key]
				if !ok {
					// Repository guarantees all required keys exist
					return nil, fmt.Errorf(
						"auxiliary ledger not found for account=%s, category=%s, auxiliary=%s",
						record.accountId,
						record.auxiliaryCategoryId,
						record.auxiliaryAccount.Key(),
					)
				}

				auxiliaryLedger.UpdateBalance(record.amount)
			}

			return auxiliaryLedgers, nil
		},
	)
}
