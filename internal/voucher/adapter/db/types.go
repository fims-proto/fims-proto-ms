package db

import (
	"github/fims-proto/fims-proto-ms/internal/voucher/domain"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type voucher struct {
	sobId              string          `db:"sob_id"`
	uuid               uuid.UUID       `db:"voucher_id"`
	voucherType        string          `db:"voucher_type"`
	number             string          `db:"number"`
	createdAt          time.Time       `db:"created_at"`
	attachmentQuantity uint            `db:"attachment_quantity"`
	debit              decimal.Decimal `db:"debit"`
	credit             decimal.Decimal `db:"credit"`
	creator            string          `db:"creator"`
	reviewer           string          `db:"reviewer"`
	isReviewed         bool            `db:"is_reviewed"`
	auditor            string          `db:"auditor"`
	isAudited          bool            `db:"is_audited"`
	isPosted           bool            `db:"is_posted"`
	// one to many lineitems
}

type lineItem struct {
	voucherUUID   uuid.UUID       `db:"voucher_id"`
	summary       string          `db:"summary"`
	accountNumber string          `db:"account_number"`
	debit         decimal.Decimal `db:"debit"`
	credit        decimal.Decimal `db:"credit"`
}

func marshallFromDomain(dv *domain.Voucher) (voucher, []lineItem) {
	v := voucher{
		sobId:              dv.Sob(),
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
