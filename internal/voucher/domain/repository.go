package domain

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github.com/google/uuid"
)

type Repository interface {
	CreateVoucher(ctx context.Context, d *voucher.Voucher) error
	UpdateVoucher(
		ctx context.Context,
		voucherId uuid.UUID,
		updateFn func(d *voucher.Voucher) (*voucher.Voucher, error),
	) error

	Migrate(ctx context.Context) error
}
