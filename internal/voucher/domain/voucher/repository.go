package voucher

import (
	"context"
	"fmt"
)

type NotFoundError struct {
	VoucherUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("voucher %s not found", e.VoucherUUID)
}

type Repository interface {
	AddVoucher(ctx context.Context, v *Voucher) error
	UpdateVoucher(
		ctx context.Context,
		voucherUUID string,
		updateFn func(v *Voucher) (*Voucher, error),
	) error
}
