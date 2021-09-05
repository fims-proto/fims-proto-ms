package query

import (
	"context"

	"github.com/google/uuid"
)

type VouchersReadModel interface {
	ReadAllVouchers(ctx context.Context, sob string) ([]Voucher, error)
	ReadByUUID(ctx context.Context, uuid uuid.UUID) (Voucher, error)
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
	return h.readModel.ReadAllVouchers(ctx, sob)
}

func (h ReadVouchersHandler) HandleReadByUUID(ctx context.Context, uuid uuid.UUID) (Voucher, error) {
	return h.readModel.ReadByUUID(ctx, uuid)
}
