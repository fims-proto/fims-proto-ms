package query

import "context"

type AllVouchersReadModel interface {
	AllVouchers(ctx context.Context) ([]Voucher, error)
}

type AllVouchersHandler struct {
	readModel AllVouchersReadModel
}

func NewAllVouchersHandler(readModel AllVouchersReadModel) AllVouchersHandler {
	if readModel == nil {
		panic("nil readModel")
	}
	return AllVouchersHandler{readModel: readModel}
}

func (h AllVouchersHandler) Handle(ctx context.Context) ([]Voucher, error) {
	return h.readModel.AllVouchers(ctx)
}