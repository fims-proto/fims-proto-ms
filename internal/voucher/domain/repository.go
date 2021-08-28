package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	AddVoucher(ctx context.Context, v *Voucher) (uuid.UUID, error)
	UpdateVoucher(
		ctx context.Context,
		voucherUUID uuid.UUID,
		updateFn func(v *Voucher) (*Voucher, error),
	) error
}
