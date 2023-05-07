package command

import (
	"context"
	"time"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/query"
	"github/fims-proto/fims-proto-ms/internal/general_ledger/app/service"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/general_ledger/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type UpdateVoucherCmd struct {
	VoucherId       uuid.UUID
	HeaderText      string
	LineItems       []LineItemCmd
	TransactionTime time.Time
	Updater         uuid.UUID
}

type UpdateVoucherHandler struct {
	repo             domain.Repository
	readModel        query.GeneralLedgerReadModel
	numberingService service.NumberingService
}

func NewUpdateVoucherHandler(
	repo domain.Repository,
	readModel query.GeneralLedgerReadModel,
	numberingService service.NumberingService,
) UpdateVoucherHandler {
	if repo == nil {
		panic("nil repo")
	}

	if readModel == nil {
		panic("nil read model")
	}

	if numberingService == nil {
		panic("nil numbering service")
	}

	return UpdateVoucherHandler{
		repo:             repo,
		readModel:        readModel,
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
				lineItems, err := prepareLineItems(ctx, h.readModel, v.SobId(), cmd.LineItems)
				if err != nil {
					return nil, errors.Wrap(err, "failed to prepare line items")
				}

				if err = v.UpdateLineItems(lineItems, cmd.Updater); err != nil {
					return nil, err
				}
			}

			// update transaction time (and period and document number, if needed)
			if !cmd.TransactionTime.IsZero() {
				periodId, err := readOrCreatePeriodForVoucher(ctx, h.repo, h.readModel, h.numberingService, v.SobId(), cmd.TransactionTime)
				if err != nil {
					return nil, errors.Wrap(err, "failed to read or create period")
				}

				if periodId != v.PeriodId() {
					// different period, need to regenerate voucher id
					identifier, err := h.numberingService.GenerateIdentifier(ctx, periodId, v.VoucherType().String())
					if err != nil {
						return nil, errors.Wrap(err, "failed to re-generate voucher number")
					}
					if err = v.UpdatePeriodAndDocumentNumber(periodId, identifier, cmd.Updater); err != nil {
						return nil, err
					}
				}

				if err = v.UpdateTransactionTime(cmd.TransactionTime, cmd.Updater); err != nil {
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
