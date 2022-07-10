package db

import (
	"time"

	"github/fims-proto/fims-proto-ms/internal/voucher/app/query"
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type voucher struct {
	Id                 uuid.UUID `gorm:"type:uuid"`
	SobId              uuid.UUID `gorm:"type:uuid;uniqueIndex:vouchers_sob_period_number_key"`
	PeriodId           uuid.UUID `gorm:"type:uuid;uniqueIndex:vouchers_sob_period_number_key"`
	Number             string    `gorm:"uniqueIndex:vouchers_sob_period_number_key"`
	VoucherType        string
	AttachmentQuantity uint
	Debit              decimal.Decimal
	Credit             decimal.Decimal
	Creator            uuid.UUID `gorm:"type:uuid"`
	Reviewer           uuid.UUID `gorm:"type:uuid"`
	IsReviewed         bool
	Auditor            uuid.UUID `gorm:"type:uuid"`
	IsAudited          bool
	IsPosted           bool
	TransactionTime    time.Time
	LineItems          []lineItem
	CreatedAt          time.Time `gorm:"<-:create"`
	UpdatedAt          time.Time
}

type lineItem struct {
	Id        uuid.UUID `gorm:"type:uuid"`
	VoucherId uuid.UUID `gorm:"type:uuid"`
	AccountId uuid.UUID `gorm:"type:uuid"`
	Summary   string
	Debit     decimal.Decimal
	Credit    decimal.Decimal
}

// from domain to db
func marshall(dv *domain.Voucher) *voucher {
	v := voucher{
		Id:                 dv.Id(),
		SobId:              dv.SobId(),
		PeriodId:           dv.PeriodId(),
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
		TransactionTime:    dv.TransactionTime(),
	}
	for _, item := range dv.LineItems() {
		v.LineItems = append(v.LineItems, lineItem{
			Id:        item.Id(),
			VoucherId: dv.Id(),
			Summary:   item.Summary(),
			AccountId: item.AccountId(),
			Debit:     item.Debit(),
			Credit:    item.Credit(),
		})
	}
	return &v
}

// from db to domain
func unmarshallToDomain(dbv *voucher) (*domain.Voucher, error) {
	var items []*domain.LineItem

	for _, dbItem := range dbv.LineItems {
		item, err := domain.NewLineItem(
			dbItem.Id,
			dbItem.AccountId,
			dbItem.Summary,
			dbItem.Debit,
			dbItem.Credit,
		)
		if err != nil {
			return nil, errors.Wrap(err, "unmarshall lineItem failed")
		}
		items = append(items, item)
	}

	return domain.NewVoucher(
		dbv.Id,
		dbv.SobId,
		dbv.PeriodId,
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
		dbv.TransactionTime,
	)
}

func unmarshallToQuery(dbv *voucher) query.Voucher {
	var items []query.LineItem

	for _, dbItem := range dbv.LineItems {
		items = append(items, query.LineItem{
			Id:        dbItem.Id,
			AccountId: dbItem.AccountId,
			Summary:   dbItem.Summary,
			Debit:     dbItem.Debit,
			Credit:    dbItem.Credit,
		})
	}

	return query.Voucher{
		Id:                 dbv.Id,
		SobId:              dbv.SobId,
		Period:             query.Period{Id: dbv.PeriodId},
		VoucherType:        dbv.VoucherType,
		Number:             dbv.Number,
		AttachmentQuantity: dbv.AttachmentQuantity,
		LineItems:          items,
		Debit:              dbv.Debit,
		Credit:             dbv.Credit,
		Creator:            query.User{Id: dbv.Creator},
		Reviewer:           query.User{Id: dbv.Reviewer},
		IsReviewed:         dbv.IsReviewed,
		Auditor:            query.User{Id: dbv.Auditor},
		IsAudited:          dbv.IsAudited,
		IsPosted:           dbv.IsPosted,
		TransactionTime:    dbv.TransactionTime,
		CreatedAt:          dbv.CreatedAt,
		UpdatedAt:          dbv.UpdatedAt,
	}
}
