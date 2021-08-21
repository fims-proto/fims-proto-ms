package db

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
)

type VoucherPostgresRepository struct{}

func NewVoucherPostgresRepository() *VoucherPostgresRepository {
	return &VoucherPostgresRepository{}
}

// implementation methods

func (r VoucherPostgresRepository) AddVoucher(ctx context.Context, v *domain.Voucher) (createdUUID uuid.UUID, err error) {
	panic("not implemented") // TODO: Implement
}

func (r VoucherPostgresRepository) UpdateVoucher(
	ctx context.Context,
	sob string,
	voucherUUID uuid.UUID,
	updateFn func(v *domain.Voucher) (*domain.Voucher, error),
) (err error) {
	panic("not implemented") // TODO: Implement
}

func (r VoucherPostgresRepository) AllVouchers(ctx context.Context, sob string) ([]query.Voucher, error) {
	panic("not implemented") // TODO: Implement
}

func (r VoucherPostgresRepository) VoucherByUUID(ctx context.Context, sob string, uuid uuid.UUID) (query.Voucher, error) {
	panic("not implemented") // TODO: Implement
}

func readDBFromCtx(ctx context.Context) uuid.UUID {
	return ctx.Value("db").(uuid.UUID)
}
