package db

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"time"

	"github.com/google/uuid"
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

func marshallFromDomain(dv *domain.Voucher) voucher {
	v := voucher{
		Id:                 dv.UUID(),
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
			VoucherId:     dv.UUID(),
			Summary:       item.Summary(),
			AccountNumber: item.AccountNumber(),
			Debit:         item.Debit(),
			Credit:        item.Credit(),
		})
	}
	return v
}
