package voucher

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type NotFoundError struct {
	VoucherUUID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("voucher %s not found", e.VoucherUUID)
}

type Repository interface {
	AddVoucher(ctx context.Context, v *Voucher) (uuid.UUID, error)
	UpdateVoucher(
		ctx context.Context,
		voucherUUID uuid.UUID,
		updateFn func(v *Voucher) (*Voucher, error),
	) error
}
