package query

import "context"

type vouchersReadModel interface {
	AllVouchers(ctx context.Context) ([]Voucher, error)
	VoucherByUUID(ctx context.Context, uuid string) (Voucher, error)
}

type ReadVouchersHandler struct {
	readModel vouchersReadModel
}

func NewReadVouchersHandler(readModel vouchersReadModel) ReadVouchersHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	return ReadVouchersHandler{readModel: readModel}
}

func (h ReadVouchersHandler) HandleReadAll(ctx context.Context) ([]Voucher, error) {
	return h.readModel.AllVouchers(ctx)
}

func (h ReadVouchersHandler) HandleReadForUUID(uuid string, ctx context.Context) (Voucher, error) {
	return h.readModel.VoucherByUUID(ctx, uuid)
}
