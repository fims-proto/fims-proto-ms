package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_account"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/auxiliary_ledger"

	"github/fims-proto/fims-proto-ms/internal/common/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"
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
	periodId uuid.UUID
	records  []postAuxiliaryLedgersRecordCmd
}

type postAuxiliaryLedgersRecordCmd struct {
	auxiliaryAccount auxiliary_account.AuxiliaryAccount
	debit            decimal.Decimal
	credit           decimal.Decimal
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
						auxiliaryAccount: *auxiliaryAccount,
						debit:            item.Debit(),
						credit:           item.Credit(),
					})
				}
			}

			if err := h.postAuxiliaryLedgers(ctx, postAuxiliaryLedgersCmd{
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

	// merge same auxiliary accounts
	auxiliaryAccountsMap := utils.SliceToMapMerge(cmd.records, func(c postAuxiliaryLedgersRecordCmd) uuid.UUID {
		return c.auxiliaryAccount.Id()
	}, func(c postAuxiliaryLedgersRecordCmd) postAuxiliaryLedgersRecordCmd {
		return c
	}, func(existing, replacement postAuxiliaryLedgersRecordCmd) postAuxiliaryLedgersRecordCmd {
		existing.debit = existing.debit.Add(replacement.debit)
		existing.credit = existing.credit.Add(replacement.credit)
		return existing
	})

	return h.repo.UpdateAuxiliaryLedgersByPeriodAndAccountIds(
		ctx,
		cmd.periodId,
		utils.MapToKeySlice(auxiliaryAccountsMap),
		func(auxiliaryLedgers []*auxiliary_ledger.AuxiliaryLedger) ([]*auxiliary_ledger.AuxiliaryLedger, error) {
			for _, l := range auxiliaryLedgers {
				record, ok := auxiliaryAccountsMap[l.AuxiliaryAccount().Id()]
				if !ok {
					return nil, fmt.Errorf("should not happen, failed to find auxiliary account %s", l.AuxiliaryAccount().Key())
				}

				l.UpdateBalance(record.debit, record.credit)
			}

			return auxiliaryLedgers, nil
		},
	)
}
