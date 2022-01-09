package domain

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	CreateVoucher(ctx context.Context, v *Voucher) (uuid.UUID, error)
	UpdateVoucher(ctx context.Context, id uuid.UUID, updateFn func(voucher *Voucher) (*Voucher, error)) error
	Migrate(ctx context.Context) error
}
