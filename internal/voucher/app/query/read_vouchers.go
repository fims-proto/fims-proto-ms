package query

import (
	"context"

	"github.com/google/uuid"
)

type VouchersReadModel interface {
	ReadAllVouchers(ctx context.Context, sobId uuid.UUID) ([]Voucher, error)
	ReadById(ctx context.Context, id uuid.UUID) (Voucher, error)
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

func (h ReadVouchersHandler) HandleReadAll(ctx context.Context, sobId uuid.UUID) ([]Voucher, error) {
	return h.readModel.ReadAllVouchers(ctx, sobId)
}

func (h ReadVouchersHandler) HandleReadById(ctx context.Context, id uuid.UUID) (Voucher, error) {
	return h.readModel.ReadById(ctx, id)
}
