package db

import (
	"context"

	"github/fims-proto/fims-proto-ms/internal/common/data"
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"

	"github/fims-proto/fims-proto-ms/internal/voucher/domain/voucher"

	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"

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

	if err := db.AutoMigrate(&voucherPO{}, &lineItemPO{}); err != nil {
		return errors.Wrap(err, "DB migration failed")
	}
	return nil
}

func (r VoucherPostgresRepository) CreateVoucher(ctx context.Context, d *voucher.Voucher) error {
	db := readDBFromCtx(ctx)

	po := voucherBOToPO(*d)

	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&po).Error
	})
}

func (r VoucherPostgresRepository) UpdateVoucher(ctx context.Context, voucherId uuid.UUID, updateFn func(d *voucher.Voucher) (*voucher.Voucher, error)) error {
	db := readDBFromCtx(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		po := voucherPO{}
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("LineItems").First(&po, "id = ?", voucherId).Error; err != nil {
			return err
		}

		bo, err := voucherPOToBO(po)
		if err != nil {
			return errors.Wrap(err, "failed to map voucher")
		}

		updatedBO, err := updateFn(bo)
		if err != nil {
			return errors.Wrap(err, "update voucher in transaction failed")
		}

		po = voucherBOToPO(*updatedBO)

		// remove existing first
		if err = tx.Where("voucher_id = ?", po.Id).Delete(&lineItemPO{}).Error; err != nil {
			return errors.Wrap(err, "delete voucher items failed")
		}

		return tx.Save(&po).Error
	})
}

// queries

func (r VoucherPostgresRepository) SearchVouchers(ctx context.Context, sobId uuid.UUID, pageRequest data.PageRequest) (data.Page[query.Voucher], error) {
	if sobId != uuid.Nil {
		sobIdFilter, _ := filterable.NewFilter("sobId", filterable.OptEq, sobId.String())
		sobIdFilterable := filterable.NewFilterableAtom(sobIdFilter)
		pageRequest.AddAndFilterable(sobIdFilterable)
	}
	return data.SearchEntities(ctx, pageRequest, voucherPO{}, voucherPOToDTO, readDBFromCtx(ctx).Preload("LineItems"))
}

func (r VoucherPostgresRepository) VoucherById(ctx context.Context, voucherId uuid.UUID) (query.Voucher, error) {
	voucherIdFilter, _ := filterable.NewFilter("id", filterable.OptEq, voucherId)
	pageRequest := data.NewPageRequest(pageable.Unpaged(), sortable.Unsorted(), filterable.NewFilterableAtom(voucherIdFilter))

	vouchers, err := r.SearchVouchers(ctx, uuid.Nil, pageRequest)
	if err != nil {
		return query.Voucher{}, err
	}

	if vouchers.NumberOfElements() != 1 {
		return query.Voucher{}, errors.Errorf("voucher not found by id: %s", voucherId)
	}

	return vouchers.Content()[0], nil
}

func readDBFromCtx(ctx context.Context) *gorm.DB {
	return ctx.Value("db").(*gorm.DB)
}
