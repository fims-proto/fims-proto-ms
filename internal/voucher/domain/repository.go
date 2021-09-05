package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	AddVoucher(ctx context.Context, v *Voucher) (uuid.UUID, error)
	UpdateVoucher(ctx context.Context, id uuid.UUID, updateFn func(v *Voucher) (*Voucher, error)) error
	Migrate(ctx context.Context) error
}
