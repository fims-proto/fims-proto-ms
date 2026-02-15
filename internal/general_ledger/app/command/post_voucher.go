package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_ledger"

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
	debit     decimal.Decimal
	credit    decimal.Decimal
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
	debit               decimal.Decimal
	credit              decimal.Decimal
}

// auxiliaryLedgerKey represents the composite natural key for auxiliary ledgers
type auxiliaryLedgerKey struct {
	accountId           uuid.UUID
	auxiliaryCategoryId uuid.UUID
	auxiliaryAccountId  uuid.UUID
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

// postVoucher updates voucher, and triggers ledgers and auxiliary ledgers (if applicable) posting
func (h PostVoucherHandler) postVoucher(ctx context.Context, cmd PostVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherId,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			if err := v.Post(cmd.Poster); err != nil {
				return nil, err
			}

			var records []postLedgersRecordCmd
			for _, item := range v.LineItems() {
				records = append(records, postLedgersRecordCmd{
					accountId: item.AccountId(),
					debit:     item.Debit(),
					credit:    item.Credit(),
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
						debit:               item.Debit(),
						credit:              item.Credit(),
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
				debit:     record.debit,
				credit:    record.credit,
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
		existing.debit = existing.debit.Add(replacement.debit)
		existing.credit = existing.credit.Add(replacement.credit)
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

				l.UpdateEndingBalance(record.debit, record.credit)
			}

			return ledgers, nil
		},
	)
}

func (h PostVoucherHandler) postAuxiliaryLedgers(ctx context.Context, cmd postAuxiliaryLedgersCmd) error {
	if len(cmd.records) == 0 {
		return nil
	}

	// Group records by (accountId, auxiliaryCategoryId, auxiliaryAccountId) composite key
	auxiliaryLedgerMap := utils.SliceToMapMerge(cmd.records,
		// Key function: composite key
		func(c postAuxiliaryLedgersRecordCmd) auxiliaryLedgerKey {
			return auxiliaryLedgerKey{
				accountId:           c.accountId,
				auxiliaryCategoryId: c.auxiliaryCategoryId,
				auxiliaryAccountId:  c.auxiliaryAccount.Id(),
			}
		},
		// Value function: pass through
		func(c postAuxiliaryLedgersRecordCmd) postAuxiliaryLedgersRecordCmd {
			return c
		},
		// Merge function: sum debits and credits for same composite key
		func(existing, replacement postAuxiliaryLedgersRecordCmd) postAuxiliaryLedgersRecordCmd {
			existing.debit = existing.debit.Add(replacement.debit)
			existing.credit = existing.credit.Add(replacement.credit)
			return existing
		},
	)

	// Group by account
	recordsByAccount := make(map[uuid.UUID][]postAuxiliaryLedgersRecordCmd)
	for _, record := range auxiliaryLedgerMap {
		recordsByAccount[record.accountId] = append(recordsByAccount[record.accountId], record)
	}

	// Update each account's auxiliary ledgers
	for accountId, records := range recordsByAccount {
		// Extract category and auxiliary account IDs for this account
		var categoryIds []uuid.UUID
		var auxiliaryAccountIds []uuid.UUID
		for _, r := range records {
			categoryIds = append(categoryIds, r.auxiliaryCategoryId)
			auxiliaryAccountIds = append(auxiliaryAccountIds, r.auxiliaryAccount.Id())
		}

		if err := h.repo.UpdateAuxiliaryLedgersByPeriodAndAccounts(
			ctx,
			cmd.periodId,
			accountId,
			categoryIds,
			auxiliaryAccountIds,
			func(auxiliaryLedgers []*auxiliary_ledger.AuxiliaryLedger) ([]*auxiliary_ledger.AuxiliaryLedger, error) {
				// Build lookup map using composite key
				ledgerMap := make(map[auxiliaryLedgerKey]*auxiliary_ledger.AuxiliaryLedger)
				for _, l := range auxiliaryLedgers {
					key := auxiliaryLedgerKey{
						accountId:           l.AccountId(),
						auxiliaryCategoryId: l.AuxiliaryCategoryId(),
						auxiliaryAccountId:  l.AuxiliaryAccount().Id(),
					}
					ledgerMap[key] = l
				}

				// Apply updates
				for _, record := range records {
					key := auxiliaryLedgerKey{
						accountId:           record.accountId,
						auxiliaryCategoryId: record.auxiliaryCategoryId,
						auxiliaryAccountId:  record.auxiliaryAccount.Id(),
					}

					auxiliaryLedger, ok := ledgerMap[key]
					if !ok {
						return nil, fmt.Errorf(
							"auxiliary ledger not found for account=%s, category=%s, auxiliary=%s",
							record.accountId,
							record.auxiliaryCategoryId,
							record.auxiliaryAccount.Key(),
						)
					}

					// Update balance using domain method
					auxiliaryLedger.UpdateBalance(record.debit, record.credit)
				}

				return auxiliaryLedgers, nil
			},
		); err != nil {
			return fmt.Errorf("failed to update auxiliary ledgers for account %s: %w", accountId, err)
		}
	}

	return nil
}
