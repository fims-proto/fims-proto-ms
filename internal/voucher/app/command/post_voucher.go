package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/log"
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

func (h PostVoucherHandler) Handle(ctx context.Context, cmd PostVoucherCmd) (err error) {
	log.Info(ctx, "handle posting voucher")
	log.Debug(ctx, "handle posting voucher, cmd: %+v", cmd)
	defer func() {
		if err != nil {
			log.Err(ctx, err, "handle posting failed")
		}
	}()

	voucher, err := h.readModel.ReadById(ctx, cmd.VoucherUUID)
	if err != nil {
		return errors.Wrap(err, "failed to read voucher while posting")
	}

	if !voucher.IsReviewed {
		return errors.Errorf("voucher %s not reviewed", cmd.VoucherUUID)
	}

	if !voucher.IsAudited {
		return errors.Errorf("voucher %s not audited", cmd.VoucherUUID)
	}

	if voucher.IsPosted {
		return errors.Errorf("voucher %s is already posted", cmd.VoucherUUID)
	}

	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			log.Info(ctx, "posting voucher")
			if err := v.Post(); err != nil {
				return nil, err
			}

			if err = h.ledgerService.PostVoucher(ctx, voucher); err != nil {
				return nil, errors.Wrap(err, "failed to post voucher to ledgers")
			}

			return v, err
		},
	)
}
