package db

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type VoucherPostgresRepository struct{}

func NewVoucherPostgresRepository() *VoucherPostgresRepository {
	return &VoucherPostgresRepository{}
}

// implementation methods

func (r VoucherPostgresRepository) AddVoucher(ctx context.Context, v *domain.Voucher) (createdUUID uuid.UUID, err error) {
	db := readDBFromCtx(ctx)

	dbVoucher := marshallFromDomain(v)

	if err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&dbVoucher).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return uuid.Nil, errors.Wrap(err, "create voucher failed")
	}

	return v.UUID(), nil
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

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
