package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type PostVoucherCmd struct {
	VoucherUUID uuid.UUID
}

type PostVoucherHandler struct {
	readModel     query.VouchersReadModel
	repo          domain.Repository
	ledgerService LedgerService
}

func NewPostVoucherHandler(readModel query.VouchersReadModel, repo domain.Repository, ledgerService LedgerService) PostVoucherHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return PostVoucherHandler{
		readModel:     readModel,
		repo:          repo,
		ledgerService: ledgerService,
	}
}

func (h PostVoucherHandler) Handler(ctx context.Context, cmd PostVoucherCmd) error {
	voucher, err := h.readModel.VoucherByUUID(ctx, cmd.VoucherUUID)
	if err != nil {
		return errors.Wrap(err, "failed to read voucher while posting")
	}
	if err = h.ledgerService.PostVoucher(ctx, voucher); err != nil {
		return errors.Wrap(err, "failed to post voucher to ledgers")
	}
	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			if err := v.Post(); err != nil {
				return nil, err
			}
			return v, err
		},
	)
}
