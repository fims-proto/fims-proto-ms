package command

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type PostVoucherCmd struct {
	VoucherUUID uuid.UUID
}

type PostVoucherHandler struct {
	readModel     query.VouchersReadModel
	ledgerService LedgerService
}

func NewPostVoucherHandler(readModel query.VouchersReadModel, ledgerService LedgerService) PostVoucherHandler {
	if readModel == nil {
		panic("nil read model")
	}
	return PostVoucherHandler{
		readModel:     readModel,
		ledgerService: ledgerService,
	}
}

func (h PostVoucherHandler) Handler(ctx context.Context, cmd PostVoucherCmd) error {
	voucher, err := h.readModel.VoucherByUUID(ctx, cmd.VoucherUUID)
	if err != nil {
		return errors.Wrap(err, "failed to read voucher while posting")
	}
	return h.ledgerService.PostVoucher(ctx, voucher)
}
