package db

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type voucher struct {
	sob                string
	uuid               uuid.UUID
	voucherType        string
	number             string
	createdAt          time.Time
	attachmentQuantity uint
	debit              decimal.Decimal
	credit             decimal.Decimal
	creator            string
	reviewer           string
	isReviewed         bool
	auditor            string
	isAudited          bool
	isPosted           bool
	// one to many lineitems
}

type lineItem struct {
	voucherUUID   uuid.UUID
	summary       string
	accountNumber string
	debit         decimal.Decimal
	credit        decimal.Decimal
}

func marshallFromDomain(dv *domain.Voucher) (voucher, []lineItem) {
	v := voucher{
		sob:                dv.Sob(),
		uuid:               dv.UUID(),
		voucherType:        dv.Type().String(),
		number:             dv.Number(),
		createdAt:          dv.CreatedAt(),
		attachmentQuantity: dv.AttachmentQuantity(),
		debit:              dv.Debit(),
		credit:             dv.Credit(),
		creator:            dv.Creator(),
		reviewer:           dv.Reviewer(),
		isReviewed:         dv.IsReviewed(),
		auditor:            dv.Auditor(),
		isAudited:          dv.IsAudited(),
		isPosted:           dv.IsPosted(),
	}
	items := []lineItem{}
	for _, item := range dv.LineItems() {
		items = append(items, lineItem{
			voucherUUID:   dv.UUID(),
			summary:       item.Summary(),
			accountNumber: item.AccountNumber(),
			debit:         item.Debit(),
			credit:        item.Credit(),
		})
	}
	return v, items
}
