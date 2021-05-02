package query

import (
	"time"

	"github.com/google/uuid"
)

type LineItem struct {
	Summary       string
	AccountNumber string
	Debit         string
	Credit        string
}

type Voucher struct {
	UUID               uuid.UUID
	VoucherType        string
	Number             string
	CreatedAt          time.Time
	AttachmentQuantity uint
	LineItems          []LineItem
	Debit              string
	Credit             string
	Creator            string
	Reviewer           string
	IsReviewed         bool
	Auditor            string
	IsAudited          bool
	IsPosted           bool
}
