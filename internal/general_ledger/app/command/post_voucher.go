package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/utils"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/ledger"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"
)

type PostVoucherCmd struct {
	VoucherId uuid.UUID
	Poster    uuid.UUID
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
				return nil, errors.Wrap(err, "failed to post voucher to ledger")
			}

			return v, nil
		},
	)
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

func (h PostVoucherHandler) postLedgers(ctx context.Context, cmd postLedgersCmd) error {
	accountCommands := cmd.records
	for _, record := range cmd.records {
		//  read all superior accounts
		superiorAccounts, err := h.repo.ReadSuperiorAccountsById(ctx, record.accountId)
		if err != nil {
			return errors.Wrap(err, "failed to read superior accounts")
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

	// merge duplicated accounts
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
					return nil, errors.Errorf("should not happen, failed to find account %s in accountsMap", l.AccountId())
				}

				l.UpdateBalance(record.debit, record.credit)
			}

			return ledgers, nil
		},
	)
}
