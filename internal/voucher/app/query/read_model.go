package query

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github.com/google/uuid"
)

type VoucherReadModel interface {
	SearchVouchers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[Voucher], error)

	VoucherById(ctx context.Context, voucherId uuid.UUID) (Voucher, error)
}
