package query

import (
	"context"

	"github.com/google/uuid"
)

type vouchersReadModel interface {
	AllVouchers(ctx context.Context) ([]Voucher, error)
	VoucherByUUID(ctx context.Context, uuid uuid.UUID) (Voucher, error)
}

type ReadVouchersHandler struct {
	readModel vouchersReadModel
}

func NewAllVouchersHandler(readModel vouchersReadModel) ReadVouchersHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	return ReadVouchersHandler{readModel: readModel}
}

func (h ReadVouchersHandler) HandleReadAll(ctx context.Context) ([]Voucher, error) {
	return h.readModel.AllVouchers(ctx)
}

func (h ReadVouchersHandler) HandleReadByUUID(uuid uuid.UUID, ctx context.Context) (Voucher, error) {
	return h.readModel.VoucherByUUID(ctx, uuid)
}
