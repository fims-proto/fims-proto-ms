package db

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type voucher struct {
	Id                 uuid.UUID `gorm:"type:uuid"`
	SobId              string    `gorm:"uniqueIndex:uni_sobid_number"`
	Number             string    `gorm:"uniqueIndex:uni_sobid_number"`
	VoucherType        string
	AttachmentQuantity uint
	Debit              decimal.Decimal
	Credit             decimal.Decimal
	Creator            string
	Reviewer           string
	IsReviewed         bool
	Auditor            string
	IsAudited          bool
	IsPosted           bool
	LineItems          []lineItem
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          gorm.DeletedAt `gorm:"index"`
}

type lineItem struct {
	Id            uuid.UUID `gorm:"type:uuid"`
	VoucherId     uuid.UUID `gorm:"type:uuid"`
	Summary       string
	AccountNumber string
	Debit         decimal.Decimal
	Credit        decimal.Decimal
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// from domain to db
func marshall(dv *domain.Voucher) *voucher {
	v := voucher{
		Id:                 dv.Id(),
		SobId:              dv.Sob(),
		Number:             dv.Number(),
		VoucherType:        dv.Type().String(),
		AttachmentQuantity: dv.AttachmentQuantity(),
		Debit:              dv.Debit(),
		Credit:             dv.Credit(),
		Creator:            dv.Creator(),
		Reviewer:           dv.Reviewer(),
		IsReviewed:         dv.IsReviewed(),
		Auditor:            dv.Auditor(),
		IsAudited:          dv.IsAudited(),
		IsPosted:           dv.IsPosted(),
		LineItems:          []lineItem{},
	}
	for _, item := range dv.LineItems() {
		v.LineItems = append(v.LineItems, lineItem{
			Id:            item.Id(),
			VoucherId:     dv.Id(),
			Summary:       item.Summary(),
			AccountNumber: item.AccountNumber(),
			Debit:         item.Debit(),
			Credit:        item.Credit(),
		})
	}
	return &v
}

// from db to domain
func unmarshallToDomain(dbv *voucher) (*domain.Voucher, error) {
	items := []*domain.LineItem{}

	for _, dbItem := range dbv.LineItems {
		item, err := domain.NewLineItem(
			dbItem.Id,
			dbItem.Summary,
			dbItem.AccountNumber,
			dbItem.Debit.String(),
			dbItem.Credit.String(),
		)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshall lineItem failed")
		}
		items = append(items, item)
	}

	v, err := domain.NewVoucher(
		dbv.Id,
		dbv.SobId,
		dbv.VoucherType,
		dbv.Number,
		dbv.AttachmentQuantity,
		items,
		dbv.Creator,
		dbv.Reviewer,
		dbv.Auditor,
		dbv.IsReviewed,
		dbv.IsAudited,
		dbv.IsPosted,
	)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshall voucher failed")
	}

	return v, nil
}

func unmarshallToQuery(dbv *voucher) query.Voucher {
	items := []query.LineItem{}

	for _, dbItem := range dbv.LineItems {
		items = append(items, query.LineItem{
			Summary:       dbItem.Summary,
			AccountNumber: dbItem.AccountNumber,
			Debit:         dbItem.Debit.String(),
			Credit:        dbItem.Credit.String(),
		})
	}

	return query.Voucher{
		Sob:                dbv.SobId,
		UUID:               dbv.Id,
		VoucherType:        dbv.VoucherType,
		Number:             dbv.Number,
		AttachmentQuantity: dbv.AttachmentQuantity,
		LineItems:          items,
		Debit:              dbv.Debit.String(),
		Credit:             dbv.Credit.String(),
		Creator:            dbv.Creator,
		Reviewer:           dbv.Reviewer,
		IsReviewed:         dbv.IsReviewed,
		Auditor:            dbv.Auditor,
		IsAudited:          dbv.IsAudited,
		IsPosted:           dbv.IsPosted,
	}
}
