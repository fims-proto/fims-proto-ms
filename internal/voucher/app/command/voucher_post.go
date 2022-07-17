package command

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/voucher/app/service"

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
	ledgerService service.LedgerService
}

func NewPostVoucherHandler(readModel query.VouchersReadModel, repo domain.Repository, ledgerService service.LedgerService) PostVoucherHandler {
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

	return h.repo.UpdateVoucher(
		ctx,
		cmd.VoucherUUID,
		func(v *domain.Voucher) (*domain.Voucher, error) {
			log.Info(ctx, "posting voucher")
			if err := v.Post(); err != nil {
				return nil, err
			}

			if err = h.ledgerService.PostVoucher(ctx, *v); err != nil {
				return nil, errors.Wrap(err, "failed to post voucher to ledgers")
			}

			return v, err
		},
	)
}
