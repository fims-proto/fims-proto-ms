package db

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"

	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VoucherPostgresRepository struct{}

func NewVoucherPostgresRepository() *VoucherPostgresRepository {
	return &VoucherPostgresRepository{}
}

func (r VoucherPostgresRepository) Migrate(ctx context.Context) error {
	db := readDBFromCtx(ctx)

	if err := db.AutoMigrate(&voucher{}, &lineItem{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r VoucherPostgresRepository) CreateVoucher(ctx context.Context, v *domain.Voucher) (uuid.UUID, error) {
	db := readDBFromCtx(ctx)

	dbVoucher := marshall(v)

	if err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(dbVoucher).Error
	}); err != nil {
		return uuid.Nil, errors.Wrap(err, "create voucher failed")
	}

	return v.Id(), nil
}

func (r VoucherPostgresRepository) UpdateVoucher(ctx context.Context, id uuid.UUID, updateFn func(voucher *domain.Voucher) (*domain.Voucher, error)) error {
	db := readDBFromCtx(ctx)

	if err := db.Transaction(func(tx *gorm.DB) error {
		dbVoucher := &voucher{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("LineItems").First(dbVoucher, "id = ?", id).Error; err != nil {
			return err
		}

		voucher, err := unmarshallToDomain(dbVoucher)
		if err != nil {
			return errors.Wrap(err, "unmarshall voucher failed")
		}

		updatedVoucher, err := updateFn(voucher)
		if err != nil {
			return errors.Wrap(err, "update voucher in transaction failed")
		}

		dbVoucher = marshall(updatedVoucher)
		// remove existing first
		if err := tx.Where("voucher_id = ?", dbVoucher.Id).Delete(&lineItem{}).Error; err != nil {
			return errors.Wrap(err, "delete voucher items failed")
		}
		if err := tx.Save(dbVoucher).Error; err != nil {
			return errors.Wrap(err, "save voucher failed")
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "update voucher failed")
	}

	return nil
}

func (r VoucherPostgresRepository) ReadAllVouchers(ctx context.Context, sobId uuid.UUID, pageable data.Pageable) (data.Page[query.Voucher], error) {
	db := readDBFromCtx(ctx)

	var dbVouchers []voucher

	db = data.EnrichDb(pageable, db).Where("sob_id = ?", sobId)

	if err := db.Preload("LineItems").Find(&dbVouchers).Error; err != nil {
		return data.Page[query.Voucher]{}, errors.Wrap(err, "find vouchers by sob failed")
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return data.Page[query.Voucher]{}, errors.Wrap(err, "count vouchers failed")
	}

	var qvs []query.Voucher
	for _, dbVoucher := range dbVouchers {
		qvs = append(qvs, unmarshallToQuery(&dbVoucher))
	}
	return data.NewPage(qvs, pageable.Page(), pageable.Size(), int(count))
}

func (r VoucherPostgresRepository) ReadById(ctx context.Context, uuid uuid.UUID) (query.Voucher, error) {
	db := readDBFromCtx(ctx)

	dbVoucher := voucher{}
	if err := db.Preload("LineItems").First(&dbVoucher, "id = ?", uuid).Error; err != nil {
		return query.Voucher{}, errors.Wrap(err, "find voucher by uuid failed")
	}

	return unmarshallToQuery(&dbVoucher), nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
