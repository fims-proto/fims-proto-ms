package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/ledger/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type LineItemCmd struct {
	AccountNumber string
	Debit         decimal.Decimal
	Credit        decimal.Decimal
}

type UpdateLedgerBalanceCmd struct {
	Sob         string
	VoucherUUID uuid.UUID
	LineItems   []LineItemCmd
}

type UpdateLedgerBalanceHandler struct {
	repo           domain.Repository
	accountService AccountService
	voucherService VoucherService
}

func NewUpdateLedgerBalanceHandler(repo domain.Repository, accountService AccountService, voucherService VoucherService) UpdateLedgerBalanceHandler {
	if repo == nil {
		panic("nil repo")
	}
	if accountService == nil {
		panic("nil account service")
	}
	if voucherService == nil {
		panic("nil voucher service")
	}
	return UpdateLedgerBalanceHandler{
		repo:           repo,
		accountService: accountService,
		voucherService: voucherService,
	}
}

func (h UpdateLedgerBalanceHandler) Handle(ctx context.Context, cmd UpdateLedgerBalanceCmd) error {
	// 1. check again if voucher already posted
	voucherPosted, err := h.voucherService.CheckVoucherPosted(ctx, cmd.Sob, cmd.VoucherUUID)
	if err != nil {
		return errors.Wrap(err, "check voucher posted failed")
	}
	if voucherPosted {
		return errors.New("voucher already posted")
	}

	// 2. get all superior account numbers for lineitems
	numberItemMap := make(map[string]LineItemCmd)
	var allLedgerNums []string
	for _, item := range cmd.LineItems {
		accNums, err := h.accountService.ReadSuperiorNumbers(ctx, cmd.Sob, item.AccountNumber)
		if err != nil {
			return errors.Wrap(err, "read account superior number failed")
		}

		for _, accNum := range accNums {
			numberItemMap[accNum] = item
		}
		allLedgerNums = append(allLedgerNums, accNums...)
	}

	if len(allLedgerNums) <= 0 {
		return errors.New("no ledgers to be updated")
	}

	// 3. do the update
	return h.repo.UpdateLedgers(
		ctx,
		cmd.Sob,
		allLedgerNums,
		func(ledgers []*domain.Ledger) ([]*domain.Ledger, error) {
			if len(ledgers) <= 0 {
				return nil, errors.New("no ledgers found")
			}

			// use sequencial process for now, as the data volumne should not be so large
			// say 5 line items in one voucher, each line item should update 4 ledgers
			// then it's 20 ledger entries to update
			// let's see if this is gonna be bottle neck of performance, we change to goroutine in parallel
			for _, l := range ledgers {
				item, ok := numberItemMap[l.Number()]
				if !ok {
					return nil, errors.Errorf("shoul not happen: debit/credit not mapped to account number %s", l.Number())
				}
				if err := l.UpdateBalance(item.Debit, item.Credit); err != nil {
					return nil, errors.Wrapf(err, "update balance of ledger %s failed", l.Number())
				}
			}
			return ledgers, nil
		},
	)
}
