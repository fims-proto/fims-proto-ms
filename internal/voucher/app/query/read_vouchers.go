package query

import (
	"context"

	"github.com/google/uuid"
)

type VouchersReadModel interface {
	AllVouchers(ctx context.Context, sob string) ([]Voucher, error)
	VoucherByUUID(ctx context.Context, sob string, uuid uuid.UUID) (Voucher, error)
}

type ReadVouchersHandler struct {
	readModel VouchersReadModel
}

func NewReadVouchersHandler(readModel VouchersReadModel) ReadVouchersHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	return ReadVouchersHandler{readModel: readModel}
}

func (h ReadVouchersHandler) HandleReadAll(ctx context.Context, sob string) ([]Voucher, error) {
	return h.readModel.AllVouchers(ctx, sob)
}

func (h ReadVouchersHandler) HandleReadByUUID(ctx context.Context, sob string, uuid uuid.UUID) (Voucher, error) {
	return h.readModel.VoucherByUUID(ctx, sob, uuid)
}
