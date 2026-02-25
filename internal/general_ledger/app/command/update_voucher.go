package command

import (
	"context"
	"fmt"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/transaction_date"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
)

type UpdateVoucherCmd struct {
	VoucherId       uuid.UUID
	HeaderText      string
	LineItems       []LineItemCmd
	TransactionDate transaction_date.TransactionDate
	Updater         uuid.UUID
}

type UpdateVoucherHandler struct {
	repo             domain.Repository
	numberingService service.NumberingService
}

func NewUpdateVoucherHandler(repo domain.Repository, numberingService service.NumberingService) UpdateVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return UpdateVoucherHandler{
		repo:             repo,
		numberingService: numberingService,
	}
}

func (h UpdateVoucherHandler) Handle(ctx context.Context, cmd UpdateVoucherCmd) error {
	return h.repo.EnableTx(ctx, func(txCtx context.Context) error {
		return h.updateVoucher(txCtx, cmd)
	})
}

func (h UpdateVoucherHandler) updateVoucher(ctx context.Context, cmd UpdateVoucherCmd) error {
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherId,
		func(v *voucher.Voucher) (*voucher.Voucher, error) {
			// update line items
			if len(cmd.LineItems) > 0 {
				lineItems, err := prepareLineItems(ctx, h.repo, v.SobId(), cmd.LineItems)
				if err != nil {
					return nil, fmt.Errorf("failed to prepare line items: %w", err)
				}

				if err = v.UpdateLineItems(lineItems, cmd.Updater); err != nil {
					return nil, err
				}
			}

			// update transaction date (and period and document number, if needed)
			if !cmd.TransactionDate.IsZero() {
				p, err := readPeriodIdAndCheck(ctx, h.repo, h.numberingService, v.SobId(), cmd.TransactionDate)
				if err != nil {
					return nil, fmt.Errorf("failed to read or create period: %w", err)
				}

				if p.Id() != v.PeriodId() {
					// different period, need to regenerate voucher id
					identifier, err := h.numberingService.GenerateIdentifier(ctx, p.Id(), v.VoucherType().String())
					if err != nil {
						return nil, fmt.Errorf("failed to re-generate voucher number: %w", err)
					}
					if err = v.UpdatePeriodAndDocumentNumber(p.Id(), identifier, cmd.Updater); err != nil {
						return nil, err
					}
				}

				if err = v.UpdateTransactionDate(cmd.TransactionDate, cmd.Updater); err != nil {
					return nil, err
				}
			}

			if cmd.HeaderText != "" {
				if err := v.UpdateHeaderText(cmd.HeaderText, cmd.Updater); err != nil {
					return nil, err
				}
			}

			return v, nil
		},
	)
}
